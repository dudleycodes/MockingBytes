package mockingbytes

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// buffer is just here to make bytes.Buffer an io.ReadWriteCloser.
type buffer struct {
	bytes.Buffer
}

// Add a Close method to our buffer so that we satisfy io.ReadWriteCloser.
func (b *buffer) Close() error {
	b.Buffer.Reset()
	return nil
}

func lagReader(reader io.Reader, tickDelay time.Duration) io.Reader {
	var rwc io.ReadWriteCloser

	c := make(chan []byte, 16)

	go func(r io.Reader, o chan<- []byte) {
		defer close(o)

		for {
			b := make([]byte, rand.Intn(3)+3)
			_, err := r.Read(b)

			if len(b) > 0 {
				o <- b
			}

			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
		}
	}(reader, c)

	go func(wc io.WriteCloser, i <-chan []byte, d time.Duration) {
		defer wc.Close()

		for {
			start := time.Now()

			b, more := <-i

			if _, err := wc.Write(b); err != nil {
				panic(err)
			}

			if w := time.Since(start); w < d {
				time.Sleep(d - w)
			}

			if !more {
				break
			}
		}
	}(rwc, c, tickDelay)

	return rwc
}

func randomReader(size int) io.Reader {
	var rwc io.ReadWriteCloser
	rwc = &buffer{}

	if size < 1 {
		defer rwc.Close()
		return rwc
	}

	const chunkSize = 8

	go func(wc io.WriteCloser) {
		defer wc.Close()

		b := make([]byte, chunkSize)

		for i := 0; ; i++ {
			rand.Read(b)
			remaining := size - (chunkSize * i)

			if remaining >= chunkSize {
				wc.Write(b)
				continue
			}

			if remaining > 0 {
				fmt.Println(remaining)
				n, _ := rwc.Write(b)

				fmt.Println("AAA", remaining, n, "AAA")
				continue
			}

			break
		}
	}(rwc)

	return rwc
}

func repeaterStream(header byte, headerSize int, body byte, bodySize int, totalTime time.Duration) io.Reader {
	rand.Seed(time.Now().UnixNano())

	if headerSize < 0 {
		headerSize = 0
	}

	if bodySize < 0 {
		bodySize = 0
	}

	if totalTime < 0 {
		totalTime = 0
	}

	var rwc io.ReadWriteCloser

	go func(wc io.WriteCloser) {
		defer wc.Close()

		totalBytes := headerSize + bodySize
		if totalBytes < 1 {
			time.Sleep(totalTime)
			return
		}

		delayTick := time.Duration(int64(totalTime) / int64(headerSize+bodySize))
		buf := new(bytes.Buffer)

		// Done in a single loop so that header and body contents can blend together in a single write
		for bytesRemaining := totalBytes; bytesRemaining > 0; {
			bytesToSend := rand.Intn(3) + 1
			if bytesToSend > bytesRemaining {
				bytesToSend = bytesRemaining
			}
			buf.Grow(bytesToSend)

			if headerRemaining := bytesRemaining - bodySize; headerRemaining > 0 {
				headerToSend := bytesToSend
				if headerToSend > headerRemaining {
					headerToSend = headerRemaining
				}

				buf.Write(bytes.Repeat([]byte{header}, headerToSend))
				bytesToSend -= buf.Len()
			}

			if bytesToSend > 0 {
				if bodyRemaining := bytesRemaining - headerSize; bodyRemaining > 0 {
					bodyToSend := bytesToSend
					if bodyToSend > bodyRemaining {
						bodyToSend = bodyRemaining
					}

					buf.Write(bytes.Repeat([]byte{body}, bodyToSend))
					bytesToSend = 0
				}
			}

			n, _ := io.Copy(rwc, buf)
			buf.Reset()

			for i := 0; i < int(n); i++ {
				time.Sleep(delayTick)
			}

			bytesRemaining -= int(n)
		}
	}(rwc)

	return rwc
}
