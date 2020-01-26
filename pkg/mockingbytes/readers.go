package mockingbytes

import (
	"errors"
	"io"
	"math/rand"
)

type ReaderOpt func(*rCfg) error

type rCfg struct {
	chunkMax     int
	chunkMin     int
	getChunkSize func() int
}

func defaultReaderCfg() *rCfg {
	return &rCfg{
		chunkMax:     8,
		chunkMin:     8,
		getChunkSize: func() int { return 8 },
	}
}

func SetChunkJitter(min, max uint16) func(*rCfg) error {
	return func(c *rCfg) error {
		if min < 1 {
			return errors.New("Can not have a min value less than 1")
		}

		if max < 1 {
			return errors.New("Can not have a max value less than 1")
		}

		if min > max {
			t := min
			min = max
			max = t
		}

		c.chunkMin, c.chunkMax = int(min), int(max)

		c.getChunkSize = func() int {
			return rand.Intn(c.chunkMax-c.chunkMin) + c.chunkMin
		}

		return nil
	}
}

func RandomReader(size int, options ...ReaderOpt) (io.Reader, error) {
	cfg, rwc := defaultReaderCfg(), &buffer{}

	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil, err
		}
	}

	var chunkSize int
	go func(wc io.WriteCloser, s int) {
		b := make([]byte, cfg.chunkMax)
		defer wc.Close()

		chunkSize = cfg.getChunkSize()

		for s >= chunkSize {
			rand.Read(b)
			n, _ := wc.Write(b[0:chunkSize])
			s -= n
		}

		if s > cfg.chunkMin {
			rand.Read(b)
			n, _ := wc.Write(b[0:cfg.chunkMin])
			s -= n
		}

		if s > 0 {
			rand.Read(b)
			wc.Write(b[0:s])
		}
	}(rwc, size)

	return rwc, nil
}

// func lagReader(reader io.Reader, tickDelay time.Duration) io.Reader {
// 	var rwc io.ReadWriteCloser

// 	c := make(chan []byte, 16)

// 	go func(r io.Reader, o chan<- []byte) {
// 		defer close(o)

// 		for {
// 			b := make([]byte, rand.Intn(3)+3)
// 			_, err := r.Read(b)

// 			if len(b) > 0 {
// 				o <- b
// 			}

// 			if err == io.EOF {
// 				break
// 			} else if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}(reader, c)

// 	go func(wc io.WriteCloser, i <-chan []byte, d time.Duration) {
// 		defer wc.Close()

// 		for {
// 			start := time.Now()

// 			b, more := <-i

// 			if _, err := wc.Write(b); err != nil {
// 				panic(err)
// 			}

// 			if w := time.Since(start); w < d {
// 				time.Sleep(d - w)
// 			}

// 			if !more {
// 				break
// 			}
// 		}
// 	}(rwc, c, tickDelay)

// 	return rwc
// }

// func repeaterStream(header byte, headerSize int, body byte, bodySize int, totalTime time.Duration) io.Reader {
// 	rand.Seed(time.Now().UnixNano())

// 	if headerSize < 0 {
// 		headerSize = 0
// 	}

// 	if bodySize < 0 {
// 		bodySize = 0
// 	}

// 	if totalTime < 0 {
// 		totalTime = 0
// 	}

// 	var rwc io.ReadWriteCloser

// 	go func(wc io.WriteCloser) {
// 		defer wc.Close()

// 		totalBytes := headerSize + bodySize
// 		if totalBytes < 1 {
// 			time.Sleep(totalTime)
// 			return
// 		}

// 		delayTick := time.Duration(int64(totalTime) / int64(headerSize+bodySize))
// 		buf := new(bytes.Buffer)

// 		// Done in a single loop so that header and body contents can blend together in a single write
// 		for bytesRemaining := totalBytes; bytesRemaining > 0; {
// 			bytesToSend := rand.Intn(3) + 1
// 			if bytesToSend > bytesRemaining {
// 				bytesToSend = bytesRemaining
// 			}
// 			buf.Grow(bytesToSend)

// 			if headerRemaining := bytesRemaining - bodySize; headerRemaining > 0 {
// 				headerToSend := bytesToSend
// 				if headerToSend > headerRemaining {
// 					headerToSend = headerRemaining
// 				}

// 				buf.Write(bytes.Repeat([]byte{header}, headerToSend))
// 				bytesToSend -= buf.Len()
// 			}

// 			if bytesToSend > 0 {
// 				if bodyRemaining := bytesRemaining - headerSize; bodyRemaining > 0 {
// 					bodyToSend := bytesToSend
// 					if bodyToSend > bodyRemaining {
// 						bodyToSend = bodyRemaining
// 					}

// 					buf.Write(bytes.Repeat([]byte{body}, bodyToSend))
// 					bytesToSend = 0
// 				}
// 			}

// 			n, _ := io.Copy(rwc, buf)
// 			buf.Reset()

// 			for i := 0; i < int(n); i++ {
// 				time.Sleep(delayTick)
// 			}

// 			bytesRemaining -= int(n)
// 		}
// 	}(rwc)

// 	return rwc
// }
