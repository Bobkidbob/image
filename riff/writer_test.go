// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riff

import (
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	w := NewWriter(FourCC{'T', 'E', 'S', 'T'})
	w.AppendChunk(FourCC{'e', 'v', 'e', 'n'}, []byte{0x00, 0x11, 0x22, 0x33})
	w.AppendChunk(FourCC{'o', 'd', 'd', ' '}, []byte{0x44, 0x55, 0x66})

	l := NewList(FourCC{'l', 's', 't', ' '})
	w.AppendList(l)

	var buf bytes.Buffer
	if err := w.WriteTo(&buf); err != nil {
		t.Errorf("WriteTo failed: %v", err)
	}

	b := buf.Bytes()
	sliceAssertion(t, "RIFF ID", b[:4], []byte{'R', 'I', 'F', 'F'})
	sliceAssertion(t, "RIFF size", b[4:8], []byte{0, 0, 0, 40})
	sliceAssertion(t, "RIFF form type", b[8:12], []byte{'T', 'E', 'S', 'T'})
	sliceAssertion(t, "Even chunk ID", b[12:16], []byte{'e', 'v', 'e', 'n'})
	sliceAssertion(t, "Even chunk size", b[16:20], []byte{0, 0, 0, 4})
	sliceAssertion(t, "Even chunk data", b[20:24], []byte{0x00, 0x11, 0x22, 0x33})
	sliceAssertion(t, "Odd chunk ID", b[24:28], []byte{'o', 'd', 'd', ' '})
	sliceAssertion(t, "Odd chunk size", b[28:32], []byte{0, 0, 0, 3})
	sliceAssertion(t, "Odd chunk data", b[32:36], []byte{0x44, 0x55, 0x66, 0x00})
	sliceAssertion(t, "List chunk ID", b[36:40], []byte{'L', 'I', 'S', 'T'})
	sliceAssertion(t, "List chunk size", b[40:44], []byte{0, 0, 0, 4})
	sliceAssertion(t, "List chunk type", b[44:], []byte{'l', 's', 't', ' '})
}

func TestList(t *testing.T) {
	l := NewList(FourCC{'T', 'E', 'S', 'T'})

	l1 := NewList(FourCC{'S', 'U', 'B', '1'})
	l1.AppendChunk(FourCC{'c', 'h', 'n', 'k'}, []byte{1, 2, 3, 4})
	l1.AppendChunk(FourCC{'c', 'h', 'n', 'k'}, []byte{1, 2, 3})
	l.AppendList(l1)

	l2 := NewList(FourCC{'S', 'U', 'B', '2'})
	l2.AppendChunk(FourCC{'c', 'h', 'n', 'k'}, []byte{1, 2, 3})
	l.AppendList(l2)

	b := l.c.bytes()
	sliceAssertion(t, "List ID", b[:4], []byte{'L', 'I', 'S', 'T'})
	sliceAssertion(t, "List size", b[4:8], []byte{0, 0, 0, 64})
	sliceAssertion(t, "List type", b[8:12], []byte{'T', 'E', 'S', 'T'})
	sliceAssertion(t, "Sublist 1 ID", b[12:16], []byte{'L', 'I', 'S', 'T'})
	sliceAssertion(t, "Sublist 1 size", b[16:20], []byte{0, 0, 0, 28})
	sliceAssertion(t, "Sublist 1 type", b[20:24], []byte{'S', 'U', 'B', '1'})
	sliceAssertion(t, "Sublist 1 chunk 1 ID", b[24:28], []byte{'c', 'h', 'n', 'k'})
	sliceAssertion(t, "Sublist 1 chunk 1 size", b[28:32], []byte{0, 0, 0, 4})
	sliceAssertion(t, "Sublist 1 chunk 1 data", b[32:36], []byte{1, 2, 3, 4})
	sliceAssertion(t, "Sublist 1 chunk 2 ID", b[36:40], []byte{'c', 'h', 'n', 'k'})
	sliceAssertion(t, "Sublist 1 chunk 2 size", b[40:44], []byte{0, 0, 0, 3})
	sliceAssertion(t, "Sublist 1 chunk 2 data", b[44:48], []byte{1, 2, 3, 0})
	sliceAssertion(t, "Sublist 2 ID", b[48:52], []byte{'L', 'I', 'S', 'T'})
	sliceAssertion(t, "Sublist 2 size", b[52:56], []byte{0, 0, 0, 16})
	sliceAssertion(t, "Sublist 2 type", b[56:60], []byte{'S', 'U', 'B', '2'})
	sliceAssertion(t, "Sublist 2 chunk 1 ID", b[60:64], []byte{'c', 'h', 'n', 'k'})
	sliceAssertion(t, "Sublist 2 chunk 1 size", b[64:68], []byte{0, 0, 0, 3})
	sliceAssertion(t, "Sublist 2 chunk 1 data", b[68:], []byte{1, 2, 3, 0})
}

func sliceAssertion(t *testing.T, name string, got, want []byte) {
	if !bytes.Equal(got, want) {
		t.Errorf("%v: got %v, want %v", name, got, want)
	}
}
