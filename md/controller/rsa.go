package controller

import (
	"md/model/common"
	"md/util"

	"github.com/kataras/iris/v12"
)

type rsaCondition struct {
	Message    string        `json:"message"`
	Bits       int           `json:"bits"`
	IsPKCS8    bool          `json:"isPKCS8"`
	PrivateKey string        `json:"privateKey"`
	PublicKey  string        `json:"publicKey"`
	Sign       string        `json:"sign"`
	SignType   util.SignType `json:"signType"`
}

type rsaResult struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

// 生成RSA密钥对
func RSAGenerateKey(ctx iris.Context) {
	rsaCondition := rsaCondition{}
	resolveParam(ctx, &rsaCondition)
	privateKey, publicKey, err := util.GenerateRSAKey(rsaCondition.Bits, rsaCondition.IsPKCS8)
	if err != nil {
		panic(common.NewErr("RSA密钥对生成失败", err))
	}

	ctx.JSON(common.NewSuccessData("RSA密钥对生成成功", rsaResult{privateKey, publicKey}))
}

// RSA公钥加密
func RSAEncrypt(ctx iris.Context) {
	rsaCondition := rsaCondition{}
	resolveParam(ctx, &rsaCondition)
	if rsaCondition.Message == "" {
		panic(common.NewError("加密内容不可为空"))
	}

	if rsaCondition.PublicKey == "" {
		panic(common.NewError("公钥不可为空"))
	}

	result, err := util.EncryptRSA(rsaCondition.Message, rsaCondition.PublicKey)
	if err != nil {
		panic(common.NewErr("公钥加密失败", err))
	}

	ctx.JSON(common.NewSuccessData("公钥加密成功", result))
}

// RSA私钥解密
func RSADecrypt(ctx iris.Context) {
	rsaCondition := rsaCondition{}
	resolveParam(ctx, &rsaCondition)
	if rsaCondition.Message == "" {
		panic(common.NewError("解密内容不可为空"))
	}

	if rsaCondition.PrivateKey == "" {
		panic(common.NewError("私钥不可为空"))
	}

	result, err := util.DecryptRSA(rsaCondition.Message, rsaCondition.PrivateKey, rsaCondition.IsPKCS8)
	if err != nil {
		panic(common.NewErr("私钥解密失败", err))
	}

	ctx.JSON(common.NewSuccessData("私钥解密成功", result))
}

// RSA私钥签名
func RSASign(ctx iris.Context) {
	rsaCondition := rsaCondition{}
	resolveParam(ctx, &rsaCondition)
	if rsaCondition.Message == "" {
		panic(common.NewError("签名内容不可为空"))
	}

	if rsaCondition.PrivateKey == "" {
		panic(common.NewError("私钥不可为空"))
	}

	if rsaCondition.SignType == "" {
		panic(common.NewError("签名方式不可为空"))
	}

	result, err := util.SignRSA(rsaCondition.Message, rsaCondition.PrivateKey, rsaCondition.SignType, rsaCondition.IsPKCS8)
	if err != nil {
		panic(common.NewErr("私钥签名失败", err))
	}

	ctx.JSON(common.NewSuccessData("私钥签名成功", result))
}

// RSA公钥验签
func RSAVerify(ctx iris.Context) {
	rsaCondition := rsaCondition{}
	resolveParam(ctx, &rsaCondition)
	if rsaCondition.Message == "" {
		panic(common.NewError("验签内容不可为空"))
	}

	if rsaCondition.PublicKey == "" {
		panic(common.NewError("公钥不可为空"))
	}

	if rsaCondition.Sign == "" {
		panic(common.NewError("签名不可为空"))
	}

	if rsaCondition.SignType == "" {
		panic(common.NewError("验签方式不可为空"))
	}

	err := util.VerifyRSA(rsaCondition.Message, rsaCondition.PublicKey, rsaCondition.Sign, rsaCondition.SignType)
	if err != nil {
		ctx.JSON(common.NewSuccessData("公钥验签失败", false))
	} else {
		ctx.JSON(common.NewSuccessData("公钥验签成功", true))
	}
}
