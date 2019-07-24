package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// CompanyDAO persists company data in database
type CompanyDAO struct{}

// NewCompanyDAO creates a new CompanyDAO
func NewCompanyDAO() *CompanyDAO {
	return &CompanyDAO{}
}

// Get reads the company with the specified ID from the database.
func (dao *CompanyDAO) Get(rs app.RequestScope, id int) (*models.Companies, error) {
	var company models.Companies
	err := rs.Tx().Select().Model(id, &company)
	return &company, err
}

// Create saves a new company record in the database.
// The Company.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *CompanyDAO) Create(rs app.RequestScope, company *models.Companies) error {
	return rs.Tx().Model(company).Exclude("UserID").Insert()
}

// Update saves the changes to an company in the database.
func (dao *CompanyDAO) Update(rs app.RequestScope, id int, company *models.Companies) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	company.CompanyID = id
	return rs.Tx().Model(company).Exclude("CompnayID, UserID").Update()
}

// Delete deletes an company with the specified ID from the database.
func (dao *CompanyDAO) Delete(rs app.RequestScope, id int) error {
	company, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(company).Delete()
}

// Count returns the number of the company records in the database.
func (dao *CompanyDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("companies").Row(&count)
	return count, err
}

// Query retrieves the company records with the specified offset and limit from the database.
func (dao *CompanyDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Companies, error) {
	companys := []models.Companies{}
	err := rs.Tx().Select().OrderBy("company_id").Offset(int64(offset)).Limit(int64(limit)).All(&companys)
	return companys, err
}
