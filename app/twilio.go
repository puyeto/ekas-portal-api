package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ekas-portal-api/models"
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
		// Set account keys & information
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
				fmt.Println(data)
				// Save sent message
			}
		} else {
			fmt.Println(resp.Status)
		}
	}
}
