package processor

import (
	"bufio"
	"encoding/base64"
	"io"
)

// Base64Processor is the Processor implementation for Base64 encoding
type Base64Processor struct {
	s    *bufio.Scanner
	w    io.Writer
	proc func(buf, src []byte) ([]byte, error)
}

// NewBase64EncodeProcessor creates and initializes a Base64Processor for encoding
func NewBase64EncodeProcessor(src io.Reader, dst io.Writer) *Base64Processor {
	return &Base64Processor{
		s: bufio.NewScanner(src),
		w: dst,
		proc: func(buf, src []byte) ([]byte, error) {
			n := base64.StdEncoding.EncodedLen(len(src))
			if cap(buf) < n {
				buf = make([]byte, n)
			}
			base64.StdEncoding.Encode(buf, src)
			return buf[:n], nil
		},
	}
}

// NewBase64DecodeProcessor creates and initializes a Base64Processor for decoding
func NewBase64DecodeProcessor(src io.Reader, dst io.Writer) *Base64Processor {
	return &Base64Processor{
		s: bufio.NewScanner(src),
		w: dst,
		proc: func(buf, src []byte) ([]byte, error) {
			n := base64.StdEncoding.DecodedLen(len(src))
			if cap(buf) < n {
				buf = make([]byte, n)
			}
			var err error
			n, err = base64.StdEncoding.Decode(buf, src)
			return buf[:n], err
		},
	}
}

// Process processes the stream and returns possible fatal error
func (b *Base64Processor) Process() error {
	var err error
	var buf []byte
	for b.s.Scan() {
		if buf, err = b.proc(buf, b.s.Bytes()); err != nil {
			return err
		}
		if _, err = b.w.Write(buf); err != nil {
			return err
		}
		if _, err = b.w.Write(LineBreakBytes); err != nil {
			return err
		}
	}
	if err = b.s.Err(); err != nil {
		return err
	}
	return nil
}
