package mint

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func FormTransactions(client *rpcclient.Client, chainParams chaincfg.Params, sourceTxHashStr string, output uint32, privateKeyStr string, contentType string, content string) (commitTx *wire.MsgTx, revealTx *wire.MsgTx, err error) {
	inscription, err := createInscription(contentType, []byte(content), chainParams)
	if err != nil {
		return
	}
	sourceTxHash, _ := chainhash.NewHashFromStr(sourceTxHashStr)
	coinSource := &wire.OutPoint{
		Hash:  *sourceTxHash,
		Index: output,
	}
	sourceTx, _ := client.GetRawTransaction(sourceTxHash)
	inputAmount := sourceTx.MsgTx().TxOut[output].Value
	privateKeyBytes, _ := hex.DecodeString(privateKeyStr)
	privateKey, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	pk := publicKey.SerializeCompressed()
	p2wkhAddr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pk), &chainParams)
	returnToSourceScript, _ := txscript.PayToAddrScript(p2wkhAddr)
	commitTx, err = createCommitTransaction(p2wkhAddr, privateKey, coinSource, inputAmount, 251, 250, returnToSourceScript, inscription.PkScript)
	if err != nil {
		return
	}
	revealTx, err = createRevealTransaction(*commitTx, inscription.PrivateKey, inscription.Script, inscription.ControlBlock, 250, returnToSourceScript)
	return
}

func hexifyTx(tx *wire.MsgTx) string {
	var buf bytes.Buffer
	tx.Serialize(&buf)
	return hex.EncodeToString(buf.Bytes())
}
