package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// ---------------------Add / Update Fitter--------------------------------------

// CreateFitter Add Fitter
func (dao *VehicleDAO) CreateFitter(rs app.RequestScope, fd *models.FitterDetails) error {
	exists, _ := dao.FitterExists(rs, fd.FitterIDNo)
	if exists == 1 {
		return dao.UpdateFitter(rs, fd)
	}
	return rs.Tx().Model(fd).Insert("FitterID", "UserID", "FitterIDNo", "FittingCenterName", "FitterLocation", "FitterEmail", "FitterAddress", "FitterPhone", "FittingDate", "FitterBizRegNo")
}

// UpdateFitter update fitter
func (dao *VehicleDAO) UpdateFitter(rs app.RequestScope, fd *models.FitterDetails) error {
	_, err := rs.Tx().Update("fitter_details", dbx.Params{
		"fitting_center_name": fd.FittingCenterName,
		"user_id":             fd.UserID,
		"fitter_location":     fd.FitterLocation,
		"fitter_email":        fd.FitterEmail,
		"fitter_address":      fd.FitterAddress,
		"fitter_phone":        fd.FitterPhone,
		"fitting_date":        fd.FittingDate,
		"fitter_biz_reg_no":   fd.FitterBizRegNo},
		dbx.HashExp{"fitting_id_no": fd.FitterIDNo}).Execute()
	return err
}

// FitterExists check if fitter exists
func (dao *VehicleDAO) FitterExists(rs app.RequestScope, id string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM fitter_details WHERE fitting_id_no='" + id + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}
