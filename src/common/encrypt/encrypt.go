package encrypt

type IEncrypt interface {
	// Encrypt
	AesEncrypt(plain, key []byte) ([]byte, error)
	// Decrypt
	AesDecrypt(cipher, key []byte) ([]byte, error)
}