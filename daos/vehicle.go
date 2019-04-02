package daos

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// VehicleDAO persists vehicle data in database
type VehicleDAO struct{}

// NewVehicleDAO creates a new VehicleDAO
func NewVehicleDAO() *VehicleDAO {
	return &VehicleDAO{}
}

// GetVehicleByStrID ...
func (dao *VehicleDAO) GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error) {
	var vdetails models.VehicleConfigDetails
	q := rs.Tx().NewQuery("SELECT conf_id, vehicle_id, owner_id, fitter_id, data FROM vehicle_configuration WHERE vehicle_string_id='" + strid + "' LIMIT 1")
	err := q.Row(&vdetails.ConfigID, &vdetails.VehicleID, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.Data)

	return &vdetails, err
}

// CreateVehicle saves a new vehicle record in the database.
func (dao *VehicleDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	if exists, _ := dao.VehicleExists(rs, v.VehicleID); exists == 1 {
		return dao.UpdateVehicle(rs, v)
	}

	return rs.Tx().Model(v).Insert("VehicleID", "VehicleStringID", "VehicleRegNo", "ChassisNo", "MakeType")
}

// UpdateVehicle ....
func (dao *VehicleDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	_, err := rs.Tx().Update("vehicle_details", dbx.Params{
		"vehicle_string_id": v.VehicleStringID,
		"vehicle_reg_no":    v.VehicleRegNo,
		"chassis_no":        v.ChassisNo,
		"make_type":         v.MakeType},
		dbx.HashExp{"vehicle_string_id": v.VehicleID}).Execute()
	return err
}

// VehicleExists ...
func (dao *VehicleDAO) VehicleExists(rs app.RequestScope, id uint64) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// CreateVehicleOwner Add vehicle owner
func (dao *VehicleDAO) CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) error {
	if exists, _ := dao.VehicleOwnerExists(rs, vo.OwnerID); exists == 1 {
		return dao.UpdateVehicleOwners(rs, vo)
	}
	return rs.Tx().Model(vo).Insert("OwnerID", "OwnerIDNo", "OwnerName", "OwnerEmail", "OwnerPhone")

}

// UpdateVehicleOwners ....
func (dao *VehicleDAO) UpdateVehicleOwners(rs app.RequestScope, vo *models.VehicleOwner) error {
	_, err := rs.Tx().Update("vehicle_owner", dbx.Params{
		"owner_name":  vo.OwnerName,
		"owner_email": vo.OwnerEmail,
		"owner_phone": vo.OwnerPhone},
		dbx.HashExp{"owner_id": vo.OwnerID}).Execute()
	return err
}

// VehicleOwnerExists check if a driver owner exists
func (dao *VehicleDAO) VehicleOwnerExists(rs app.RequestScope, id uint64) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_owner WHERE owner_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// CreateFitter Add Fitter
func (dao *VehicleDAO) CreateFitter(rs app.RequestScope, fd *models.FitterDetails) error {
	if exists, _ := dao.FitterExists(rs, fd.FitterID); exists == 1 {
		return dao.UpdateFitter(rs, fd)
	}
	return rs.Tx().Model(fd).Insert("FitterID", "FitterIDNo", "FittingCenterName", "FitterLocation", "FitterEmail", "FitterAddress", "FitterPhone", "FittingDate", "FitterBizRegNo")
}

// UpdateFitter update fitter
func (dao *VehicleDAO) UpdateFitter(rs app.RequestScope, fd *models.FitterDetails) error {
	_, err := rs.Tx().Update("fitter_details", dbx.Params{
		"fitting_center_name": fd.FittingCenterName,
		"fitter_location":     fd.FitterLocation,
		"fitter_email":        fd.FitterEmail,
		"fitter_address":      fd.FitterAddress,
		"fitter_phone":        fd.FitterPhone,
		"fitting_date":        fd.FittingDate,
		"fitter_biz_reg_no":   fd.FitterBizRegNo},
		dbx.HashExp{"fitting_id": fd.FitterID}).Execute()
	return err
}

// FitterExists check if fitter exists
func (dao *VehicleDAO) FitterExists(rs app.RequestScope, id uint64) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM fitter_details WHERE fitting_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

//CreateConfiguration Add configuartion details to db
func (dao *VehicleDAO) CreateConfiguration(rs app.RequestScope, cd *models.Vehicle, ownerid uint64, fitterid uint64, vehicleid uint64) error {
	a, _ := json.Marshal(cd)
	_, err := rs.Tx().Insert("vehicle_configuration", dbx.Params{
		"conf_id":           app.GenerateNewID(),
		"vehicle_id":        vehicleid,
		"owner_id":          ownerid,
		"fitter_id":         fitterid,
		"vehicle_string_id": strings.ToLower(strings.Replace(cd.DeviceDetails.RegistrationNO, " ", "", -1)),
		"data":              string(a)}).Execute()
	return err
}
