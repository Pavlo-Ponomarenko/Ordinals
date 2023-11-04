package mint

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	RevealTxFeeRate = 1
)

func createRevealTransaction(commitTx wire.MsgTx, privateKey *btcec.PrivateKey, inscription []byte, controlBlock []byte, fee int64, pkScript []byte) (*wire.MsgTx, error) {
	tx := new(wire.MsgTx)
	prevOutPoint := &wire.OutPoint{
		Hash:  commitTx.TxHash(),
		Index: 0,
	}
	in := wire.NewTxIn(prevOutPoint, nil, nil)
	in.Sequence = TxSequence
	tx.AddTxIn(in)
	spendingAmount := commitTx.TxOut[0].Value - fee
	tx.AddTxOut(&wire.TxOut{Value: spendingAmount, PkScript: pkScript})
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	prevOutFetcher.AddPrevOut(*prevOutPoint, commitTx.TxOut[0])
	sigHash, err := txscript.CalcTapscriptSignaturehash(txscript.NewTxSigHashes(tx, prevOutFetcher), txscript.SigHashDefault,
		tx, 0, prevOutFetcher, txscript.NewBaseTapLeaf(inscription))
	if err != nil {
		return nil, err
	}
	signature, err := schnorr.Sign(privateKey, sigHash)
	if err != nil {
		return nil, err
	}
	in.Witness = wire.TxWitness{signature.Serialize(), inscription, controlBlock}
	if calcFee(tx, RevealTxFeeRate) > btcutil.Amount(fee) {
		return nil, errors.New("Fee is too low")
	}
	return tx, nil
}
