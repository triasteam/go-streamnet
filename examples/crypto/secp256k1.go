package main

import (
	"github.com/triasteam/go-streamnet/utils/crypto"
	"github.com/triasteam/go-streamnet/utils/crypto/secp256k1"
)

func main() {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()

	msg := crypto.CRandBytes(128)
	sig, err := privKey.Sign(msg)
	println("err = ", err)
	//require.Nil(t, err)

	//assert.True(t, pubKey.VerifyBytes(msg, sig))
	r := pubKey.VerifyBytes(msg, sig)
	println(r)

	// Mutate the signature, just one bit.
	sig[3] ^= byte(0x01)
	//assert.False(t, pubKey.VerifyBytes(msg, sig))
	r = pubKey.VerifyBytes(msg, sig)
	println(r)
}
