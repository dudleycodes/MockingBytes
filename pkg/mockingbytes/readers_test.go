package mockingbytes

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRandomReader(t *testing.T) {
	t.Parallel()

	tests := []int{-42, 0, 1, 7, 8, 42, 256, 512, 1024, 1048576}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test), func(t *testing.T) {
			r, err := RandomReader(test, SetChunkJitter(2, 8))
			if err != nil {
				t.Fatalf("Error while creating reader: %s", err.Error())
			}

			actual, err := ioutil.ReadAll(r)

			if err != nil {
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
