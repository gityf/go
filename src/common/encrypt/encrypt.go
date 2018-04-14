package encrypt

type IEncrypt interface {
	// Encrypt
	Encrypt(plain, key []byte) ([]byte, error)
	// Decrypt
	Decrypt(cipher, key []byte) ([]byte, error)
}