package mint

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/wire"
)

func calcFee(tx *wire.MsgTx, feeRate int) btcutil.Amount {
	return btcutil.Amount(mempool.GetTxVirtualSize(btcutil.NewTx(tx))) * btcutil.Amount(feeRate) / btcutil.Amount(10)
}
