package data

import (
	"ordinals/internal/service/requests"
	"time"
)

type InscriptionEntity struct {
	Id        int64  `db:"id" json:"id"`
	Content   string `db:"content" json:"content"`
	TxHash    string `db:"tx_hash" json:"tx_hash"`
	Date      string `db:"date" json:"date"`
	AddressId string `db:"address_id" json:"address_id"`
}

func InscriptionRequestToEntity(request requests.CreateInscriptionRequest) *InscriptionEntity {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return &InscriptionEntity{
		Content:   request.Content,
		TxHash:    request.Hash,
		Date:      formattedTime,
		AddressId: request.Address,
	}
}
