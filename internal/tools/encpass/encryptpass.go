package encpass

import "golang.org/x/crypto/bcrypt"

// EncPassword - returns hash pwd.
func EncPassword(pwd string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

// ComparePassword - compare password and hash.
func ComparePassword(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
