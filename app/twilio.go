package app

import "github.com/sfreiberg/gotwilio"

var (
	twilio *gotwilio.Twilio
)

// InitializeTwilio ...
func InitializeTwilio() {
	twilio = gotwilio.NewTwilioClient(Config.TwilioAccountSID, Config.TwilioAuthToken)
}

// SendMessage ...
func SendMessage(from, to, message string) {
	twilio.SendSMS(from, to, message, "", "")
}
