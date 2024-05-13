// 加密工具类
package util

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

type SignType string

const (
	MD5    SignType = "MD5"
	SHA1   SignType = "SHA1"
	SHA256 SignType = "SHA256"
	SHA512 SignType = "SHA512"
)

// MD5加密
func EncryptMD5(message []byte) string {
	hash := md5.New()
	hash.Write(message)
	sum := hash.Sum(nil)
	hashCode := hex.EncodeToString(sum)
	return hashCode
}

// SHA1加密
func EncryptSHA1(message []byte) string {
	hash := sha1.New()
	hash.Write(message)
	sum := hash.Sum(nil)
	hashCode := hex.EncodeToString(sum)
	return hashCode
}

// SHA256加密
func EncryptSHA256(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	sum := hash.Sum(nil)
	hashCode := hex.EncodeToString(sum)
	return hashCode
}

// SHA512加密
func EncryptSHA512(message []byte) string {
	hash := sha512.New()
	hash.Write(message)
	sum := hash.Sum(nil)
	hashCode := hex.EncodeToString(sum)
	return hashCode
}

// BASE64编码
func EncryptBASE64(message []byte) string {
	return base64.StdEncoding.EncodeToString(message)
}

// BASE64解码
func DecryptBASE64(message string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(message)
}

// GenerateRSAKey 函数用于生成RSA密钥对
// 参数 bits 表示密钥的位数
// 参数 isPKCS8 表示是否使用PKCS8编码格式
// 返回生成的私钥和公钥的字符串表示以及可能的错误
func GenerateRSAKey(bits int, isPKCS8 bool) (string, string, error) {
	// 检查密钥位数是否在合法范围内
	if bits < 512 || bits > 4096 {
		return "", "", errors.New("密钥位数需在512-4096之间")
	}

	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	var privateDer []byte
	if isPKCS8 {
		// 使用PKCS8编码格式生成私钥的DER格式
		privateDer, err = x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return "", "", err
		}
	} else {
		// 使用PKCS1编码格式生成私钥的DER格式
		privateDer = x509.MarshalPKCS1PrivateKey(privateKey)
	}

	// 生成公钥
	publicDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	// 将私钥和公钥的DER格式转换为BASE64编码的字符串
	privateKeyStr := EncryptBASE64(privateDer)
	publicKeyStr := EncryptBASE64(publicDer)

	return privateKeyStr, publicKeyStr, nil
}

// EncryptRSA 函数用于使用RSA公钥加密数据
// 参数 message 表示要加密的数据
// 参数 publicKey 表示RSA公钥的字符串
// 返回加密后的数据的字符串表示以及可能的错误
func EncryptRSA(message, publicKey string) (string, error) {
	// base64解码公钥字符串
	key, err := DecryptBASE64(publicKey)
	if err != nil {
		return "", err
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return "", err
	}

	// 使用公钥加密数据
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(message))
	if err != nil {
		return "", err
	}

	// 将加密后的数据转换为BASE64编码的字符串
	return EncryptBASE64(encryptedData), nil
}

// DecryptRSA 函数用于使用RSA私钥解密数据
// 参数 message 表示要解密的数据的字符串表示
// 参数 privateKey 表示RSA私钥的字符串表示
// 参数 isPKCS8 表示私钥是否使用PKCS8编码格式
// 返回解密后的数据的字符串表示以及可能的错误
func DecryptRSA(message, privateKey string, isPKCS8 bool) (string, error) {
	// base64解码要解密的数据
	messageBytes, err := DecryptBASE64(message)
	if err != nil {
		return "", err
	}

	// base64解码私钥
	key, err := DecryptBASE64(privateKey)
	if err != nil {
		return "", err
	}

	var priKey interface{}
	if isPKCS8 {
		// 使用PKCS8编码格式解析私钥
		priKey, err = x509.ParsePKCS8PrivateKey(key)
	} else {
		// 使用PKCS1编码格式解析私钥
		priKey, err = x509.ParsePKCS1PrivateKey(key)
	}
	if err != nil {
		return "", err
	}

	// 使用私钥解密数据
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), messageBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

// SignRSA 函数用于使用 RSA 私钥对消息进行签名
// 参数 message 表示要签名的消息
// 参数 privateKey 表示 RSA 私钥
// 参数 signType 表示签名类型
// 参数 isPKCS8 表示私钥是否采用 PKCS8 格式
// 返回签名结果和可能的错误
func SignRSA(message, privateKey string, signType SignType, isPKCS8 bool) (string, error) {
	// 解码私钥
	key, err := DecryptBASE64(privateKey)
	if err != nil {
		return "", err
	}

	var priKey interface{}

	// 根据 isPKCS8 参数决定解析私钥的方式
	if isPKCS8 {
		priKey, err = x509.ParsePKCS8PrivateKey(key)
	} else {
		priKey, err = x509.ParsePKCS1PrivateKey(key)
	}
	if err != nil {
		return "", err
	}

	var signature []byte

	// 根据签名类型选择哈希算法，并对消息进行哈希计算
	switch signType {
	case MD5:
		h := md5.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		signature, err = rsa.SignPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), crypto.MD5, hash)
	case SHA1:
		h := sha1.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		signature, err = rsa.SignPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), crypto.SHA1, hash)
	case SHA256:
		h := sha256.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		signature, err = rsa.SignPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), crypto.SHA256, hash)
	case SHA512:
		h := sha512.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		signature, err = rsa.SignPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), crypto.SHA512, hash)
	default:
		return "", errors.New("不支持的签名类型")
	}
	if err != nil {
		return "", err
	}

	// DER格式转换为BASE64编码的字符串
	return EncryptBASE64(signature), nil
}

