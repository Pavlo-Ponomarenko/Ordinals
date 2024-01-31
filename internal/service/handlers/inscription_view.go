package handlers

import "ordinals/internal/data"

type InscriptionView struct {
	Id      int64
	Content string
	TxHash  string
	Date    string
}

func InscriptionEntityToView(entity data.InscriptionEntity) InscriptionView {
	return InscriptionView{
		Id:      entity.Id,
		Content: entity.Content,
		TxHash:  entity.TxHash,
		Date:    entity.Date,
	}
}

func InscriptionEntitiesToViews(entities []data.InscriptionEntity) []InscriptionView {
	views := make([]InscriptionView, 0)
	for _, entity := range entities {
		views = append(views, InscriptionEntityToView(entity))
	}
	return views
}
