package mockingbytes

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRandomReader(t *testing.T) {
	//t.Parallel()

	//	tests := []int{-42, 0, 1, 7, 8, 42, 256, 512, 1024, 1048576}
	tests := []int{-42, 0, 1, 7}

	for _, test := range tests {
		t.Run(string(test), func(t *testing.T) {
			r := randomReader(test)
			b, err := ioutil.ReadAll(r)

			fmt.Println("BBBB", string(b), "BBBB")

			if err != nil {
				t.Fatalf("Created a random reader of size %d got an while reading error: %q", test, err.Error())
			}

			if test < 1 {
				if len(b) != 0 {
					t.Errorf("Created a random reader of size %d, expected to read 0 bytes, but read %d bytes", test, len(b))
				}
			} else if len(b) != test {
				t.Errorf("Created a random reader of size %d but read %d from it: %v", test, len(b), b)
			}
		})
	}
}
