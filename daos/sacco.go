package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// SaccoDAO persists sacco data in database
type SaccoDAO struct{}

// NewSaccoDAO creates a new SaccoDAO
func NewSaccoDAO() *SaccoDAO {
	return &SaccoDAO{}
}

// Get reads the sacco with the specified ID from the database.
func (dao *SaccoDAO) Get(rs app.RequestScope, id int) (*models.Saccos, error) {
	var sacco models.Saccos
	err := rs.Tx().Select("id", "name", "short_name", "address").
		Model(id, &sacco)
	return &sacco, err
}

// GetSaccoUser Get sacco details associated with a user
func (dao *SaccoDAO) GetSaccoUser(rs app.RequestScope, userid int) (*models.Saccos, error) {
	var sacco models.Saccos
	err := rs.Tx().Select("saccos.sacco_id", "sacco_name", "sacco_contacts", "sacco_contact_name", "sacco_email", "sacco_location", "user_id AS user", "contact_id", "COALESCE(sacco_phone, '') AS sacco_phone", "COALESCE(business_reg_no, '') AS business_reg_no").
		From("saccos").LeftJoin("sacco_users", dbx.NewExp("sacco_users.sacco_id = saccos.sacco_id")).
		Where(dbx.HashExp{"user_id": userid}).One(&sacco)
	return &sacco, err
}

// Create saves a new sacco record in the database.
// The Sacco.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *SaccoDAO) Create(rs app.RequestScope, sacco *models.Saccos) error {
	return rs.Tx().Model(sacco).Exclude("User").Insert()
}

// CreateSaccoUser create user relationship to sacco.
func (dao *SaccoDAO) CreateSaccoUser(rs app.RequestScope, saccoid int32, userid int32) error {
	_, err := rs.Tx().Insert("sacco_users", dbx.Params{
		"user_id":  userid,
		"sacco_id": saccoid}).Execute()
	return err
}

// Update saves the changes to an sacco in the database.
func (dao *SaccoDAO) Update(rs app.RequestScope, id int, sacco *models.Saccos) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	sacco.ID = int32(id)
	return rs.Tx().Model(sacco).Exclude("SaccoID", "User").Update()
}

// Delete deletes an sacco with the specified ID from the database.
func (dao *SaccoDAO) Delete(rs app.RequestScope, id int) error {
	sacco, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(sacco).Delete()
}

// Count returns the number of the sacco records in the database.
func (dao *SaccoDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("saccos").Row(&count)
	return count, err
}

// Query retrieves the sacco records with the specified offset and limit from the database.
func (dao *SaccoDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Saccos, error) {
	saccos := []models.Saccos{}
	err := rs.Tx().Select("id", "name", "short_name", "address").
		OrderBy("short_name ASC").Offset(int64(offset)).Limit(int64(limit)).All(&saccos)
	return saccos, err
}

// IsExistSaccoName Check if sacco name exists.
func (dao *SaccoDAO) IsExistSaccoName(rs app.RequestScope, sacconame string) (int, error) {
	var saccoid int
	q := rs.Tx().NewQuery("SELECT COALESCE(sacco_id, 0) AS sacco_id FROM saccos WHERE sacco_name='" + sacconame + "' LIMIT 1")
	q.Row(&saccoid)
	return saccoid, nil
}
