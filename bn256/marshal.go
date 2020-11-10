// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gurvy DO NOT EDIT

package bn256

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"

	"github.com/consensys/gurvy/bn256/fp"
	"github.com/consensys/gurvy/bn256/fr"
)

// To encode G1 and G2 points, we mask the most significant bits with these bits to specify without ambiguity
// metadata needed for point (de)compression
// we have less than 3 bits available on the msw, so we can't follow BLS381 style encoding.
// the difference is the case where a point is infinity and uncompressed is not flagged
const (
	mMask               uint64 = 0b11 << 62
	mUncompressed       uint64 = 0b00 << 62
	mCompressedSmallest uint64 = 0b10 << 62
	mCompressedLargest  uint64 = 0b11 << 62
	mCompressedInfinity uint64 = 0b01 << 62
)

// Encoder writes bn256 object values to an output stream
type Encoder struct {
	w   io.Writer
	n   int64 // written bytes
	raw bool  // raw vs compressed encoding
}

// Decoder reads bn256 object values from an inbound stream
type Decoder struct {
	r io.Reader
	n int64 // read bytes
}

// NewDecoder returns a binary decoder supporting curve bn256 objects in both
// compressed and uncompressed (raw) forms
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the binary encoding of v from the stream
// type must be *uint64, *fr.Element, *fp.Element, *G1Affine, *G2Affine, *[]G1Affine or *[]G2Affine
func (dec *Decoder) Decode(v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || !rv.Elem().CanSet() {
		return errors.New("bn256 decoder: unsupported type, need pointer")
	}

	// implementation note: code is a bit verbose (abusing code generation), but minimize allocations on the heap
	// TODO double check memory usage and factorize this

	var buf [SizeOfG2Uncompressed]byte
	var read int
	var msw uint64

	switch t := v.(type) {
	case *uint64:
		msw, err = dec.readUint64()
		if err != nil {
			return
		}
		*t = msw
		return
	case *fr.Element:
		read, err = io.ReadFull(dec.r, buf[:fr.Limbs*8])
		dec.n += int64(read)
		if err != nil {
			return
		}
		t.SetBytes(buf[:fr.Limbs*8])
		return
	case *fp.Element:
		read, err = io.ReadFull(dec.r, buf[:fp.Limbs*8])
		dec.n += int64(read)
		if err != nil {
			return
		}
		t.SetBytes(buf[:fp.Limbs*8])
		return
	case *G1Affine:
		// read the most significant word
		read, err = io.ReadFull(dec.r, buf[:8])
		dec.n += int64(read)
		if err != nil {
			return
		}
		msw = binary.BigEndian.Uint64(buf[:8])
		nbBytes := SizeOfG1Uncompressed
		if isCompressed(msw) {
			nbBytes = SizeOfG1Compressed
		}
		read, err = io.ReadFull(dec.r, buf[8:nbBytes])
		dec.n += int64(read)
		if err != nil {
			return
		}
		_, err = t.SetBytes(buf[:nbBytes])
		return
	case *G2Affine:
		read, err = io.ReadFull(dec.r, buf[:8])
		dec.n += int64(read)
		if err != nil {
			return
		}
		msw = binary.BigEndian.Uint64(buf[:8])
		nbBytes := SizeOfG2Uncompressed
		if isCompressed(msw) {
			nbBytes = SizeOfG2Compressed
		}
		read, err = io.ReadFull(dec.r, buf[8:nbBytes])
		dec.n += int64(read)
		if err != nil {
			return
		}
		_, err = t.SetBytes(buf[:nbBytes])
		return
	case *[]G1Affine:
		msw, err = dec.readUint64()
		if err != nil {
			return
		}
		if len(*t) != int(msw) {
			*t = make([]G1Affine, msw)
		}
		for i := 0; i < len(*t); i++ {
			// read the most significant word
			read, err = io.ReadFull(dec.r, buf[:8])
			dec.n += int64(read)
			if err != nil {
				return
			}
			msw = binary.BigEndian.Uint64(buf[:8])
			nbBytes := SizeOfG1Uncompressed
			if isCompressed(msw) {
				nbBytes = SizeOfG1Compressed
			}
			read, err = io.ReadFull(dec.r, buf[8:nbBytes])
			dec.n += int64(read)
			if err != nil {
				return
			}
			_, err = (*t)[i].SetBytes(buf[:nbBytes])
			if err != nil {
				return
			}
		}
		return nil
	case *[]G2Affine:
		msw, err = dec.readUint64()
		if err != nil {
			return
		}
		if len(*t) != int(msw) {
			*t = make([]G2Affine, msw)
		}
		for i := 0; i < len(*t); i++ {
			read, err = io.ReadFull(dec.r, buf[:8])
			dec.n += int64(read)
			if err != nil {
				return
			}
			msw = binary.BigEndian.Uint64(buf[:8])
			nbBytes := SizeOfG2Uncompressed
			if isCompressed(msw) {
				nbBytes = SizeOfG2Compressed
			}
			read, err = io.ReadFull(dec.r, buf[8:nbBytes])
			dec.n += int64(read)
			if err != nil {
				return
			}
			_, err = (*t)[i].SetBytes(buf[:nbBytes])
			if err != nil {
				return
			}
		}
		return nil
	default:
		return errors.New("bn256 encoder: unsupported type")
	}
}

