package mockingbytes

import (
	"io"
	"math/rand"
	"testing"
	"time"
)

func Benchmark_Buffer(b *testing.B) {
	rawBuffer := make([]byte, 64)
	sut := newBuffer(rawBuffer)

	readBuffer := make([]byte, 9)

	content := make([]byte, 512)
	if _, err := rand.Read(content); err != nil {
		b.Fatal("Was unable to generate random byte slice")
	}

	b.Run("write and read in 16 byte chunks", func(b *testing.B) {
		b.SetBytes(16)
		for n := 0; n < b.N; n++ {
			if _, err := sut.Write(content[0:15]); err != nil {
				b.Error(err)
			}

			if _, err := sut.Read(readBuffer); err == io.EOF {
				b.Error(io.ErrUnexpectedEOF)
			} else if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("write and read in 128 byte chunks", func(b *testing.B) {
		b.SetBytes(128)
		for n := 0; n < b.N; n++ {
			if _, err := sut.Write(content[0:127]); err != nil {
				b.Error(err)
			}

			if _, err := sut.Read(readBuffer); err == io.EOF {
				b.Error(io.ErrUnexpectedEOF)
			} else if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("write and read in 256 byte chunks", func(b *testing.B) {
		b.SetBytes(256)
		for n := 0; n < b.N; n++ {
			if _, err := sut.Write(content[0:255]); err != nil {
				b.Error(err)
			}

			if _, err := sut.Read(readBuffer); err == io.EOF {
				b.Error(io.ErrUnexpectedEOF)
			} else if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("1 to 512 byte writes and reads", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			size := rand.Intn(511) + 1
			b.SetBytes(int64(size))

			if _, err := sut.Write(content[0:size]); err != nil {
				b.Error(err)
			}

			if _, err := sut.Read(readBuffer); err == io.EOF {
				b.Error(io.ErrUnexpectedEOF)
			} else if err != nil {
				b.Error(err)
			}
		}
	})
	sut.Close()
}

func Test_Buffer_ThreadSafety(t *testing.T) {
	t.Parallel()

	rawBuffer := make([]byte, 64)
	sut, read, written := newBuffer(rawBuffer), make(chan int), make(chan int)

	go func(wc io.WriteCloser, c chan<- int) {
		defer wc.Close()
		start, i := time.Now(), 0

		for time.Since(start) < time.Second {
			wc.Write([]byte("a"))
			time.Sleep(time.Millisecond)
			wc.Write([]byte("b"))
			time.Sleep(time.Millisecond)
			wc.Write([]byte("c"))
			i += 3
		}

		c <- i
		close(c)
	}(sut, read)

	go func(r io.Reader, c chan<- int) {
		b, i := make([]byte, 3), 0
		time.Sleep(time.Microsecond * 5)

		for {
			n, err := r.Read(b)
			if n > 0 {
				i += n
				time.Sleep(time.Millisecond)
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				t.Errorf("Received non EOF error while reading %s", err.Error())
				break
			}
		}

		c <- i
		close(c)
	}(sut, written)

	actualRead, actualWritten := <-read, <-written

	if actualRead != actualWritten {
		t.Errorf("Read and writes were mismatched. Wrote %d but Read %d", actualRead, actualWritten)
	}
}
