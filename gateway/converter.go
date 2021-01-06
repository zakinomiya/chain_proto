package gateway

import (
	"chain_proto/block"
	"chain_proto/common"
	"chain_proto/gateway/gw"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"encoding/hex"
	"fmt"
)

func toPbBlock(b *block.Block) (*gw.Block, error) {
	return &gw.Block{}, nil
}

func toBlock(b *gw.Block) (*block.Block, error) {
	hash, err := hex.DecodeString(b.GetHash())
	if err != nil {
		return nil, err
	}

	prevBlockHash, err := hex.DecodeString(b.GetPrevBlockHash())
	if err != nil {
		return nil, err
	}

	merkleRoot, err := hex.DecodeString(b.GetMerkleRoot())
	if err != nil {
		return nil, err
	}

	txs := make([]*transaction.Transaction, 0, len(b.GetTransactions()))
	for _, pbTx := range b.GetTransactions() {
		tx, err := toTransaction(pbTx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	blk := block.New(b.GetHeight(), b.GetBits(), common.ReadByteInto32(prevBlockHash), txs)

	blk.ExtraNonce = b.GetExtraNonce()
	blk.Hash = common.ReadByteInto32(hash)
	blk.MerkleRoot = merkleRoot
	blk.Nonce = b.GetNonce()
	blk.Timestamp = int64(b.GetTimestamp())

	return blk, nil
}

func toTransaction(t *gw.Transaction) (*transaction.Transaction, error) {
	tx := transaction.New()

	hash, err := hex.DecodeString(t.GetTxHash())
	if err != nil {
		return nil, err
	}

	sig, err := wallet.DecodeSigString(t.GetSignature())
	if err != nil {
		return nil, err
	}

	tx.TxHash = common.ReadByteInto32(hash)
	tx.Fee = t.GetFee()
	tx.SenderAddr = t.GetSenderAddr()
	tx.TxType = transaction.TxType(gw.Transaction_TxType_name[int32(t.GetTxType())])
	tx.Timestamp = t.GetTimestamp()
	tx.Signature = sig
	tx.TotalValue = t.GetTotalValue()
	tx.Outs = toOutputs(t.GetOuts())

	return tx, nil
}

func toPbTransaction(t *transaction.Transaction) *gw.Transaction {

	return &gw.Transaction{
		TxHash:     fmt.Sprintf("%x", t.TxHash[:]),
		TxType:     gw.Transaction_TxType(gw.Transaction_TxType_value[string(t.TxType)]),
		TotalValue: t.TotalValue,
		Timestamp:  t.Timestamp,
		Fee:        t.Fee,
		SenderAddr: t.SenderAddr,
		Signature:  t.Signature.String(),
		Outs:       toPbOutputs(t.Outs),
	}

}

func toOutputs(pbOuts []*gw.Output) []*transaction.Output {
	outs := make([]*transaction.Output, 0, len(pbOuts))
	for _, pbOut := range pbOuts {
		out := &transaction.Output{
			RecipientAddr: pbOut.GetRecipientAddr(),
			Value:         pbOut.GetValue(),
		}
		outs = append(outs, out)
	}

	return outs
}

func toPbOutputs(outs []*transaction.Output) []*gw.Output {
	pbOuts := make([]*gw.Output, 0, len(outs))
	for _, out := range outs {
		pbOut := &gw.Output{
			Value:         out.Value,
			RecipientAddr: out.RecipientAddr,
		}
		pbOuts = append(pbOuts, pbOut)

	}

	return pbOuts
}
