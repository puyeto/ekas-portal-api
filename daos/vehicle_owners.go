package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// ------------------------Add / Update Owner-----------------------------------

// CreateVehicleOwner Add vehicle owner
func (dao *VehicleDAO) CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) error {
	exists, _ := dao.VehicleOwnerExists(rs, vo.OwnerIDNo)
	if exists == 1 {
		return dao.UpdateVehicleOwners(rs, vo)
	}
	return rs.Tx().Model(vo).Exclude("OwnerID").Insert("UserID", "OwnerIDNo", "OwnerName", "OwnerEmail", "OwnerPhone")
}

// UpdateVehicleOwners ....
func (dao *VehicleDAO) UpdateVehicleOwners(rs app.RequestScope, vo *models.VehicleOwner) error {
	_, err := rs.Tx().Update("vehicle_owner", dbx.Params{
		"owner_name":  vo.OwnerName,
		"user_id":     vo.UserID,
		"owner_email": vo.OwnerEmail,
		"owner_phone": vo.OwnerPhone},
		dbx.HashExp{"owner_id_no": vo.OwnerIDNo}).Execute()
	return err
}

// VehicleOwnerExists ...
func (dao *VehicleDAO) VehicleOwnerExists(rs app.RequestScope, id string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_owner WHERE owner_id_no='" + id + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}
