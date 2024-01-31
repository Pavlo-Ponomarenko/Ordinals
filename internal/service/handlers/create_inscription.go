package handlers

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"net/http"
	"ordinals/internal/data"
	"ordinals/internal/mint"
	"ordinals/internal/ord_math"
	"ordinals/internal/service/requests"
	"time"
)

func CreateInscription(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	inscriptionData, err := requests.NewCreateInscriptionRequest(r)
	if err != nil {
		log.Debug(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entity := data.InscriptionRequestToEntity(*inscriptionData)
	client := BtcdClient(r)
	commitTx, revealTx, err := mint.FormTransactions(client, chaincfg.TestNet3Params, inscriptionData.Hash, inscriptionData.Output, inscriptionData.Key, "text", inscriptionData.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = client.SendRawTransaction(commitTx, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Debug("Commit tx was not accepted")
		return
	}
	revealTxHash, err := client.SendRawTransaction(revealTx, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Debug("Reveal tx was not accepted")
		return
	}
	log.Debug("Reveal tx hash: ", revealTxHash.String())
	go saveAfterConfirmation(client, r, revealTx, entity)
	w.WriteHeader(http.StatusOK)
}

var blockPeriod = 10

func saveAfterConfirmation(client *rpcclient.Client, r *http.Request, tx *wire.MsgTx, entity *data.InscriptionEntity) {
	log := Log(r)
	txHash := tx.TxHash()
	for i := 0; i < 6; i++ {
		time.Sleep(time.Duration(blockPeriod) * time.Minute)
		txResult, err := client.GetRawTransactionVerbose(&txHash)
		if err != nil {
			log.Debug(err)
			return
		}
		if txResult.Confirmations == 0 {
			continue
		}
		log.Debug("Tx ", txHash, " succeeded")
		entity.Id = ord_math.FindOutputOrder(client, tx, 0)
		q := InscriptionQ(r)
		err = q.SaveInscription(entity)
		if err != nil {
			log.Debug(err)
			return
		}
		log.Debug("Inscription ", entity.Id, " saved")
		return
	}
	log.Debug("Tx ", txHash, " failed")
}
