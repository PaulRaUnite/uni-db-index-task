package limit

import (
	"io"
)

type LineLimitReader struct {
	r     io.Reader
	count int
	limit int
}

func NewLineLimitReader(r io.Reader, limit int) *LineLimitReader {
	return &LineLimitReader{
		r:     r,
		count: 0,
		limit: limit,
	}
}
func (clr *LineLimitReader) Read(p []byte) (int, error) {
	if clr.count >= clr.limit {
		return 0, io.EOF
	}
	n, err := clr.r.Read(p)
	if n > 0 {
		for _, char := range p {
			if char == '\n' {
				clr.count++
			}
		}
	}
	return n, err
}
