package grpcdump

import (
	"bufio"
	"io"

	"golang.org/x/net/http2"
)

type skipClientPrefaceReader struct {
	buf *bufio.Reader
}

func newSkipClientPrefaceReader(r io.Reader) *skipClientPrefaceReader {
	return &skipClientPrefaceReader{
		buf: bufio.NewReader(r),
	}
}

func (s *skipClientPrefaceReader) skip() {
	readed := 0
	for {
		b, err := s.buf.ReadByte()
		if err != nil {
			break
		}
		readed++
		if b != http2.ClientPreface[readed-1] {
			break
		}
		if readed == len(http2.ClientPreface) {
			return
		}
	}

	if readed != len(http2.ClientPreface) && readed != 0 {
		for i := 0; i < readed; i++ {
			s.buf.UnreadByte()
		}
	}
}

func (s *skipClientPrefaceReader) Read(p []byte) (int, error) {
	s.skip()
	return s.buf.Read(p)
}
