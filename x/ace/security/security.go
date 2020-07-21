package security

type PublicKey interface {
	Encrypt(plain []byte) ([]byte, error)
	Add(cipher1, cipher2 []byte) ([]byte, error)
	Mul(cipher1, cipher2 []byte) ([]byte, error)
	Secret() (Secret, error)
}

type PrivateKey interface {
	PublicKey() PublicKey
	Decrypt(cipher []byte) ([]byte, error)
}