// BytesRead return total bytes read from reader
func (dec *Decoder) BytesRead() int64 {
	return dec.n
}

func (dec *Decoder) readUint64() (r uint64, err error) {
	var read int
	var buf [8]byte
	read, err = io.ReadFull(dec.r, buf[:8])
	dec.n += int64(read)
	if err != nil {
		return
	}
	r = binary.BigEndian.Uint64(buf[:8])
	return
}

func isCompressed(msw uint64) bool {
	mData := msw & mMask
	return !(mData == mUncompressed)
}

// NewEncoder returns a binary encoder supporting curve bn256 objects
func NewEncoder(w io.Writer, options ...func(*Encoder)) *Encoder {
	// default settings
	enc := &Encoder{
		w:   w,
		n:   0,
		raw: false,
	}

	// handle options
	for _, option := range options {
		option(enc)
	}

	return enc
}

// Encode writes the binary encoding of v to the stream
// type must be uint64, *fr.Element, *fp.Element, *G1Affine, *G2Affine, []G1Affine or []G2Affine
func (enc *Encoder) Encode(v interface{}) (err error) {
	if enc.raw {
		return enc.encodeRaw(v)
	}
	return enc.encode(v)
}

// BytesWritten return total bytes written on writer
func (enc *Encoder) BytesWritten() int64 {
	return enc.n
}

// RawEncoding returns an option to use in NewEncoder(...) which sets raw encoding mode to true
// points will not be compressed using this option
func RawEncoding() func(*Encoder) {
	return func(enc *Encoder) {
		enc.raw = true
	}
}

func (enc *Encoder) encode(v interface{}) (err error) {

	// implementation note: code is a bit verbose (abusing code generation), but minimize allocations on the heap
	// TODO double check memory usage and factorize this

	var written int
	switch t := v.(type) {
	case uint64:
		err = binary.Write(enc.w, binary.BigEndian, t)
		enc.n += 8
		return
	case *fr.Element:
		buf := t.Bytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case *fp.Element:
		buf := t.Bytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case *G1Affine:
		buf := t.Bytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case *G2Affine:
		buf := t.Bytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case []G1Affine:
		// write slice length
		err = binary.Write(enc.w, binary.BigEndian, uint64(len(t)))
		if err != nil {
			return
		}
		enc.n += 8

		var buf [SizeOfG1Compressed]byte

		for i := 0; i < len(t); i++ {
			buf = t[i].Bytes()
			written, err = enc.w.Write(buf[:])
			enc.n += int64(written)
			if err != nil {
				return
			}
		}
		return nil
	case []G2Affine:
		// write slice length
		err = binary.Write(enc.w, binary.BigEndian, uint64(len(t)))
		if err != nil {
			return
		}
		enc.n += 8

		var buf [SizeOfG2Compressed]byte

		for i := 0; i < len(t); i++ {
			buf = t[i].Bytes()
			written, err = enc.w.Write(buf[:])
			enc.n += int64(written)
			if err != nil {
				return
			}
		}
		return nil
	default:
		return errors.New("<no value> encoder: unsupported type")
	}
}

func (enc *Encoder) encodeRaw(v interface{}) (err error) {

	// implementation note: code is a bit verbose (abusing code generation), but minimize allocations on the heap
	// TODO double check memory usage and factorize this

	var written int
	switch t := v.(type) {
	case uint64:
		err = binary.Write(enc.w, binary.BigEndian, t)
		enc.n += 8
		return
	case *fr.Element:
		buf := t.Bytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case *fp.Element:
		buf := t.Bytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case *G1Affine:
		buf := t.RawBytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case *G2Affine:
		buf := t.RawBytes()
		written, err = enc.w.Write(buf[:])
		enc.n += int64(written)
		return
	case []G1Affine:
		// write slice length
		err = binary.Write(enc.w, binary.BigEndian, uint64(len(t)))
		if err != nil {
			return
		}
		enc.n += 8

		var buf [SizeOfG1Uncompressed]byte

		for i := 0; i < len(t); i++ {
			buf = t[i].RawBytes()
			written, err = enc.w.Write(buf[:])
			enc.n += int64(written)
			if err != nil {
				return
			}
		}
		return nil
	case []G2Affine:
		// write slice length
		err = binary.Write(enc.w, binary.BigEndian, uint64(len(t)))
		if err != nil {
			return
		}
		enc.n += 8

		var buf [SizeOfG2Uncompressed]byte

		for i := 0; i < len(t); i++ {
			buf = t[i].RawBytes()
			written, err = enc.w.Write(buf[:])
			enc.n += int64(written)
			if err != nil {
				return
			}
		}
		return nil
	default:
		return errors.New("<no value> encoder: unsupported type")
	}
}
