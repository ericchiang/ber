package ber

import "testing"

type Human struct {
	Name string
	Age  int
}

func TestUnmarshal(t *testing.T) {
	input := []byte{0x02, 0x01, 0x05}
	var i int
	if err := Unmarshal(input, &i); err != nil {
		t.Fatal(err)
	}
	if i != 5 {
		t.Errorf("expected 5 got %d", i)
	}

	input = []byte{0x02, 0x01, 0xfc}
	if err := Unmarshal(input, &i); err != nil {
		t.Fatal(err)
	}
	if i != -4 {
		t.Errorf("expected -4 got %d", i)
	}
}

func TestType(t *testing.T) {
	if err := Unmarshal([]byte{}, 5); err == nil {
		t.Errorf("expected int to return an error, got nil")
	}
}
