package revetdeviceids

import (
	"github.com/ekas-portal-api/app"
)

// Status ...
type Status struct {
	// filtered
}

// Run REVERT DEVICE DATA TO ORIGINAL CONFIGURED IDs
func (c Status) Run() {
	query := "UPDATE IGNORE `vehicle_configuration` SET `device_id`=json_value(DATA, '$.governor_details.device_id')"
	app.DBCon.NewQuery(query).Execute()
}
