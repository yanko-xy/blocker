package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "e3562065491278c4e14020fe0284d19fa98ee9774141c651d862642c00a337c4"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "450262425e70058ff8c79364c6b9c91d84e804ff"
	)
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)
	address := privKey.Public().Address()
	assert.Equal(t, addressStr, address.String())

}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))

	// test with invaild msg
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	// test with invalid publicKey
	invalidPirvKey := GeneratePrivateKey()
	invalidPubKey := invalidPirvKey.Public()
	assert.False(t, sig.Verify(invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()
	assert.Equal(t, len(address.Bytes()), addressLen)
	fmt.Println(address)
}
