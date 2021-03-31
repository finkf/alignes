package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"git.sr.ht/~flobar/lev"
)

var args = struct {
	ocrext, gtext string
}{}

func usage(prog string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [Options] [JSON...]\nOptions:\n", prog)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.StringVar(&args.ocrext, "ocrext", ".pred.txt", "set file extension for input ocr files")
	flag.StringVar(&args.gtext, "gtext", ".gt.txt", "set file extension for output gt files")
	flag.Usage = usage(os.Args[0])
	flag.Parse()
	for i := 1; i < len(os.Args); i++ {
		chk(align(os.Args[i]))
	}
}

func align(name string) error {
	d, err := readJSON(name)
	if err != nil {
		return fmt.Errorf("align: %v", err)
	}
	dir := d["Dir"].(string)
	if !exists(dir) {
		log.Printf("warning: directory %s does not exit; skipping", dir)
		return nil
	}
	files, ocr, err := gatherOCRFiles(d["Dir"].(string))
	if err != nil {
		return fmt.Errorf("align: %v", err)
	}
	gt := strings.Split(d["Text"].(string), "\n")
	m, trace := alignLines(gt, ocr)
	m.print(os.Stdout, gt, ocr)
	var as []alignment
	i, j := 0, 0
	for _, t := range trace {
		switch t {
		case '#':
			as = append(as, alignment{
				BaseName: files[i],
				OCR:      ocr[i],
				GT:       gt[j],
				Distance: lev.Distance(ocr[i], gt[j]),
			})
			i++
			j++
		case 'd':
			i++
		case 'i':
			j++
		default:
			panic("bad trace")
		}
	}
	d["Alignments"] = as
	for i := range as {
		log.Printf("%s GT:  %s", as[i].BaseName, as[i].GT)
		log.Printf("%s OCR: %s", as[i].BaseName, as[i].OCR)
		if err := ioutil.WriteFile(as[i].BaseName+args.gtext, []byte(as[i].GT+"\n"), 0666); err != nil {
			return err
		}
	}
	return writeJSON(name, d)
}

func gatherOCRFiles(dir string) ([]string, []string, error) {
	var files []string
	err := filepath.Walk(dir, func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(name, args.ocrext) {
			return nil
		}
		files = append(files, name[0:len(name)-len(args.ocrext)])
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("gather ocr files %s: %v", dir, err)
	}
	ocr := make([]string, len(files))
	for i := range files {
		line, err := readOCRFile(files[i] + args.ocrext)
		if err != nil {
			return nil, nil, fmt.Errorf("gather ocr files %s: %v", dir, err)
		}
		ocr[i] = line
	}
	return files, ocr, nil
}

func exists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

type mat struct {
	r, c int
	tab  []int
}

func newMat(r, c int) mat {
	return mat{r: r, c: c, tab: make([]int, r*c)}
}

func (m mat) at(i, j int) int {
	idx := i*m.c + j
	if idx >= len(m.tab) {
		return math.MaxInt32
	}
	return m.tab[i*m.c+j]
}

func (m mat) set(i, j, val int) int {
	m.tab[i*m.c+j] = val
	return val
}

func (m mat) trace() string {
	var x []byte
	for i, j := m.r-1, m.c-1; i > 0 && j > 0; {
		switch m.at(i, j) {
		case 0:
			i--
			j--
			x = append(x, '#')
		case 1:
			i--
			x = append(x, 'd')
		case 2:
			j--
			x = append(x, 'i')
		default:
			panic("bad entry")
		}
	}
	for i, j := 0, len(x); i < j; i, j = i+1, j-1 {
		x[i], x[j-1] = x[j-1], x[i]
	}
	return string(x)
}

func (m mat) print(out io.Writer, gt, ocr []string) {
	max := 0
	for i := range ocr {
		if len(tostr(ocr[i], 10)) > max {
			max = len(tostr(ocr[i], 10))
		}
	}
	for i := range gt {
		if len(tostr(gt[i], 10)) > max {
			max = len(tostr(gt[i], 10))
		}
	}
	var w tabwriter.Writer
	w.Init(out, 1, 8, 1, ' ', 0)
	defer w.Flush()
	fmt.Fprint(&w, " \t ")
	for i := range gt {
		fmt.Fprintf(&w, "\t%s", tostr(gt[i], 10))
	}
	fmt.Fprintln(&w)
	for i := 0; i < m.r; i++ {
		if i == 0 {
			fmt.Fprint(&w, " ")
		} else {
			fmt.Fprintf(&w, "%s", tostr(ocr[i-1], 10))
		}
		for j := 0; j < m.c; j++ {
			fmt.Fprintf(&w, "\t%d", m.at(i, j))
		}
		fmt.Fprintln(&w)
	}
}

func tostr(str string, n int) string {
	if len(str) > n {
		return str[:n-3] + "..."
	}
	return str
}

func alignLines(gt, ocr []string) (mat, string) {
	m := newMat(len(ocr)+1, len(gt)+1)
	t := newMat(len(ocr)+1, len(gt)+1)
	for i := range ocr {
		m.set(i+1, 0, len(ocr[i])+m.at(i, 0))
		t.set(i+1, 0, 1)
	}
	for i := range gt {
		m.set(0, i+1, len(gt[i])+m.at(0, i))
		t.set(i+1, 0, 2)
	}
	for i := 1; i < m.r; i++ {
		for j := 1; j < m.c; j++ {
			a := m.at(i-1, j-1) + lev.Distance(gt[j-1], ocr[i-1])
			b := m.at(i-1, j) + len(ocr[i-1])
			c := m.at(i, j-1) + len(gt[j-1])
			min, pos := min(a, b, c)
			m.set(i, j, min)
			t.set(i, j, pos)
		}
	}
	return m, t.trace()
}

func readOCRFile(name string) (string, error) {
	in, err := os.Open(name)
	if err != nil {
		return "", fmt.Errorf("read ocr file %s: %v", name, err)
	}
	defer in.Close()
	line, err := ioutil.ReadAll(in)
	if err != nil {
		return "", fmt.Errorf("read ocr file %s: %v", name, err)
	}
	return string(line), nil
}

func writeJSON(name string, data interface{}) (err error) {
	out, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("write json %s: %v", name, err)
	}
	defer func() {
		if err != nil {
			err = out.Close()
		}
	}()
	if err := json.NewEncoder(out).Encode(data); err != nil {
		return fmt.Errorf("write json %s: encode: %v", name, err)
	}
	return nil
}

func readJSON(name string) (map[string]interface{}, error) {
	in, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("read json %s: %v", name, err)
	}
	defer in.Close()
	d := make(map[string]interface{})
	if json.NewDecoder(in).Decode(&d); err != nil {
		return nil, fmt.Errorf("read json %s: decode: %v", name, err)
	}
	return d, nil
}

type alignment struct {
	BaseName, GT, OCR string
	Distance          int
}

func min(xs ...int) (int, int) {
	min := xs[0]
	arg := 0
	for i, x := range xs[1:] {
		if x < min {
			min = x
			arg = i + 1
		}
	}
	return min, arg
}

func chk(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
