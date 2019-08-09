package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// OwnerDAO persists owner data in database
type OwnerDAO struct{}

// NewOwnerDAO creates a new OwnerDAO
func NewOwnerDAO() *OwnerDAO {
	return &OwnerDAO{}
}

// Get reads the owner with the specified ID from the database.
func (dao *OwnerDAO) Get(rs app.RequestScope, id uint32) (*models.VehicleOwner, error) {
	var owner models.VehicleOwner
	err := rs.Tx().Select().Model(id, &owner)
	return &owner, err
}

// Create saves a new owner record in the database.
// The Owner.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *OwnerDAO) Create(rs app.RequestScope, owner *models.VehicleOwner) error {
	owner.OwnerID = 0
	return rs.Tx().Model(owner).Insert()
}

// Update saves the changes to an owner in the database.
func (dao *OwnerDAO) Update(rs app.RequestScope, id uint32, owner *models.VehicleOwner) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	owner.OwnerID = id
	return rs.Tx().Model(owner).Exclude("Id").Update()
}

// Delete deletes an owner with the specified ID from the database.
func (dao *OwnerDAO) Delete(rs app.RequestScope, id uint32) error {
	owner, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(owner).Delete()
}

// Count returns the number of the owner records in the database.
func (dao *OwnerDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("vehicle_owner").Row(&count)
	return count, err
}

// Query retrieves the owner records with the specified offset and limit from the database.
func (dao *OwnerDAO) Query(rs app.RequestScope, offset, limit int) ([]models.VehicleOwner, error) {
	owners := []models.VehicleOwner{}
	err := rs.Tx().Select().OrderBy("owner_id").Offset(int64(offset)).Limit(int64(limit)).All(&owners)
	return owners, err
}
