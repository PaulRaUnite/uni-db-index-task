package splitter

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"os"
	"path"
	"unicode/utf8"

	"github.com/pkg/errors"
)

func Run(splitFrom io.Reader, fromPath, splitInto string, lineCount int) error {
	boundary := int(math.Round(float64(lineCount) * 0.8))
	dir := path.Dir(fromPath)
	file80, file20, err := openFiles(dir, splitInto)
	if err != nil {
		return err
	}
	if file80 != nil {
		defer file80.Close()
	}
	if file20 != nil {
		defer file20.Close()
	}
	file80W := bufio.NewWriter(file80)
	file20W := bufio.NewWriter(file20)

	s := bufio.NewScanner(splitFrom)
	if !s.Scan() {
		return s.Err()
	}
	err = writeHeader(s.Bytes(), file80W, file20W)
	if err != nil {
		return errors.Wrap(err, "failed to write header")
	}

	for i := 0; s.Scan() && i < boundary; i++ {
		err = writeLine(file80W, s.Bytes())
		if err != nil {
			return err
		}
	}
	for i := 0; s.Scan(); i++ {
		err = writeLine(file20W, s.Bytes())
		if err != nil {
			return err
		}
	}
	err = s.Err()
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "failed to scan input file")
	}
	return nil
}

func writeHeader(header []byte, file80W *bufio.Writer, file20W *bufio.Writer) error {
	_, err := file80W.Write(header)
	if err != nil {
		return errors.Wrap(err, "failed to write header into 80 percent data file")
	}
	_, err = file20W.Write(header)
	if err != nil {
		return errors.Wrap(err, "failed to write header into 20 percent data file")
	}
	return nil
}

var check = []byte(",check,")
var comma = []byte(",")
var spaceComma = []byte(" ,")

func writeLine(file *bufio.Writer, line []byte) error {

	if !(48 <= line[0] && line[0] <= 57) ||
		bytes.Contains(line, check) ||
		!utf8.Valid(line) ||
		bytes.Count(line, comma) != 7 {
		return nil
	}
	line = bytes.ReplaceAll(line, spaceComma, comma)

	_, err := file.WriteString("\n")
	if err != nil {
		return errors.Wrap(err, "failed to write new line to file")
	}
	_, err = file.Write(line)
	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}
	return err
}

func openFiles(dir string, splitInto string) (*os.File, *os.File, error) {
	file80, err := os.Create(path.Join(dir, splitInto+"80.csv"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create file for 80 percent of data")
	}
	file20, err := os.Create(path.Join(dir, splitInto+"20.csv"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create file for 20 percent of data")
	}
	return file80, file20, nil
}
