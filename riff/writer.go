// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riff

import (
	"io"
)

func NewWriter(formType FourCC) *Writer {
	return &Writer{chunk{FourCC{'R', 'I', 'F', 'F'}, 4, formType[:]}}
}

type Writer struct {
	c chunk
}

func (w *Writer) AppendChunk(id FourCC, data []byte) {
	w.c.len += appendChunk(&w.c.data, id, data)
}

func (w *Writer) AppendList(l *List) {
	w.c.len += appendList(&w.c.data, l)
}

func (w *Writer) WriteTo(s io.Writer) error {
	n, err := s.Write(w.c.bytes())
	if uint32(n) != 8+w.c.len {
		err = io.ErrShortWrite
	}
	return err
}

func NewList(listType FourCC) *List {
	return &List{chunk{FourCC{'L', 'I', 'S', 'T'}, 4, listType[:]}}
}

type List struct {
	c chunk
}

func (l *List) AppendChunk(id FourCC, data []byte) {
	l.c.len += appendChunk(&l.c.data, id, data)
}

func (l *List) AppendList(sub *List) {
	l.c.len += appendList(&l.c.data, sub)
}

type chunk struct {
	id   FourCC
	len  uint32
	data []byte
}

func (c chunk) bytes() []byte {
	b := make([]byte, 0, 8+c.len)
	b = []byte{
		c.id[0],
		c.id[1],
		c.id[2],
		c.id[3],
		byte(c.len >> 24),
		byte(c.len >> 16),
		byte(c.len >> 8),
		byte(c.len),
	}
	return append(b, c.data...)
}

func appendChunk(to *[]byte, id FourCC, data []byte) uint32 {
	len := uint32(len(data))
	*to = append(*to, chunk{id, len, data}.bytes()...)
	if len&1 == 1 {
		*to = append(*to, byte(0))
		len++
	}
	return 8 + len
}

func appendList(to *[]byte, l *List) uint32 {
	*to = append(*to, l.c.bytes()...)
	return 8 + l.c.len
}
