package mockingbytes

import (
	"io"
	"io/ioutil"
	"math"
	"testing"
)

func TestRandomReader(t *testing.T) {
	//t.Parallel()

	tests := []int{-42, 0, 1, 7, 8, 42, 256, 512, 1024, 1048576, math.MaxInt32}

	for _, test := range tests {
		t.Run(string(test), func(t *testing.T) {
			actual := []byte{}
			var err error
			r := randomReader(test)

			for {
				b, err := ioutil.ReadAll(r)
				actual = append(actual, b...)

				if err != nil {
					break
				}
			}

			if err != io.EOF {
				t.Fatalf("Created a random reader of size %d got an unexpected error reading: %q", test, err.Error())
			}

			if test < 1 {
				if len(actual) != 0 {
					t.Errorf("Created a random reader of size %d, expected to read 0 bytes, but read %d bytes", test, len(actual))
				}
			} else if len(actual) != test {
				t.Errorf("Created a random reader of size %d but read %d from it", test, len(actual))
			}
		})
	}
}
