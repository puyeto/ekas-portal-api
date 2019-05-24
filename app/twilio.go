package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Message ...
var Message = make(chan models.MessageDetails)

func init() {
	// message := <-Message
	go SendSMSMessages(Message)
}

// check if messages have been sent

// SendSMSMessages ...
func SendSMSMessages(message chan models.MessageDetails) {
	for {
		message := <-message
		toNumber := "+254723436438"
		// Set account keys & information

		if message.MessageID > 0 {
			data, _ := checkSentMessages(toNumber)
			t1 := time.Now()
			diff := t1.Sub(data.DateTime)
			dif := int64(diff.Minutes())
			fmt.Println(dif, data.MessageID)
			if (dif > 10 && data.MessageID == message.MessageID) || (data.MessageID != message.MessageID) {
				accountSid := "ACeab16ebd80a48c1f4318f09c6ad6e33e"
				authToken := "8812492c587bf5cda4ee01a0bfedff3d"
				urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
				fmt.Println("accountSid = ", accountSid)

				// Pack up the data for our message
				msgData := url.Values{}
				msgData.Set("To", "+254723436438")
				msgData.Set("From", "+14086101380")
				msgData.Set("Body", message.Message)
				msgDataReader := *strings.NewReader(msgData.Encode())

				// Create HTTP request client
				client := &http.Client{}
				req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
				req.SetBasicAuth(accountSid, authToken)
				req.Header.Add("Accept", "application/json")
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

// SMSCheck ...
type SMSCheck struct {
	MessageID int       `json:"message_id"`
	DateTime  time.Time `json:"date_time"`
}

// check if user has sent message in the last 5min
func checkSentMessages(tonumber string) (SMSCheck, error) {
	var data SMSCheck
	q := DBCon.NewQuery("SELECT date_time, message_id FROM saved_messages WHERE `to`='" + tonumber + "' ORDER BY id DESC  LIMIT 1 ")
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
