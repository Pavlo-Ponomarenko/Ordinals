package mint

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	TxSequence = 1
	CommitTxFeeRate
)

func createCommitTransaction(client *rpcclient.Client, prevOut *wire.OutPoint, spendingAmount int64, pkScript []byte) (*wire.MsgTx, error) {
	tx := new(wire.MsgTx)
	in := wire.NewTxIn(prevOut, nil, nil)
	in.Sequence = TxSequence
	tx.AddTxIn(in)
	tx.AddTxOut(&wire.TxOut{Value: spendingAmount, PkScript: pkScript})
	tx, signed, err := client.SignRawTransactionWithWallet(tx)
	if err != nil || signed == false {
		return nil, errors.New("Signing error")
	}
	if calcFee(tx, CommitTxFeeRate) <= btcutil.Amount(tx.TxOut[0].Value) {
		return nil, errors.New("Fee is too low")
	}
	return tx, nil
}
