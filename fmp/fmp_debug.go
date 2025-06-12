package fmp

import (
	"fmt"
	"os"
	"strings"
)

func (f *FmpFile) ToDebugFile(fname string) {
	f_chunks, err := os.Create(fname + ".chunks")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f_chunks.Close(); err != nil {
			panic(err)
		}
	}()
	for _, chunk := range f.Chunks {
		fmt.Fprintf(f_chunks, "%s, %s\n", chunk.String(), string(chunk.Value))
	}

	f_dicts, err := os.Create(fname + ".dicts")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f_dicts.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Fprint(f_dicts, f.Dictionary.String())
}

func (c *FmpChunk) String() string {
	return fmt.Sprintf("<%v(%v)>", c.Type, c.Length)
}

func (dict *FmpDict) String() string {
	s := ""
	for k, v := range *dict {
		ns := strings.ReplaceAll(v.Children.String(), "\n", "\n\t")
		s += fmt.Sprintf("%v: %v\n%v\n", k, string(v.Value), ns)
	}
	return s
}
