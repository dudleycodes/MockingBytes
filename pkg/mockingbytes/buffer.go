package mockingbytes

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// buffer is an in-memory byte buffer that implements the io.ReadWriteCloser interface
type buffer struct {
	bytes.Buffer
	isClosed bool
	m        sync.Mutex
}

// Read reads the next len(p) bytes from the buffer or until the buffer is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero); otherwise it is nil.
func (b *buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()

	if b.Buffer.Len() < 1 {
		if b.isClosed {
			return 0, io.EOF
		}

		return 0, nil
	}

	n, err = b.Buffer.Read(p)

	if b.isClosed && b.Buffer.Len() == 0 {
		return n, io.EOF
	}

	return n, nil
}

// Write appends the contents of p to the buffer, growing the buffer as needed. The return value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with ErrTooLarge.
func (b *buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()

	if b.isClosed {
		return 0, errors.New("Cannot write to closed stream")
	}

	n, err = b.Buffer.Write(p)
	return
}

// Close closes the buffer, rendering it unusable for I/O. Close will return an error if it has already been called.
func (b *buffer) Close() error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.isClosed {
		return errors.New("already closed")
	}

	b.isClosed = true
	return nil
}
