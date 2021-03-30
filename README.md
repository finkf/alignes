# alignes
Align ocr lines based on ground-truth.

## Usage
```
alignes [Options] [JSON...]
Options:
  -ocrext set file extension of input ocr files (default ".pred.txt")
  -gtext set file extension of output gt files (default ".gt.txt")
```

Aligns ocr lines with the ground-truth lines in the region directory
and write the ground-truth files.  The region directory and the ground
truth lines are read from the json files.  The resulting alignments
are written back to the json files.

## Installation
To install just type `go get github.com/finkf/alignes`.
