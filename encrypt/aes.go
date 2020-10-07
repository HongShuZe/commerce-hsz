package encrypt

import (
	"bytes"
	"errors"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 高级加密标准（Advanced Encryption Standard ,AES）

//16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
// key不能泄露
var PwdKey = []byte("hszz123456789hsz")

// PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// Repeat()函数功能是把切片[]byte{byte(padding)}复制padding个, 然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反操作, 删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	// 获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串失败")
	} else {
		// 获取填充字符串长度(AES会把填充字符串长度附在加密串上)
		unpadding := int(origData[length-1])
		// 截取切片, 删除填充字节, 返回明文
		return origData[:(length - unpadding)], nil
	}
}

// 实现加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块大小
	blockSize := block.BlockSize()
	// 对数据进行填充, 然数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	// 采用AES加密方法中的CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 执行加密
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 实现解密
func AesDecrypt(crypted []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块大小
	blockSize := block.BlockSize()
	// 创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// 该函数可以加解密
	blockMode.CryptBlocks(origData, crypted)
	// 去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// 加密base64
func EnPwdCode(pwd []byte) (string, error) {
	// AES加密
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	// base64加密
	return base64.StdEncoding.EncodeToString(result), err
}

// 解密
func DePwdCode(pwd string) ([]byte, error) {
	// 解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// 执行AES解密
	return AesDecrypt(pwdByte, PwdKey)
}
















