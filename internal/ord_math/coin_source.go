package ord_math

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func FindOutputOrder(client *rpcclient.Client, tx *wire.MsgTx, outputIndex uint32) int64 {
	coinbaseTx, coins, output := FindCoinSource(client, tx, outputIndex)
	for i := uint32(0); i < output; i++ {
		coins += coinbaseTx.TxOut[i].Value
	}
	coinbaseHash := coinbaseTx.TxHash()
	txData, _ := client.GetRawTransactionVerbose(&coinbaseHash)
	blockHash, _ := chainhash.NewHashFromStr(txData.BlockHash)
	block, _ := client.GetBlockVerbose(blockHash)
	return firstOrdinal(block.Height) + coins - 1
}

func FindCoinSource(client *rpcclient.Client, tx *wire.MsgTx, outputIndex uint32) (*wire.MsgTx, int64, uint32) {
	currentOrder := int64(0)
	for !isCoinbase(tx) {
		outputValues := extractOutputValues(tx.TxOut)
		previousOutputs := getPreviousOutputs(client, tx.TxIn)
		previousOutputValues := extractOutputValues(previousOutputs)
		inputIndex := findRelation(outputIndex, outputValues, previousOutputValues)
		currentOrder += calcCurrentOutputOrder(outputIndex, inputIndex, outputValues, previousOutputValues)
		input := tx.TxIn[inputIndex]
		outPoint := input.PreviousOutPoint
		txData, _ := client.GetRawTransaction(&outPoint.Hash)
		tx = txData.MsgTx()
		outputIndex = outPoint.Index
	}
	return tx, currentOrder, outputIndex
}

func extractOutputValues(outList []*wire.TxOut) []int64 {
	result := make([]int64, 0)
	for _, out := range outList {
		result = append(result, out.Value)
	}
	return result
}

func getPreviousOutputs(client *rpcclient.Client, inList []*wire.TxIn) []*wire.TxOut {
	result := make([]*wire.TxOut, 0)
	for _, in := range inList {
		outPoint := in.PreviousOutPoint
		tx, _ := client.GetRawTransaction(&outPoint.Hash)
		result = append(result, tx.MsgTx().TxOut[outPoint.Index])
	}
	return result
}

var zeroHash chainhash.Hash

func isCoinbase(tx *wire.MsgTx) bool {
	return tx.TxIn[0].PreviousOutPoint.Hash.IsEqual(&zeroHash)
}
