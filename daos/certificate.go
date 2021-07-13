package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// CertificateDAO persists certificate data in database
type CertificateDAO struct{}

// NewCertificateDAO creates a new CertificateDAO
func NewCertificateDAO() *CertificateDAO {
	return &CertificateDAO{}
}

// Get reads the certificate with the specified ID from the database.
func (dao *CertificateDAO) Get(rs app.RequestScope, id int) (*models.Certificates, error) {
	var certificate models.Certificates
	err := rs.Tx().Select().Model(id, &certificate)
	return &certificate, err
}

// Create saves a new certificate record in the database.
// The Certificate.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *CertificateDAO) Create(rs app.RequestScope, certificate *models.Certificates) error {
	return rs.Tx().Model(certificate).Exclude("Company", "Fitter").Insert()
}

// Update saves the changes to an certificate in the database.
func (dao *CertificateDAO) Update(rs app.RequestScope, id int, certificate *models.Certificates) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	certificate.ID = id
	return rs.Tx().Model(certificate).Exclude("Id").Update()
}

// Delete deletes an certificate with the specified ID from the database.
func (dao *CertificateDAO) Delete(rs app.RequestScope, id int) error {
	certificate, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(certificate).Delete()
}

// Count returns the number of the certificate records in the database.
func (dao *CertificateDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("certificates").Row(&count)
	return count, err
}

// Query retrieves the certificate records with the specified offset and limit from the database.
func (dao *CertificateDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Certificates, error) {
	certificates := []models.Certificates{}
	err := rs.Tx().Select("ce.id", "ce.company_id", "issued_on", "cert_no", "cert_serial", "created_on",
		"fitter_id", "company_name AS company", "CONCAT(first_name , ' ', last_name) AS fitter").From("certificates AS ce").
		LeftJoin("companies", dbx.NewExp("companies.company_id = ce.company_id")).
		LeftJoin("auth_users", dbx.NewExp("auth_users.auth_user_id = ce.fitter_id")).
		OrderBy("id").Offset(int64(offset)).Limit(int64(limit)).All(&certificates)
	return certificates, err
}
