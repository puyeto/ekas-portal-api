package checkexpired

import (
	"fmt"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Status ...
type Status struct{}

// Run Status.Run() will get triggered automatically.
func (s Status) Run() {
	fmt.Println("Running ...")

	vDetails, err := getExpiredVehicles()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(vDetails))

	//
	// smsCheck, err := app.CheckMessages()
	// if err != nil {
	// 	return
	// }

}

func getExpiredVehicles() ([]models.VehicleDetails, error) {
	var vdetails = []models.VehicleDetails{}
	expdate := "DATE_ADD(DATE_ADD(COALESCE(renewal_date, vd.created_on), INTERVAL -1 DAY), INTERVAL 1 YEAR)"
	err := app.DBCon.Select("vehicle_id", "vehicle_reg_no", "owner_name", "owner_phone", "vehicle_status").
		From("vehicle_details AS vd").
		LeftJoin("vehicle_owner AS vo", dbx.NewExp("vd.owner_id = vo.owner_id_no")).
		Where(dbx.And(dbx.Or(dbx.NewExp(expdate+" < CURDATE()"),
			dbx.Between(expdate, "DATE_ADD(CURDATE(), INTERVAL -30 DAY)", "CURDATE()")), dbx.HashExp{"vehicle_status": 1})).All(&vdetails)
	return vdetails, err
}

func sendSMS(message, tonumber string) {
	app.MessageChan <- app.MessageDetails{
		Message:  message,
		ToNumber: tonumber,
		Type:     "Expired",
	}
}
