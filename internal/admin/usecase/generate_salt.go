package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"crypto/rand"
)

const (
	letters    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	saltLength = 8
)

func GenerateRandomString(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func (a *AdminUseCaseImpl) generateHash(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))

	return hex.EncodeToString(hash.Sum(nil))
}
