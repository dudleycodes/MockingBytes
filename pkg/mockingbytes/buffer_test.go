package mockingbytes

import (
	"io"
	"testing"
	"time"
)

func Benchmark_Buffer(b *testing.B) {
	sut := &buffer{}
	buf := make([]byte, 4)

	b.SetBytes(4)

	for n := 0; n < b.N; n++ {
		sut.Write([]byte("ping"))
		sut.Read(buf)
	}

	sut.Close()
}

func Test_Buffer_ThreadSafety(t *testing.T) {
	t.Parallel()

	sut, read, written := &buffer{}, make(chan int), make(chan int)

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
				t.Fatalf("Received non EOF error while reading %s", err.Error())
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
