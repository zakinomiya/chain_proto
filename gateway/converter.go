package gateway

import (
	"chain_proto/block"
	"chain_proto/common"
	"chain_proto/config"
	"chain_proto/gateway/gw"
	"chain_proto/peer"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"encoding/hex"
	"fmt"
)

// Height        uint32         `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
//	Hash          string         `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
//	PrevBlockHash string         `protobuf:"bytes,3,opt,name=prev_block_hash,json=prevBlockHash,proto3" json:"prev_block_hash,omitempty"`
//	ExtraNonce    uint32         `protobuf:"varint,4,opt,name=extra_nonce,json=extraNonce,proto3" json:"extra_nonce,omitempty"`
//	MerkleRoot    string         `protobuf:"bytes,5,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
//	Timestamp     uint64         `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
//	Bits          uint32         `protobuf:"varint,7,opt,name=bits,proto3" json:"bits,omitempty"`
//	Nonce         uint32         `protobuf:"varint,8,opt,name=nonce,proto3" json:"nonce,omitempty"`
//	Transactions  []*Transaction `protobuf:"bytes,9,rep,name=transactions,proto3" json:"transactions,omitempty"`

func toPbBlock(b *block.Block) (*gw.Block, error) {
	pbTxs := make([]*gw.Transaction, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		pbTxs = append(pbTxs, toPbTransaction(tx))
	}

	return &gw.Block{
		Height:        b.Height,
		Hash:          fmt.Sprintf("%x", b.Hash),
		PrevBlockHash: fmt.Sprintf("%x", b.PrevBlockHash),
		ExtraNonce:    b.ExtraNonce,
		MerkleRoot:    fmt.Sprintf("%x", b.MerkleRoot),
		Timestamp:     uint64(b.Timestamp),
		Bits:          b.Bits,
		Nonce:         b.Nonce,
		Transactions:  pbTxs,
	}, nil
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

	totalValue, err := common.ToDecimal(t.GetTotalValue(), "")
	if err != nil {
		return nil, err
	}

	fee, err := common.ToDecimal(t.GetFee(), "")
	if err != nil {
		return nil, err
	}

	outs, err := toOutputs(t.GetOuts())
	if err != nil {
		return nil, err
	}

	tx.TxHash = common.ReadByteInto32(hash)
	tx.Fee = fee
	tx.SenderAddr = t.GetSenderAddr()
	tx.TxType = transaction.TxType(gw.Transaction_TxType_name[int32(t.GetTxType())])
	tx.Timestamp = t.GetTimestamp()
	tx.Signature = sig
	tx.TotalValue = totalValue
	tx.Outs = outs

	return tx, nil
}

func toPbTransaction(t *transaction.Transaction) *gw.Transaction {
	return &gw.Transaction{
		TxHash:     fmt.Sprintf("%x", t.TxHash[:]),
		TxType:     gw.Transaction_TxType(gw.Transaction_TxType_value[string(t.TxType)]),
		TotalValue: t.TotalValue.StringFixed(config.MaxDecimalDigit),
		Timestamp:  t.Timestamp,
		Fee:        t.Fee.StringFixed(config.MaxDecimalDigit),
		SenderAddr: t.SenderAddr,
		Signature:  t.Signature.String(),
		Outs:       toPbOutputs(t.Outs),
	}

}

func toOutputs(pbOuts []*gw.Output) ([]*transaction.Output, error) {
	outs := make([]*transaction.Output, 0, len(pbOuts))
	for _, pbOut := range pbOuts {
		value, err := common.ToDecimal(pbOut.GetValue(), "")
		if err != nil {
			return nil, err
		}

		out := &transaction.Output{
			RecipientAddr: pbOut.GetRecipientAddr(),
			Value:         value,
		}
		outs = append(outs, out)
	}

	return outs, nil
}

func toPbOutputs(outs []*transaction.Output) []*gw.Output {
	pbOuts := make([]*gw.Output, 0, len(outs))
	for _, out := range outs {
		pbOut := &gw.Output{
			Value:         out.Value.StringFixed(config.MaxDecimalDigit),
			RecipientAddr: out.RecipientAddr,
		}
		pbOuts = append(pbOuts, pbOut)
	}

	return pbOuts
}

func toPbPeer(p *peer.Peer) *gw.Peer {
	return &gw.Peer{
		Addr:    p.Addr(),
		Network: p.Network(),
	}
}

func toPeer(p *gw.Peer) *peer.Peer {
	return peer.New(p.Addr, p.Network)
}
