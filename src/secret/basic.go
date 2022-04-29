package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"
)

func generatPrivateAndPublicKey(privateKey, publicKey io.Writer, bits int) error {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	data := x509.MarshalPKCS1PrivateKey(priKey)
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: data,
	}
	err = pem.Encode(privateKey, &block)
	if err != nil {
		return err
	}
	pubKey := priKey.PublicKey
	pubKeyData := x509.MarshalPKCS1PublicKey(&pubKey)
	err = pem.Encode(publicKey, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyData,
	})
	if err != nil {
		return err
	}
	return nil
}

func FileRSAPrivateAndPublicKeyGenerat(bits int, publicWhere string, privateWhere string) error {
	files1, err := os.OpenFile(publicWhere, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	files2, err2 := os.OpenFile(privateWhere, os.O_RDWR|os.O_CREATE, 0766)
	if err2 != nil {
		return err
	}
	return generatPrivateAndPublicKey(files1, files2, bits)
}

func SecretPrivateFile(privateWhere string, password string) error {
	plaintext, err := ioutil.ReadFile(privateWhere)
	if err != nil {
		panic(err.Error())
	}
	key := []byte(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	f, err := os.Create(privateWhere)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	if err != nil {
		panic(err.Error())
	}
	return nil
}
