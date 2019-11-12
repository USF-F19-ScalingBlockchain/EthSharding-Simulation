package transaction

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
	"log"
)

type Identity struct {
	privateKey *rsa.PrivateKey `json:"privateKey"`
	PublicKey  *rsa.PublicKey  `json:"publicKey"`
}

func NewIdentity() Identity {
	id := Identity{}
	id.privateKey, id.PublicKey = GeneratePubPrivKeyPair()
	//id.HashForKey = GenerateHashForKey(label)
	return id
}

func (id *Identity) GetMyPrivateKey() *rsa.PrivateKey {
	pid := id.privateKey
	return pid
}

//generatePubPrivKeyPair func creates key pair pub - priv
func GeneratePubPrivKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Couldnt generate keys - err :", err)
	}
	pub := &priv.PublicKey
	fmt.Println("pub.Size() : ", pub.Size())
	fmt.Println("pub.E : ", pub.E)
	fmt.Println("pub.N : ", pub.N)

	return priv, pub
}

func EncryptMessageWithPublicKey(publicKey *rsa.PublicKey, message string) []byte {

	//hashed := sha3.Sum256([]byte(message))
	//h := sha3.New256()
	//h.Write([]byte(message))
	//h.Sum(nil)

	//encrypted, err := rsa.EncryptOAEP(h, rand.Reader, publicKey, []byte(message), []byte("OAEP"))
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		fmt.Println("Error in encryptMessage : err - ", err)
	} else {
		fmt.Println(">>>>>>>>>>>>>>>>>> Encrypted Successfully")
	}
	return encrypted
}

func DecryptMessageWithPrivateKey(privateKey *rsa.PrivateKey, encrypted []byte) []byte {
	//decrypted, err := rsa.DecryptOAEP(hash1, rand.Reader, privateKey, []byte(encrypted), []byte(label))
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
	if err != nil {
		fmt.Println("Error in decryptMessage : err - ", err)
	} else {
		fmt.Println(">>>>>>>>>>>>>>>>>> De----------crypted Successfully")
	}
	return decrypted
}

//generateHashForKey func creates a hash for id
func GenerateHashForKey(label string) hash.Hash {
	hashForKey := sha256.New()
	hashForKey.Write([]byte(label))
	hashForKey.Sum(nil)
	fmt.Println("For Label : ", label, " Hash Generated is : ", hashForKey)
	return hashForKey
}

func (id *Identity) GenSignature(message []byte) []byte {
	//digest := GenDigest(id.HashForKey, message)//id.HashForKey.Sum(message)
	//signature, err := rsa.SignPSS(rand.Reader, id.privateKey, crypto.SHA256, digest, nil)
	hashMsg := sha3.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, id.privateKey, crypto.SHA256, hashMsg[:])
	if err != nil {
		log.Fatal("Err in generating signature : err : ", err)
	}
	return signature
}

func VerifySingature(senderPubKey *rsa.PublicKey /*senderHashForKey hash.Hash,*/, message []byte, sig []byte) bool { //public key and digest > matched against > sig
	//digest := GenDigest(senderHashForKey, message)
	//err := rsa.VerifyPSS(senderPubKey, crypto.SHA256, digest, sig, nil)
	hashMsg := sha3.Sum256(message)
	err := rsa.VerifyPKCS1v15(senderPubKey, crypto.SHA256, hashMsg[:], sig)
	if err != nil {
		//log.Println("Error in verifying : err : ", err)
		return false
	}
	return true
}

//GetHashOfPubId
func GetHashOfPublicKey(pid *PublicIdentity) string {
	sum := sha3.Sum256(pid.PublicKey.N.Bytes())
	return hex.EncodeToString(sum[:])

}

// GenDigest func returns a hash for both a)GenSignature func and b)VerifySignature
func GenDigest(hash hash.Hash, message []byte) []byte {

	h := hash
	_, err := h.Write(message)
	if err != nil {
		fmt.Println("Error in generating Digest, Error - ", err)
	}
	digest := h.Sum(nil)
	fmt.Println("For Message :", string(message), "  Digest gereated is :", string(digest))
	return digest
}

//func Hello() {
//
//	i := NewIdentity("ok")
//
//	message := "Hello World !!!"
//	fmt.Println(message)
//
//
//	epub	:= 	EncryptMessageWithPublicKey(i.PublicKey, message)
//	dpriv 	:= 	DecryptMessageWithPrivateKey(i.privateKey, epub)
//	fmt.Println("dpriv :", string(dpriv))
//
//
//
//}

func (pid *Identity) PublicIdentityToJson() string {
	jsonBytes, err := json.Marshal(&pid.PublicKey)
	if err != nil {
		log.Println("Error in marshalling publicIdentity, err - ", err)
		return "{}"
	}
	return string(jsonBytes)
}
