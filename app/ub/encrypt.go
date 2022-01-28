package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"test/app/ub/aes"

	"github.com/davecgh/go-spew/spew"
)

func encrypt() {

	sign := "U4lwqIn7llDlEhH86qjr9ZFY8ABcXSO9Ec+NM/oHadvH6JOaJemqwOa4t2COXST+BqLa4HgA/VPqCIUzFnvjswfRq79zvPJOHhEPKriPSTLNVdHoiALZARDU8ribzQa0LKFpIkkKyJsDZDY8oj6JNHWXJ8qmh5X/r/P641PJ4SgIcmJpCs2f/nOFQifPoRjnAvTSE8/WFr+/+dSR4llB6AVZM0JtyLAkLy2a+LAYckd73gXBOndyZGUqXUF84ZFYCW3ZMLf2GPHcGSgxi3baH4QYZNo2Wf6NOKc/iDgyvJdW7c+Wgv0VkYnCf7c6estGcOWarCZ3bdriBtvJV3s8og=="

	mac := "mDdId6ZoHv67g1WzsSs+toERXpPzivwK4Adjt8zAOuEQR8m5NevsQ7bYqSGyn2Ba"

	data := "{\"acc\":\"004100037912\",\"date\":\"20211125\",\"txseq\":\"TYR1637819929763\",\"ubnotify\":\"record\",\"amt\":\"00000003000000\",\"wdacc\":\"0000002200910040\",\"wdbank\":\"803\",\"stan\":\"073SO00003\",\"to\":\"99327\",\"time\":\"135848\",\"ecacc\":\"99327803323902\",\"status\":\"0\",\"txnid\":\"1486\"}"

	suc, err := VerifyMac(data, mac)

	spew.Dump(suc)
	spew.Dump(err)

	err = VerifySignature(sign, mac)

	spew.Dump(err)

}

// encrypt - 加密data
func Encrypt(src string) (string, error) {
	sha := sha256.New()
	sha.Write([]byte(src))
	// Sha
	sha256Hash := sha.Sum(nil)
	// Base64
	sha256Base64 := base64.StdEncoding.EncodeToString(sha256Hash)

	dat, _ := ioutil.ReadFile("./aes_key")
	aesCrypto := aes.NewAESCrypto(dat, []byte("UBOTSECRETIVSEED"))
	// AES
	encryptedData, err := aesCrypto.AesEncrypt(sha256Base64)
	if err != nil {
		return "", err
	}
	// Base64
	encryptedString := base64.StdEncoding.EncodeToString(encryptedData)

	// Mac
	return encryptedString, nil
}

// verifyMac - 驗證mac
func VerifyMac(src, mac string) (bool, error) {
	encryptSrc, err := Encrypt(src)
	if err != nil {
		return false, err
	}
	if encryptSrc != mac {
		return false, nil
	}
	return true, nil
}

// 驗證signature
func VerifySignature(signature, mac string) error {
	dat, _ := ioutil.ReadFile("./UBOT_PUBLIC.PEM")

	pub := bytesToPublicKey(dat)

	sDec, _ := base64.StdEncoding.DecodeString(signature)
	macByte := []byte(mac)
	digest := sha256.Sum256(macByte)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sDec)
}
func bytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		spew.Dump("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			spew.Dump(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		spew.Dump(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		spew.Dump("not ok")
	}
	return key
}
