package ber

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type byteReader interface {
	io.Reader
	ReadByte() (byte, error)
}

type TypeError struct {
	Tag  Tag
	Kind reflect.Kind
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("incompatible types %s and %s", t.Tag, t.Kind)
}

func parseInteger(v reflect.Value, b []byte) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return &TypeError{TagInteger, v.Kind()}
	}

	if len(b) > 8 {
		return errors.New("value does not fit in a int64")
	}

	var n int64
	for i, v := range b {
		shift := uint((len(b) - i - 1) * 8)
		if i == 0 {
			if v&0x80 != 0 {
				n -= 0x80 << shift
				v &= 0x7f
			}
		}
		n += int64(v) << shift
	}

	reflect.Value(v).SetInt(n)
	return nil
}

func Unmarshal(data []byte, v interface{}) error {

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("ber: value must be non nil pointer")
	}
	elem := rv.Elem()

	c, isPrimative, t, data, err := readNext(bytes.NewReader(data))
	if err != nil {
		return err
	}
	_, _ = c, isPrimative
	switch t {
	case TagInteger:
		return parseInteger(elem, data)
	default:
		return fmt.Errorf("unsupported tag: %s", t)
	}
}

type reader struct {
	*bufio.Reader
}

// the bytes stop being valid
func readNext(r byteReader) (c Class, isPrimative bool, t Tag, data []byte, err error) {
	var b byte
	if b, err = r.ReadByte(); err != nil {
		return
	}
	i := identifier(b)
	c = i.Class()
	t = i.Tag()

	// check primative or constructed (6th bit)
	isPrimative = i&0x20 == 0

	if b, err = r.ReadByte(); err != nil {
		return
	}

	// indefinate length
	if b == 0x80 {
		err = errors.New("ber: indefinate length encodings not supported")
		return
	}

	var n int
	if b&0x80 == 0 {
		// short form
		n = int(b)
	} else {
		// long form
		if n = int(b & 0x7f); n > 4 {
			err = errors.New("length octet: long form longer than 4 bytes")
		}
		length := make([]byte, 4)
		if _, err = io.ReadFull(r, length[4-n:]); err != nil {
			return
		}
		n = int(binary.BigEndian.Uint32(length))
	}

	data = make([]byte, n)
	_, err = io.ReadFull(r, data)
	return
}

//go:generate stringer -type=Class

type Class byte

const (
	ClassUniversal Class = iota
	ClassApplication
	ClassContextSpecific
	ClassPrivate
)

type identifier byte

func (i identifier) Class() Class {
	switch i >> 6 {
	case 0:
		return ClassUniversal
	case 1:
		return ClassApplication
	case 2:
		return ClassContextSpecific
	default:
		return ClassPrivate
	}
}

func (i identifier) Tag() Tag {
	return Tag(i & 0x1f)
}

//go:generate stringer -type=Tag

type Tag byte

const (
	TagEOC              Tag = 0x00
	TagBoolean          Tag = 0x01
	TagInteger          Tag = 0x02
	TagBitString        Tag = 0x03
	TagOctetString      Tag = 0x04
	TagNull             Tag = 0x05
	TagObjectIdentifier Tag = 0x06
	TagObjectDescriptor Tag = 0x07
	TagExternal         Tag = 0x08
	TagReal             Tag = 0x09
	TagEnumerated       Tag = 0x0a
	TagEmbeddedPDV      Tag = 0x0b
	TagUTF8String       Tag = 0x0c
	TagRelativeOID      Tag = 0x0d
	TagSequence         Tag = 0x10
	TagSet              Tag = 0x11
	TagNumericString    Tag = 0x12
	TagPrintableString  Tag = 0x13
	TagT61String        Tag = 0x14
	TagVideotexString   Tag = 0x15
	TagIA5String        Tag = 0x16
	TagUTCTime          Tag = 0x17
	TagGeneralizedTime  Tag = 0x18
	TagGraphicString    Tag = 0x19
	TagVisiableString   Tag = 0x1A
	TagGeneralString    Tag = 0x1b
	TagUniversalString  Tag = 0x1c
	TagCharacterString  Tag = 0x1d
	TagBMPString        Tag = 0x1e
	TagUseLongForm      Tag = 0x1f
)
