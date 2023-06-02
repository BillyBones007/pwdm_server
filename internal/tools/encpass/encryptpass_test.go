package encpass

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestEncPassword(t *testing.T) {
	testCases := []struct {
		name, password string
	}{
		{"Case 1", "testpassword1"},
		{"Case 2", "testpassword2"},
		{"Case 3", "testpassword3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := EncPassword(tc.password)
			if err != nil {
				t.Errorf("EncPassword(%v) returned error: %v", tc.password, err)
			}

			if len(hash) == 0 {
				t.Errorf("EncPassword(%v) returned empty string", tc.password)
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "testpassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("Failed to generate hash from password: %v", err)
	}

	testCases := []struct {
		name, password, hash string
		expected             bool
	}{
		{"Case 1", password, string(hash), true},
		{"Case 2", "wrongpassword", string(hash), false},
		{"Case 3", password, "wronghash", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ComparePassword(tc.password, tc.hash)
			if result != tc.expected {
				t.Errorf("ComparePassword(%v, %v) = %v, expected %v", tc.password, tc.hash, result, tc.expected)
			}
		})
	}
}
