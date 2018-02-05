package audio

import "io"

// File is an audio file
type File struct {
	// implement this
}

func NewFromBytes(b []byte) *File {
	// implement this
	return nil
}

func NewFromReader(r io.Reader) *File {
	// implement this
	return nil
}

func (f *File) WAVData() []byte {
	// implement this
	return nil
}
