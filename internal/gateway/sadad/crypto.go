package sadad

import (
	"bytes"
	"crypto/des"
	"errors"
)

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padText...)
}

func encrypt(origData, key []byte) ([]byte, error) {

	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}

	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

func decrypt(crypted, key []byte) ([]byte, error) {

	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(crypted))
	dst := out

	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

func tripleEcbDesEncrypt(origData, key []byte) ([]byte, error) {

	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]

	block, err := des.NewCipher(k1)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	origData = pkcs5Padding(origData, bs)

	buf1, err := encrypt(origData, k1)
	if err != nil {
		return nil, err
	}

	buf2, err := decrypt(buf1, k2)
	if err != nil {
		return nil, err
	}

	out, err := encrypt(buf2, k3)
	if err != nil {
		return nil, err
	}

	return out, nil
}
