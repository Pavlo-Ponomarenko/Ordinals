package mint

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type Inscription struct {
	PrivateKey   *btcec.PrivateKey
	Script       []byte
	ControlBlock []byte
	PkScript     []byte
}

func createInscription(contentType string, content []byte, chainParams chaincfg.Params) (*Inscription, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	inscriptionBuilder := txscript.NewScriptBuilder().
		AddData(schnorr.SerializePubKey(privateKey.PubKey())).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_FALSE).
		AddOp(txscript.OP_IF).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte(contentType)).
		AddOp(txscript.OP_0).
		AddFullData(content).
		AddOp(txscript.OP_ENDIF)
	script, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	leaf := txscript.NewBaseTapLeaf(script)
	proof := &txscript.TapscriptProof{
		TapLeaf:  leaf,
		RootNode: leaf,
	}
	controlBlockData := proof.ToControlBlock(privateKey.PubKey())
	controlBlock, err := controlBlockData.ToBytes()
	if err != nil {
		return nil, err
	}
	tapHash := proof.RootNode.TapHash()
	commitOutputAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(privateKey.PubKey(), tapHash[:])), &chainParams)
	if err != nil {
		return nil, err
	}
	commitOutputScript, err := txscript.PayToAddrScript(commitOutputAddress)
	if err != nil {
		return nil, err
	}
	return &Inscription{
		PrivateKey:   privateKey,
		Script:       script,
		ControlBlock: controlBlock,
		PkScript:     commitOutputScript,
	}, nil
}
