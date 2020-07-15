package paillier

import (
	"errors"
	"math/big"

	"github.com/wangfeiping/saturn/x/ace/security"
)

var _ security.PrivateKey = PaillierPrivKey{}
var _ security.PublicKey = PaillierPubKey{}

// Create returns a new private key of paillier
func Create() security.PrivateKey {
	return PaillierPrivKey{}
}

// ResumePrivateKey returns a private key of paillier
func ResumePrivateKey(secret security.Secret) security.PrivateKey {
	return PaillierPrivKey{}
}

// ResumePubKey returns a private key of paillier
func ResumePubKey(secret security.Secret) security.PublicKey {
	return PaillierPubKey{}
}

// PaillierPubKey public key of paillier
type PaillierPubKey struct {
}

// Encrypt the plain text and returns the cipher text
func (pub PaillierPubKey) Encrypt(plain []byte) ([]byte, error) {
	return plain, nil
}

// Add two cipher text and returns the result(cipher)
func (pub PaillierPubKey) Add(cipher1, cipher2 []byte) ([]byte, error) {
	x := new(big.Int).SetBytes(cipher1)
	y := new(big.Int).SetBytes(cipher2)

	return new(big.Int).Add(x, y).Bytes(), nil
}

// Mul returns the result(cipher) of multiplying the two cipher text
func (pub PaillierPubKey) Mul(cipher1, cipher2 []byte) ([]byte, error) {
	return nil, errors.New("Multiplication is not supported")
}

func (pub PaillierPubKey) Secret() (security.Secret, error) {
	return security.Secret{}, nil
}

// PaillierPrivKey private key of paillier
type PaillierPrivKey struct {
}

// PublicKey returns the public key for the private key
func (priv PaillierPrivKey) PublicKey() security.PublicKey {
	return PaillierPubKey{}
}

// Decrypt the cipher text and returns the plain text
func (priv PaillierPrivKey) Decrypt(cipher []byte) ([]byte, error) {
	return cipher, nil
}
