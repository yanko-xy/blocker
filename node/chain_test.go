package node

import (
	"blocker/crypto"
	"blocker/proto"
	"blocker/types"
	"blocker/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	assert.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemeryBlockStore(), NewMemoryTXStore())
	assert.Equal(t, chain.Height(), 0)
	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemeryBlockStore(), NewMemoryTXStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)
		assert.Nil(t, chain.AddBlock(b))
		assert.Equal(t, chain.Height(), i+1)
	}

}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemeryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		blockHash := types.HashBlock(block)

		assert.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlockByHeight, block)
	}
}

func TestAddBlockWithTxInsufficientFunds(t *testing.T) {
	var (
		chain     = NewChain(NewMemeryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(godSeed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("cea8f54544ec77e1cb54150f0ea42e21eba0df8a563cc0b9aee3f7bdc38975c5")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PublicKey:    privKey.Public().Bytes(),
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
		},
	}

	outputs := []*proto.TxOutput{
		{
			Amount:  1001,
			Address: recipient,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	assert.NotNil(t, chain.AddBlock(block))
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		chain     = NewChain(NewMemeryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(godSeed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("cea8f54544ec77e1cb54150f0ea42e21eba0df8a563cc0b9aee3f7bdc38975c5")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PublicKey:    privKey.Public().Bytes(),
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
		},
	}

	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privKey.Public().Address().Bytes(),
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	assert.Nil(t, chain.AddBlock(block))
}
