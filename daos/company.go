package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
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
	err := rs.Tx().Select("company_id", "company_name", "company_contacts", "company_contact_name", "company_email", "company_location", "updated_by", "contact_id", "COALESCE(company_phone, '') AS company_phone", "COALESCE(business_reg_no, '') AS business_reg_no").
		Model(id, &company)
	return &company, err
}

// GetCompanyUser Get company details associated with a user
func (dao *CompanyDAO) GetCompanyUser(rs app.RequestScope, userid int) (*models.Companies, error) {
	var company models.Companies
	err := rs.Tx().Select("companies.company_id", "company_name", "company_contacts", "company_contact_name", "company_email", "company_location", "user_id AS user", "contact_id", "COALESCE(company_phone, '') AS company_phone", "COALESCE(business_reg_no, '') AS business_reg_no").
		From("companies").LeftJoin("company_users", dbx.NewExp("company_users.company_id = companies.company_id")).
		Where(dbx.HashExp{"user_id": userid}).One(&company)
	return &company, err
}

// Create saves a new company record in the database.
// The Company.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *CompanyDAO) Create(rs app.RequestScope, company *models.Companies) error {
	return rs.Tx().Model(company).Exclude("User").Insert()
}

// CreateCompanyUser create user relationship to company.
func (dao *CompanyDAO) CreateCompanyUser(rs app.RequestScope, companyid int32, userid int32) error {
	_, err := rs.Tx().Insert("company_users", dbx.Params{
		"user_id":    userid,
		"company_id": companyid}).Execute()
	return err
}

// Update saves the changes to an company in the database.
func (dao *CompanyDAO) Update(rs app.RequestScope, id int, company *models.Companies) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	company.CompanyID = int32(id)
	return rs.Tx().Model(company).Exclude("CompanyID", "User").Update()
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
	err := rs.Tx().Select("company_id", "company_name", "company_contacts", "company_contact_name", "company_email", "company_location", "contact_id", "COALESCE(company_phone, '') AS company_phone", "COALESCE(business_reg_no, '') AS business_reg_no").
		OrderBy("company_name ASC").Offset(int64(offset)).Limit(int64(limit)).All(&companys)
	return companys, err
}

// IsExistCompanyName Check if company name exists.
func (dao *CompanyDAO) IsExistCompanyName(rs app.RequestScope, companyname string) (int, error) {
	var companyid int
	q := rs.Tx().NewQuery("SELECT company_id FROM companies WHERE company_name='" + companyname + "' LIMIT 1")
	err := q.Row(&companyid)
	return companyid, err
}
