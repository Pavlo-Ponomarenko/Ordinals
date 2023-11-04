package mint

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	TxSequence = 1
	CommitTxFeeRate
)

func createCommitTransaction(sourceAddress btcutil.Address, privateKey *btcec.PrivateKey, prevOutPoint *wire.OutPoint, totalAmount int64, spendingAmount int64, fee int64, returnCoinsScript []byte, pkScript []byte) (*wire.MsgTx, error) {
	tx := new(wire.MsgTx)
	tx.AddTxOut(&wire.TxOut{Value: spendingAmount, PkScript: pkScript})
	returnValue := totalAmount - spendingAmount - fee
	tx.AddTxOut(&wire.TxOut{Value: returnValue, PkScript: returnCoinsScript})
	in := wire.NewTxIn(prevOutPoint, nil, nil)
	in.Sequence = TxSequence
	tx.AddTxIn(in)
	witnessProgram, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return nil, err
	}
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	prevOutFetcher.AddPrevOut(*prevOutPoint, getPreviousOut(sourceAddress, totalAmount))
	hashCache := txscript.NewTxSigHashes(tx, prevOutFetcher)
	witnessScript, err := txscript.WitnessSignature(tx, hashCache, 0,
		totalAmount, witnessProgram, txscript.SigHashAll, privateKey, true)
	if err != nil {
		return nil, err
	}
	in.Witness = witnessScript
	if calcFee(tx, CommitTxFeeRate) > btcutil.Amount(fee) {
		return nil, errors.New("Fee is too low")
	}
	return tx, nil
}

func getPreviousOut(address btcutil.Address, value int64) *wire.TxOut {
	returnToSourceScript, _ := txscript.PayToAddrScript(address)
	return &wire.TxOut{
		Value:    value,
		PkScript: returnToSourceScript,
	}
}
