package fmp

import (
	"fmt"
	"os"
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

func (dict *FmpDict) string(parentPath string) string {
	s := ""
	for k, v := range *dict {
		s += fmt.Sprintf("%v%v: %v\n", parentPath, k, string(v.Value))

		if v.Children != nil {
			s += v.Children.string(fmt.Sprintf("%v%v.", parentPath, k))
		}
	}
	return s
}

func (dict *FmpDict) String() string {
	return dict.string("")
}
