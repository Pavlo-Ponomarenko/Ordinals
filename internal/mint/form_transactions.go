package mint

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func FormTransactions(contentType string, content string, chainParams chaincfg.Params) (firstTx string, secondTx string, err error) {
	inscription, err := createInscription(contentType, []byte(content), chainParams)
	if err != nil {
		return
	}
	sourceTxHash, _ := chainhash.NewHashFromStr("18782419d37ca0a356e32a2d1d7da04466bb3f82337adb2399e1ee2690c6b9e8")
	coinSource := &wire.OutPoint{
		Hash:  *sourceTxHash,
		Index: 1,
	}
	privateKeyBytes, _ := hex.DecodeString("4aeadd94b538ab69dc2ae9136aeb0722211813d82fb2907a9d1324942be90f61")
	privateKey, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	pk := publicKey.SerializeCompressed()
	p2wkhAddr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pk), &chainParams)
	fmt.Println(p2wkhAddr.String())
	returnToSourceScript, _ := txscript.PayToAddrScript(p2wkhAddr)
	commitTx, err := createCommitTransaction(p2wkhAddr, privateKey, coinSource, 798700, 400, 250, returnToSourceScript, inscription.PkScript)
	fmt.Println(err)
	firstTx = hexifyTx(commitTx)
	revealTx, err := createRevealTransaction(*commitTx, inscription.PrivateKey, inscription.Script, inscription.ControlBlock, 200, returnToSourceScript)
	secondTx = hexifyTx(revealTx)
	return
}

func hexifyTx(tx *wire.MsgTx) string {
	var buf bytes.Buffer
	tx.Serialize(&buf)
	return hex.EncodeToString(buf.Bytes())
}
