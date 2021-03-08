package app

import (
	"fmt"

	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// MessageChan ...
var MessageChan = make(chan MessageDetails)

// MessageDetails ...
type MessageDetails struct {
	MessageID string
	Message   string
	ToNumber  string
	Type      string
}

func init() {
	// message := <-Message
	go SendSMSMessages(MessageChan)
}

// SendSMSMessages ...
func SendSMSMessages(message chan MessageDetails) {
	for {
		message := <-message
		//Call the Gateway, and pass the constants here!
		smsService := NewSMSService(Config.ATAPIUsername, Config.ATAPIKey, "production")

		//Send SMS - REPLACE Recipient and Message with REAL Values
		recipients, err := smsService.Send("EKASTECH", message.ToNumber, message.Message) //Leave blank, "", if you don't have one)
		if err != nil {
			fmt.Println(err)
		}

		// save sms

		fmt.Println(recipients)
	}

}

// CheckMessages check for sent messages
func CheckMessages(tonumber, messagetype string) (SMSCheck, error) {
	var data SMSCheck
	q := DBCon.NewQuery("SELECT date_time FROM saved_messages WHERE `to`='" + tonumber + "' AND `message_type`='" + messagetype + "' ORDER BY id DESC  LIMIT 1 ")
	err := q.One(&data)

	return data, err
}

// SaveSentMessages save sent messages
func SaveSentMessages(m models.SaveMessageDetails) {
	_, err := DBCon.Insert("saved_messages", dbx.Params{
		"message_id":   m.MessageID,
		"vehicle_id":   m.VehicleID,
		"device_id":    m.DeviceID,
		"message":      m.Message,
		"message_type": m.MessageType,
		"date_time":    m.DateTime,
		"from":         m.From,
		"to":           m.To,
		"status":       m.Status,
	}).Execute()

	fmt.Println(err)
}
