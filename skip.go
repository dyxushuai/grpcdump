//    Copyright 2018 <xu shuai <dyxushuai@gmail.com>>
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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
