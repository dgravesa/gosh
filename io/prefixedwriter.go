package io

import "io"

// PrefixedWriter is a Writer that prefixes its outputs
type PrefixedWriter struct {
	receiver io.Writer
	prefix   []byte
}

// NewPrefixedWriter creates a new PrefixedWriter, which delegates to w
func NewPrefixedWriter(w io.Writer) *PrefixedWriter {
	return &PrefixedWriter{
		receiver: w,
	}
}

// Prefix gets the current prefix in w
func (w *PrefixedWriter) Prefix() string {
	return string(w.prefix)
}

// SetPrefix sets the prefix of w to be p
func (w *PrefixedWriter) SetPrefix(p string) {
	w.prefix = []byte(p)
}

func (w *PrefixedWriter) Write(p []byte) (n int, err error) {
	return w.receiver.Write([]byte(append(w.prefix, p...)))
}
