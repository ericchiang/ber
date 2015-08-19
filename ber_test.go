package ber

import "testing"

type Human struct {
	Name string
	Age  int
}

func TestUnmarshalInteger(t *testing.T) {
	tests := []struct {
		raw []byte
		exp int64
	}{
		{[]byte{0x02, 0x08, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 72057594037927935},
		{[]byte{0x02, 0x01, 0x05}, 5},
		{[]byte{0x02, 0x01, 0xfc}, -4},
	}
	for _, test := range tests {
		var i int64
		if err := Unmarshal(test.raw, &i); err != nil {
			t.Errorf("could not unmarshal: %d: %v", test.exp, err)
			continue
		}
		if test.exp != i {
			t.Errorf("expected %d got %d", test.exp, i)
		}
	}
}
