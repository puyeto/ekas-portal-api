package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ekas-portal-api/models"
)

// Env is the environment type
type Env string

const (
	// DEV is the development env tag

	// SANDBOX is the sandbox env tag
	SANDBOX = iota
	// PRODUCTION is the production env tag
	PRODUCTION

	// APPKEY ...
	// APPKEY      = "ZpBbYCDc4ZzoJjOsGo5nTcdI4r2JCSBY"
	// APPSECRET   = "e24DmwdmMhHr52Hh"
	// CALLBACKURL = "https://618a574c.ngrok.io"
	// SHORTCODE   = "174379"
	// PASSWORD    = "MTc0Mzc5YmZiMjc5ZjlhYTliZGJjZjE1OGU5N2RkNzFhNDY3Y2QyZTBjODkzMDU5YjEwZjc4ZTZiNzJhZGExZWQyYzkxOTIwMTkwNjA3MTIwNTMw"
	// TIMESTAMP   = "20190607120530"

	APPKEY      = "V5LFV9ArzihjpjHkKg0BToA8GEhdM7vD"
	APPSECRET   = "cjIWCUTjHJNjgTFh"
	CALLBACKURL = "https://618a574c.ngrok.io"
	SHORTCODE   = "278514"
	PASSWORD    = "Mjc4NTE0ZWM2YWMzMmY2ZTFlMzY0MTE3NTEwM2I4ZDYwYTQxOWRlNDg0OWYwMWVkMTA3M2U1Y2RjOTY4MWE3Y2NkYzIxOTIwMjIwNzE0MTIwNTMw"
	TIMESTAMP   = "20220714120530"
	PASSKEY     = "ec6ac32f6e1e3641175103b8d60a419de4849f01ed1073e5cdc9681a7ccdc219"
)

// Service is an Mpesa Service
type Service struct {
	AppKey    string
	AppSecret string
	Env       int
}

// New return a new Mpesa Service
func New(appKey, appSecret string, env int) (Service, error) {
	return Service{appKey, appSecret, env}, nil
}

//Generate Mpesa Daraja Access Token
func (s Service) auth() (string, error) {
	b := []byte(s.AppKey + ":" + s.AppSecret)
	encoded := base64.StdEncoding.EncodeToString(b)

	url := s.baseURL() + "oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(encoded))
	if err != nil {
		return "", err
	}
	req.Header.Add("authorization", "Basic "+encoded)
	req.Header.Add("cache-control", "no-cache")

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("could not send auth request: %v", err)
	}

	var authResponse models.AuthResponse
	err = json.NewDecoder(res.Body).Decode(&authResponse)
	if err != nil {
		return "", fmt.Errorf("could not decode auth response: %v", err)
	}

	accessToken := authResponse.AccessToken
	fmt.Println("MPESA TOKEN ", accessToken)
	return accessToken, nil
}

// Simulation requests user device for payment
func (s Service) Simulation(express models.Express) (string, error) {
	body, err := json.Marshal(express)
	if err != nil {
		return "", nil
	}
	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["content-type"] = "application/json"
	headers["authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/stkpush/v1/processrequest"
	return s.newReq(url, body, headers)
}

// TransactionStatus gets status of a transaction
func (s Service) TransactionStatus(express models.Status) (string, error) {
	body, err := json.Marshal(express)
	if err != nil {
		return "", nil
	}

	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth

	url := s.baseURL() + "mpesa/stkpushquery/v1/query"
	return s.newReq(url, body, headers)
}

// C2BRegisterURL requests
func (s Service) C2BRegisterURL(c2bRegisterURL models.C2BRegisterURL) (string, error) {
	body, err := json.Marshal(c2bRegisterURL)
	if err != nil {
		return "", err
	}

	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["Cache-Control"] = "no-cache"

	url := s.baseURL() + "mpesa/c2b/v1/registerurl"
	return s.newReq(url, body, headers)
}

// C2BSimulation sends a new request
func (s Service) C2BSimulation(c2b models.C2B) ([]byte, error) {
	body, err := json.Marshal(c2b)
	if err != nil {
		return nil, err
	}

	auth, err := s.auth()
	if err != nil {
		return nil, nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/c2b/v1/simulate"
	return s.newReq2(url, body, headers)
}

// B2CRequest sends a new request
func (s Service) B2CRequest(b2c models.B2C) (string, error) {
	body, err := json.Marshal(b2c)
	if err != nil {
		return "", err
	}

	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/b2c/v1/paymentrequest"
	return s.newReq(url, body, headers)
}

// B2BRequest sends a new request
func (s Service) B2BRequest(b2b models.B2B) (string, error) {
	body, err := json.Marshal(b2b)
	if err != nil {
		return "", nil
	}
	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/b2b/v1/paymentrequest"
	return s.newReq(url, body, headers)
}

// Reversal requests a reversal?
func (s Service) Reversal(reversal models.Reversal) (string, error) {
	body, err := json.Marshal(reversal)
	if err != nil {
		return "", err
	}

	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "safaricom/reversal/v1/request" //TODO :: CONFIRM THIS URL/ENDPOINT???
	return s.newReq(url, body, headers)
}

// BalanceInquiry sends a balance inquiry
func (s Service) BalanceInquiry(balanceInquiry models.BalanceInquiry) (string, error) {
	auth, err := s.auth()
	if err != nil {
		return "", nil
	}

	body, err := json.Marshal(balanceInquiry)
	if err != nil {
		return "", err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"
	headers["postman-token"] = "2aa448be-7d56-a796-065f-b378ede8b136"

	url := s.baseURL() + "mpesa/accountbalance/v1/query"
	return s.newReq(url, body, headers)
}

func (s Service) newReq(url string, body []byte, headers map[string]string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", nil
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}

	stringBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(stringBody), nil
}

func (s Service) newReq2(url string, body []byte, headers map[string]string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	stringBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return stringBody, nil
}

func (s Service) baseURL() string {
	if s.Env == PRODUCTION {
		return "https://api.safaricom.co.ke/"
	}
	return "https://sandbox.safaricom.co.ke/"
}
