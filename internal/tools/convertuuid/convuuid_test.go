package convertuuid

import "testing"

func TestUUIDString(t *testing.T) {
	uuid := UUID{
		0x6f, 0xdd, 0x89, 0xf3,
		0xe7, 0x40,
		0x46, 0x4a,
		0x96, 0xd5,
		0xc9, 0x4d, 0xa4, 0x0a, 0x3a, 0x12}

	expected := "6fdd89f3-e740-464a-96d5-c94da40a3a12"
	got := uuid.String()

	if got != expected {
		t.Errorf("UUID String() method failed: expected %v but got %v", expected, got)
	}
}

func TestEncodeHex(t *testing.T) {
	uuid := UUID{
		0x6f, 0xdd, 0x89, 0xf3,
		0xe7, 0x40,
		0x46, 0x4a,
		0x96, 0xd5,
		0xc9, 0x4d, 0xa4, 0x0a, 0x3a, 0x12}

	expected := "6fdd89f3-e740-464a-96d5-c94da40a3a12"
	var buf [36]byte
	encodeHex(buf[:], uuid)
	got := string(buf[:])

	if got != expected {
		t.Errorf("encodeHex() function failed: expected %v but got %v", expected, got)
	}
}
