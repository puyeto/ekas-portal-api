package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// VehicleRecordDAO persists vehicleRecord data in database
type VehicleRecordDAO struct{}

// NewVehicleRecordDAO creates a new VehicleRecordDAO
func NewVehicleRecordDAO() *VehicleRecordDAO {
	return &VehicleRecordDAO{}
}

// Get reads the vehicleRecord with the specified ID from the database.
func (dao *VehicleRecordDAO) Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error) {
	var vehicleRecord models.VehicleDetails
	err := rs.Tx().Select().Model(id, &vehicleRecord)
	return &vehicleRecord, err
}

// Update saves the changes to an vehicleRecord in the database.
func (dao *VehicleRecordDAO) Update(rs app.RequestScope, id uint32, vehicleRecord *models.VehicleDetails) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	vehicleRecord.VehicleID = id
	return rs.Tx().Model(vehicleRecord).Exclude("Id").Update()
}

// Delete deletes an vehicleRecord with the specified ID from the database.
func (dao *VehicleRecordDAO) Delete(rs app.RequestScope, id uint32) error {
	vehicleRecord, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(vehicleRecord).Delete()
}

// Count returns the number of the vehicleRecord records in the database.
func (dao *VehicleRecordDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("vehicle_details").Row(&count)
	return count, err
}

// Query retrieves the vehicleRecord records with the specified offset and limit from the database.
func (dao *VehicleRecordDAO) Query(rs app.RequestScope, offset, limit int) ([]models.VehicleDetails, error) {
	vehicleRecords := []models.VehicleDetails{}
	err := rs.Tx().Select().OrderBy("created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	return vehicleRecords, err
}
