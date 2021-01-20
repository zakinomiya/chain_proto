package gateway

import (
	"chain_proto/account"
	"chain_proto/block"
	"chain_proto/common"
	"chain_proto/config"
	"chain_proto/gateway/gw"
	"chain_proto/peer"
	"chain_proto/transaction"
	"chain_proto/wallet"
	"encoding/hex"
	"fmt"

	"github.com/shopspring/decimal"
)

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
	tx.TxType = transaction.Normal
	tx.Timestamp = t.GetTimestamp()
	tx.Signature = sig
	tx.TotalValue = totalValue
	tx.Outs = outs

	return tx, nil
}

func toPbTransaction(t *transaction.Transaction) *gw.Transaction {
	return &gw.Transaction{
		TxHash:     fmt.Sprintf("%x", t.TxHash[:]),
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

func toPbAccount(a *account.Account) *gw.Account {
	return &gw.Account{
		Addr:    a.Addr,
		Balance: a.BalanceString(),
	}
}

func toAccount(a *gw.Account) (*account.Account, error) {
	balance, err := decimal.NewFromString(a.GetBalance())
	if err != nil {
		return nil, err
	}

	return &account.Account{
		Addr:    a.GetAddr(),
		Balance: balance,
	}, nil
}
