package daos

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// ------------------------Add / Update Owner-----------------------------------

// CreateVehicleOwner Add vehicle owner
func (dao *VehicleDAO) CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) error {
	exists, _ := dao.VehicleOwnerExists(rs, vo.OwnerID)
	if exists == 1 {
		return dao.UpdateVehicleOwners(rs, vo)
	}
	return rs.Tx().Model(vo).Insert("OwnerID", "UserID", "OwnerIDNo", "OwnerName", "OwnerEmail", "OwnerPhone")
}

// UpdateVehicleOwners ....
func (dao *VehicleDAO) UpdateVehicleOwners(rs app.RequestScope, vo *models.VehicleOwner) error {
	_, err := rs.Tx().Update("vehicle_owner", dbx.Params{
		"owner_name":  vo.OwnerName,
		"user_id":     vo.UserID,
		"owner_email": vo.OwnerEmail,
		"owner_phone": vo.OwnerPhone},
		dbx.HashExp{"owner_id": vo.OwnerID}).Execute()
	return err
}

// VehicleOwnerExists ...
func (dao *VehicleDAO) VehicleOwnerExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_owner WHERE owner_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}
