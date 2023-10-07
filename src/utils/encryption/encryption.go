package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os/exec"
)

func GetKey(privateKey []byte) (*rsa.PrivateKey) {
    block, _ := pem.Decode(privateKey)
    key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    return key.(*rsa.PrivateKey)
}

func GetRequestHashSum(str *string) []byte {
    var data = []byte(*str)
    msgHash := sha256.New()
    _, err := msgHash.Write(data)
    if err != nil {
        panic(err)
    }
    return msgHash.Sum(nil)
}

func Sign(privateKey *rsa.PrivateKey, requestHashSum []byte) []byte {
    sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, requestHashSum)
    if err != nil {
        panic(err)
    }
    return sign
}

func Verify(publicKey *rsa.PublicKey, requestHashSum, signature []byte) {
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, requestHashSum, signature)
	if err != nil {
		panic(err)
	}
}

func Encrypt(privateKey []byte, str *string) []byte {
	pKey := GetKey(privateKey)
    requestHashSum := GetRequestHashSum(str)

    signature := Sign(pKey, requestHashSum)

	return signature
}

func EncryptAesKey(taxServerPublicKey string, key string) ([]byte, error) {
    publicKeyPEM := []byte("-----BEGIN PUBLIC KEY-----\n" + taxServerPublicKey + "\n-----END PUBLIC KEY-----")
    block, _ := pem.Decode(publicKeyPEM)
    publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)

    // Convert public key to RSA public key
    rsaPublicKey := publicKey.(*rsa.PublicKey)

    label := []byte("")
    hash := sha256.New()
    cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, rsaPublicKey, []byte(key), label)
    if err != nil {
        return nil, err
    }
    return cipherText, nil
}

func XorAndEncryptData(data string, key string, iv string) (string, error) {
    out, err := exec.Command("php", "helper.php", key, iv, data).Output()
    if err != nil {
        return "", err
    }
    return string(out), nil
}