package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// SimcardDAO persists simcard data in database
type SimcardDAO struct{}

// NewSimcardDAO creates a new SimcardDAO
func NewSimcardDAO() *SimcardDAO {
	return &SimcardDAO{}
}

// Get reads the simcard with the specified ID from the database.
func (dao *SimcardDAO) Get(rs app.RequestScope, id int) (*models.Simcards, error) {
	var simcard models.Simcards
	err := rs.Tx().Select("*").Model(id, &simcard)
	return &simcard, err
}

// Create saves a new simcard record in the database.
// The Simcard.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *SimcardDAO) Create(rs app.RequestScope, simcard *models.Simcards) error {
	return rs.Tx().Model(simcard).Exclude("User").Insert()
}

// Update saves the changes to an simcard in the database.
func (dao *SimcardDAO) Update(rs app.RequestScope, id int, simcard *models.Simcards) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	simcard.ID = uint32(id)
	return rs.Tx().Model(simcard).Exclude("ID").Update()
}

// Delete deletes an simcard with the specified ID from the database.
func (dao *SimcardDAO) Delete(rs app.RequestScope, id int) error {
	simcard, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(simcard).Delete()
}

// Count returns the number of the simcard records in the database.
func (dao *SimcardDAO) Count(rs app.RequestScope, status string) (int, error) {
	var count int
	q := rs.Tx().Select("COUNT(*)").From("simcards")
	if status == "managed" {
		q = q.Where(dbx.HashExp{"managed": 1})
	} else if status == "inactive" {
		q = q.Where(dbx.HashExp{"managed": 0})
	}
	err := q.Row(&count)
	return count, err
}

// Query retrieves the simcard records with the specified offset and limit from the database.
func (dao *SimcardDAO) Query(rs app.RequestScope, offset, limit int, status string) ([]models.Simcards, error) {
	simcards := []models.Simcards{}
	q := rs.Tx().Select()
	if status == "managed" {
		q = q.Where(dbx.HashExp{"managed": 1})
	} else if status == "inactive" {
		q = q.Where(dbx.HashExp{"managed": 0})
	}
	err := q.OrderBy("identifier ASC").Offset(int64(offset)).Limit(int64(limit)).All(&simcards)
	return simcards, err
}

// IsExistSimcard Check if simcard identifier exists.
func (dao *SimcardDAO) IsExistSimcard(rs app.RequestScope, identifier string) (int, error) {
	var simcardid int
	q := rs.Tx().NewQuery("SELECT id FROM simcards WHERE identifier='" + identifier + "' LIMIT 1")
	q.Row(&simcardid)
	return simcardid, nil
}

// GetStats ...
func (dao *SimcardDAO) GetStats(rs app.RequestScope) *models.SimcardStats {
	var stats models.SimcardStats
	rs.Tx().Select("COUNT(*)").From("simcards").Row(&stats.TotalCount)
	rs.Tx().Select("COUNT(*)").Where(dbx.HashExp{"managed": 1}).From("simcards").Row(&stats.ManagedCount)
	stats.UnManagedCount = stats.TotalCount - stats.ManagedCount
	return &stats
}
