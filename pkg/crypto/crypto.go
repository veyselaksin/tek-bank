package crypto

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

//go:generate mockgen -destination=../../mocks/crypto/crypto_mock.go -package=crypto tek-bank/pkg/crypto Crypto
type Crypto interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hashedPassword string) bool
	RandomNumber() int64
	RandomPassword() string
	RandomIBAN(isoCode string) string
	GenerateToken(length int) (string, error)
}

type crypto struct{}

func NewCrypto() Crypto {
	return &crypto{}
}

func (c *crypto) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	return string(bytes), err
}

func (c *crypto) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// RandomNumber generates a random customer number with 12 digits
func (c *crypto) RandomNumber() int64 {
	return rand.Int63n(1e12)
}

func (c *crypto) RandomPassword() string {
	digitRunes := []rune("0123456789")
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	d := make([]rune, 8)
	for i := range d {
		if i%2 == 0 {
			d[i] = letterRunes[rand.Intn(len(letterRunes))]
		} else {
			d[i] = digitRunes[rand.Intn(len(digitRunes))]
		}
	}

	return string(d)
}

func (c *crypto) RandomIBAN(isoCode string) string {
	digitRunes := []rune("0123456789")

	d := make([]rune, 28)
	for i := range d {
		d[i] = digitRunes[rand.Intn(len(digitRunes))]
	}

	return isoCode + "00 " + string(d[0:4]) + " " + string(d[4:8]) + " " + string(d[8:12]) + " " + string(d[12:16]) + " " + string(d[16:20]) + " " + string(d[20:24]) + " " + string(d[24:28])
}

func (c *crypto) GenerateToken(length int) (string, error) {
	tokenChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	tokenCharsLength := len(tokenChars)
	for i := 0; i < length; i++ {
		buffer[i] = tokenChars[int(buffer[i])%tokenCharsLength]
	}

	return string(buffer), nil
}
