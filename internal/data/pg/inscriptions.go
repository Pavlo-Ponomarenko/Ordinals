package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"ordinals/internal/data"
)

type inscriptionQ struct {
	db                   *pgdb.DB
	inscriptionSQL       sq.SelectBuilder
	addressSQL           sq.SelectBuilder
	inscriptionSQLInsert sq.InsertBuilder
	addressSQLInsert     sq.InsertBuilder
}

const inscriptionsTable = "inscription"
const addressesTable = "address"

func NewInscriptionQ(db *pgdb.DB) data.InscriptionQ {
	return &inscriptionQ{
		db:                   db.Clone(),
		inscriptionSQL:       sq.Select("*").From(inscriptionsTable),
		addressSQL:           sq.Select("*").From(addressesTable),
		inscriptionSQLInsert: sq.Insert(inscriptionsTable),
		addressSQLInsert:     sq.Insert(addressesTable),
	}
}

func (q *inscriptionQ) New() data.InscriptionQ {
	return NewInscriptionQ(q.db)
}

func (q *inscriptionQ) GetInscriptions(address string) ([]data.InscriptionEntity, error) {
	var result []data.InscriptionEntity
	q.inscriptionSQL = q.inscriptionSQL.Where(sq.Eq{"address_id": address})
	err := q.db.Select(&result, q.inscriptionSQL)
	return result, err
}

func (q *inscriptionQ) SaveInscription(inscription *data.InscriptionEntity) error {
	var addresses data.AddressEntity
	q.addressSQL = q.addressSQL.Where(sq.Eq{"id": inscription.AddressId})
	err := q.db.Get(&addresses, q.addressSQL)
	if err == sql.ErrNoRows {
		err = q.SaveAddress(inscription.AddressId)
		if err != nil {
			fmt.Println(q.addressSQL.ToSql())
			return err
		}
	}
	jsonStr, _ := json.Marshal(*inscription)
	var clauses map[string]interface{}
	json.Unmarshal(jsonStr, &clauses)
	q.inscriptionSQLInsert = q.inscriptionSQLInsert.SetMap(clauses)
	err = q.db.Exec(q.inscriptionSQLInsert)
	if err != nil {
		return err
	}
	return nil
}

func (q *inscriptionQ) SaveAddress(address string) error {
	clauses := make(map[string]interface{})
	clauses["id"] = address
	q.addressSQLInsert = q.addressSQLInsert.SetMap(clauses)
	err := q.db.Exec(q.addressSQLInsert)
	return err
}
