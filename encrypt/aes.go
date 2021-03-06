package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 高级加密标准 (Adevanced Encruption Standard, AES)

// 16,24,32为字符串的话，分别对应AES-128, AES-192, AES-256 加密方法
// key不能泄露
var PwdKey = []byte("DIS**#KKKDJJSKDI")

// pkcs7填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// Repeat()的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反向操作， 删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	// 获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		// 获取填充字符串长度
		unPadding := int(origData[length-1])
		// 截取切片，删除填充字节， 并且返回明文
		return origData[:(length - unPadding)], nil
	}
}

// 实现加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取快的大小
	blockSize := block.BlockSize()
	// 对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	// 采用AES加密方法中的CBC加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	// 执行加密
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// 实现解密
func AesDeCrypt(encrypted []byte, key []byte) ([]byte, error) {
	// 实现加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取快的大小
	blockSize := block.BlockSize()
	// 创建加密客户端实例
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	// 这个函数也可以用来解密
	blockMode.CryptBlocks(origData, encrypted)
	// 去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

// 加密base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// 解密base64
func DePwdCode(pwd string) ([]byte, error) {
	// 解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// 执行AES解密
	return AesDeCrypt(pwdByte, PwdKey)
}