// VerifyRSA 函数用于RSA公钥验签
// 参数 message 表示待验签的消息
// 参数 publicKey 表示公钥
// 参数 sign 表示签名
// 参数 signType 表示签名类型
// 返回验签结果，如果验签成功返回nil，否则返回错误
func VerifyRSA(message, publicKey, sign string, signType SignType) error {
	// base64解码签名和公钥
	signBytes, err := DecryptBASE64(sign)
	if err != nil {
		return err
	}
	key, err := DecryptBASE64(publicKey)
	if err != nil {
		return err
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return err
	}

	// 根据签名类型选择哈希算法
	switch signType {
	case MD5:
		h := md5.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.MD5, hash, signBytes)
	case SHA1:
		h := sha1.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA1, hash, signBytes)
	case SHA256:
		h := sha256.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hash, signBytes)
	case SHA512:
		h := sha512.New()
		h.Write([]byte(message))
		hash := h.Sum(nil)
		err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA512, hash, signBytes)
	default:
		return errors.New("不支持的签名类型")
	}

	// 返回验签结果
	if err != nil {
		return err
	}
	return nil
}

// pkcs7Padding 函数用于进行 PKCS7 填充
// 参数 message 表示要填充的消息
// 参数 blockSize 表示块大小
// 返回填充后的消息
func pkcs7Padding(message []byte, blockSize int) []byte {
	padding := blockSize - len(message)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(message, padText...)
}

// pkcs7UnPadding 函数用于进行 PKCS7 反填充
// 参数 message 表示要进行反填充的消息
// 返回反填充后的消息
func pkcs7UnPadding(message []byte) []byte {
	length := len(message)
	unPadding := int(message[length-1])
	return message[:(length - unPadding)]
}

// paddingKey 函数用于将 key 向上填充为 16 位、24 位、32 位，超过 32 位则截取前 32 位
// 参数 key 表示要进行填充的 key
// 返回填充后的 key 字节
func paddingKey(key string) []byte {
	keyByte := []byte(key)
	keyLen := len(key)
	if keyLen < 16 {
		for i := 0; i < 16-keyLen; i++ {
			keyByte = append(keyByte, 0)
		}
	} else if keyLen < 24 {
		for i := 0; i < 24-keyLen; i++ {
			keyByte = append(keyByte, 0)
		}
	} else if keyLen < 32 {
		for i := 0; i < 32-keyLen; i++ {
			keyByte = append(keyByte, 0)
		}
	} else if keyLen > 32 {
		keyByte = keyByte[:32]
	}
	return keyByte
}

// EncryptAES 函数用于进行 AES 加密，使用 CBC 模式
// 参数 message 表示要加密的消息
// 参数 key 表示加密密钥
// 参数 isHex 表示是否返回十六进制字符串
// 返回加密后的结果字符串和可能的错误
func EncryptAES(message, key string, isHex bool) (res string, err error) {
	// 使用 defer 捕获可能的 panic，并转换为错误
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("AES加密异常")
		}
	}()

	// 将消息和密钥转换为字节数组
	messageByte := []byte(message)
	// 函数用于将 key 向上填充为 16 位、24 位、32 位，超过 32 位则截取前 32 位
	keyByte := paddingKey(key)

	// 创建 AES 加密块
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	// 获取块大小并对消息进行 PKCS7 填充
	blockSize := block.BlockSize()
	// 函数用于进行 PKCS7 填充
	messageByte = pkcs7Padding(messageByte, blockSize)

	// 创建 CBC 加密器并进行加密
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	result := make([]byte, len(messageByte))
	blockMode.CryptBlocks(result, messageByte)

	// 根据 isHex 参数决定返回结果的格式
	if isHex {
		return hex.EncodeToString(result), err
	}

	// DER格式转换为BASE64编码的字符串
	return EncryptBASE64(result), err
}

// DecryptAES 函数用于进行 AES 解密，使用 CBC 模式
// 参数 message 表示要解密的消息
// 参数 key 表示解密密钥
// 参数 isHex 表示消息是否为十六进制格式
// 返回解密后的消息和可能的错误
func DecryptAES(message, key string, isHex bool) (res string, err error) {
	// 使用 defer 和 recover 捕获可能的异常，并将错误信息设置为 "AES解密异常"
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("AES解密异常")
		}
	}()

	var messageByte []byte
	if isHex {
		// 如果消息为十六进制格式，则将其解码为字节
		messageByte, err = hex.DecodeString(message)
	} else {
		// 否则，将消息解码为 BASE64 格式的字节
		messageByte, err = DecryptBASE64(message)
	}
	if err != nil {
		return "", err
	}

	// 函数用于将 key 向上填充为 16 位、24 位、32 位，超过 32 位则截取前 32 位
	keyByte := paddingKey(key)

	// 创建 AES 密码器
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	// 获取块大小
	blockSize := block.BlockSize()
	// 创建 CBC 解密器
	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])
	// 创建结果字节切片
	result := make([]byte, len(messageByte))
	// 解密消息
	blockMode.CryptBlocks(result, messageByte)
	// 函数用于进行 PKCS7 反填充
	result = pkcs7UnPadding(result)

	// 将结果转换为字符串并返回
	return string(result), err
}
