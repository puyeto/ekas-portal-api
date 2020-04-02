package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// check if messages have been sent

// SendSMSMessages ...
func SendViolationSMSMessages(message chan models.MessageDetails) {
	for {
		message := <-message
		toNumber := "0723436438"
		accessKey := "5Ad32ZMlklBaTOBqFc8Mtk5bwOl1r09j"
		apiKey := "8812492c587bf5cda4ee01a0bfedff3d"
		clientID := "97bcef17-e7b8-451b-a75c-913db0f2a069"
		urlString := "https://api.onfonmedia.co.ke/v1/sms/SendBulkSMS"

		if message.MessageID > 0 {
			// check if user has sent message in the last 5min
			data, _ := checkMessages(toNumber)
			t1 := time.Now()
			diff := t1.Sub(data.DateTime)
			dif := int64(diff.Minutes())

			if message.Message != data.Message {
				if dif > 5 && data.MessageID != message.MessageID {

					// Pack up the data for our message
					messageParameters, _ := json.Marshal(map[string]string{
						"Number": toNumber,
						"Text":   message.Message,
					})
					requestBody, _ := json.Marshal(map[string]string{
						"SenderId":          "EKAS",
						"MessageParameters": string(messageParameters),
						"ApiKey":            apiKey,
						"ClientId":          clientID,
					})

					// Create HTTP request client
					client := &http.Client{}
					req, _ := http.NewRequest("POST", urlString, bytes.NewBuffer(requestBody))
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("AccessKey", accessKey)

					// Make HTTP POST request and return message SID
					resp, _ := client.Do(req)
					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						var data map[string]interface{}
						decoder := json.NewDecoder(resp.Body)
						err := decoder.Decode(&data)
						if err == nil {
							// Save sent message
							var savedata = models.SaveMessageDetails{
								MessageID:   message.MessageID,
								Message:     data["body"].(string),
								DateTime:    time.Now(),
								From:        data["from"].(string),
								To:          data["to"].(string),
								DateCreated: data["date_created"].(string),
								SID:         data["sid"].(string),
								Status:      data["status"].(string),
							}
							saveSentMessages(savedata)
						}
					} else {
						fmt.Println(resp.Status)
					}
				}
			}

		}
	}
}

// SMSCheck ...
type SMSCheck struct {
	MessageID int       `json:"message_id"`
	DateTime  time.Time `json:"date_time"`
	Message   string    `json:"message"`
}

// check for sent messages
func checkViolationMessages(tonumber string) (SMSCheck, error) {
	var data SMSCheck
	q := DBCon.NewQuery("SELECT date_time, message_id, message FROM saved_messages WHERE `to`='" + tonumber + "' ORDER BY id DESC  LIMIT 1 ")
	err := q.One(&data)

	return data, err
}

// save sent messages
func saveSentMessages(m models.SaveMessageDetails) {
	DBCon.Insert("saved_messages", dbx.Params{
		"message_id":   m.MessageID,
		"message":      m.Message,
		"date_time":    m.DateTime,
		"from":         m.From,
		"to":           m.To,
		"date_created": m.DateCreated,
		"sid":          m.SID,
		"status":       m.Status,
	}).Execute()
}
