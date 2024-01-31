package data

type InscriptionQ interface {
	New() InscriptionQ

	GetInscriptions(address string) ([]InscriptionEntity, error)
	SaveInscription(inscription *InscriptionEntity) error
	SaveAddress(address string) error
}
