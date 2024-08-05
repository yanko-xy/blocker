package types

import (
	"blocker/crypto"
	"blocker/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GeneratePrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))
}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, len(hash), 32)
}
