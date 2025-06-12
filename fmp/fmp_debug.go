package fmp

import (
	"fmt"
	"os"
)

func (f *FmpFile) ToDebugFile(fname string) {
	fo, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	for _, chunk := range f.Chunks {
		fmt.Fprintf(fo, "%s, %s\n", chunk.String(), string(chunk.Value))
	}
}

func (c *FmpChunk) String() string {
	return fmt.Sprintf("<%v(%v)>", c.Type, c.Length)
}
