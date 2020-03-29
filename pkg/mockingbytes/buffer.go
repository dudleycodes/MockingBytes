package mockingbytes

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// buffer is a thread-safe implementation of bytes.buffer that includes support for the io.Closer interface
type buffer struct {
	buf      *bytes.Buffer
	isClosed bool
	mu       sync.Mutex
}

// newBuffer creates a new thread-safe read, write, closer
func newBuffer(buf []byte) *buffer {
	buf = bytes.Trim(buf, "\x00")

	return &buffer{
		buf: bytes.NewBuffer(buf),
	}
}

// Read reads the next len(p) bytes from the buffer or until the buffer is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero); otherwise it is nil.
func (b *buffer) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.buf.Len() < 1 {
		if b.buf.Cap() > 0 {
			b.buf.Reset()
		}

		if b.isClosed {
			return 0, io.EOF
		}

		return 0, nil
	}

	n, err = b.buf.Read(p)

	if n == 0 && len(p) == 0 {
		return 0, errors.New("unable to read contents and write to a buffer with a length of 0-bytes")
	}

	if b.isClosed && b.buf.Len() == 0 {
		b.buf.Reset()
		return n, io.EOF
	}

	return n, nil
}

// Write appends the contents of p to the buffer, growing the buffer as needed. The return value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with ErrTooLarge.
func (b *buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isClosed {
		return 0, errors.New("cannot write to the closed buffer")
	}

	n, err = b.buf.Write(p)

	return n, err
}

// Close the buffer for writing while allowing the reading any unread content. Will return a predictable error if already closed.
func (b *buffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isClosed {
		return errors.New("already closed")
	}

	b.isClosed = true

	return nil
}
