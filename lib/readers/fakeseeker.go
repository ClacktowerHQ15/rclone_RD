package readers

import (
	"errors"
	"fmt"
	"io"
)

// FakeSeeker adapts an io.Seeker into an io.ReadReader
type FakeSeeker struct {
	io.Reader
	length int64
	offset int64
	read   bool
}

// NewFakeSeeker creates a fake seeker from an io.Reader
//
// This can be seeked before reading to discover the length passed in.
func NewFakeSeeker(in io.Reader, length int64) *FakeSeeker {
	return &FakeSeeker{
		Reader: in,
		length: length,
	}
}

// Seek the stream - possible only before reading
func (r *FakeSeeker) Seek(offset int64, whence int) (abs int64, err error) {
	if r.read {
		return 0, fmt.Errorf("FakeSeeker: can't Seek(%d, %d) after reading", offset, whence)
	}
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.offset + offset
	case io.SeekEnd:
		abs = r.length + offset
	default:
		return 0, errors.New("FakeSeeker: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("FakeSeeker: negative position")
	}
	r.offset = abs
	return abs, nil
}

// Read data from the stream. Will give an error if seeked.
func (r *FakeSeeker) Read(p []byte) (n int, err error) {
	if !r.read && r.offset != 0 {
		return 0, errors.New("FakeSeeker: not at start: can't read")
	}
	n, err = r.Reader.Read(p)
	if n != 0 {
		r.read = true
	}
	return n, err
}
