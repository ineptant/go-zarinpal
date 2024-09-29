package zarinpal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Zarinpal struct {
	MerchantID      string
	Sandbox         bool
	APIEndpoint     string
	PaymentEndpoint string
}

type paymentRequestReqBody struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	CallbackURL string `json:"callback_url"`
}

type paymentRequestRespData struct {
	Authority string `json:"authority"`
	Fee       int    `json:"fee"`
	FeeType   string `json:"fee_type"`
	Status    int    `json:"code"`
	Message   string `json:"message"`
}

type paymentRequestResp struct {
	Data   paymentRequestRespData
	Errors []any
}

type paymentVerificationReqBody struct {
	MerchantID string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Authority  string `json:"authority"`
}

type paymentVerificationRespData struct {
	Status    int         `json:"code"`
	Message   string      `json:"message"`
	CardHash  string      `json:"card_hash"`
	CardPan   string      `json:"card_pan"`
	RefID     json.Number `json:"ref_id"`
	FeeType   string      `json:"fee_type"`
	Fee       int         `json:"fee"`
}

type paymentVerificationResp struct {
	Data   paymentVerificationRespData
	Errors []any
}

type unverifiedTransactionsReqBody struct {
	MerchantID string `json:"merchant_id"`
}

type UnverifiedAuthority struct {
	Authority   string `json:"authority"`
	Amount      int    `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Referer     string `json:"referer"`
	Date        string `json:"date"`
}

type unverifiedTransactionsRespData struct {
	Status      int                   `json:"code"`
	Message     string                `json:"message"`
	Authorities []UnverifiedAuthority `json:"authorities"`
}

type unverifiedTransactionsResp struct {
	Data unverifiedTransactionsRespData
}

type transactionInquiryReqBody struct {
	MerchantID string `json:"merchant_id"`
	Authority  string `json:"authority"`
}

type transactionInquiryRespData struct {
	Code    string `json:"status"`
	Status  int    `json:"code"`
	Message string `json:"message"`
}

type transactionInquiryResp struct {
	Data    transactionInquiryRespData
	Message string
	Errors  []any
}

type reverseTransactionRespData struct {
	Status  int    `json:"code"`
	Message string `json:"message"`
}

type reverseTransactionResp struct {
	Data    reverseTransactionRespData
	Errors  []any
}

func New(merchantID string, sandbox bool) (*Zarinpal, error) {
	if len(merchantID) != 36 {
		return nil, errors.New("MerchantID must be 36 characters")
	}

	apiEndpoint := "https://api.zarinpal.com/pg/v4/payment/"
	paymentEndpoint := "https://payment.zarinpal.com/pg/StartPay/"
	if sandbox == true {
		apiEndpoint = "https://sandbox.zarinpal.com/pg/v4/payment/"
		paymentEndpoint = "https://sandbox.zarinpal.com/pg/StartPay/"
	}

	return &Zarinpal{
		MerchantID:      merchantID,
		Sandbox:         sandbox,
		APIEndpoint:     apiEndpoint,
		PaymentEndpoint: paymentEndpoint,
	}, nil
}

func (zarinpal *Zarinpal) NewPayment(amount int, callbackURL, description string) (paymentURL, authority string, statusCode int, err error) {
	if amount < 1 {
		err = errors.New("amount must be a positive number")
		return
	}
	if callbackURL == "" {
		err = errors.New("callbackURL should not be empty")
		return
	}
	if description == "" {
		err = errors.New("description should not be empty")
		return
	}

	paymentRequest := paymentRequestReqBody{
		MerchantID:  zarinpal.MerchantID,
		Amount:      amount,
		Description: description,
		CallbackURL: callbackURL,
	}
	var resp paymentRequestResp
	err = zarinpal.request(zarinpal.APIEndpoint+"request.json", &paymentRequest, &resp)
	if err != nil {
		return
	}

	statusCode = resp.Data.Status
	if resp.Data.Status == 100 {
		authority = resp.Data.Authority
		paymentURL = zarinpal.PaymentEndpoint + resp.Data.Authority
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Status))
	}

	return
}

func (zarinpal *Zarinpal) VerifyTransaction(amount int, authority string) (verified bool, refID string, statusCode int, err error) {
	if amount <= 0 {
		err = errors.New("amount must be a positive number")
		return
	}
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}

	paymentVerification := paymentVerificationReqBody{
		MerchantID: zarinpal.MerchantID,
		Amount:     amount,
		Authority:  authority,
	}
	var resp paymentVerificationResp
	err = zarinpal.request(zarinpal.APIEndpoint+"verify.json", &paymentVerification, &resp)
	if err != nil {
		return
	}

	statusCode = resp.Data.Status
	if resp.Data.Status == 100 {
		verified = true
		refID = string(resp.Data.RefID)
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Status))
	}

	return
}

func (zarinpal *Zarinpal) UnverifiedTransactions() (authorities []UnverifiedAuthority, statusCode int, err error) {
	unverifiedTransactions := unverifiedTransactionsReqBody{
		MerchantID: zarinpal.MerchantID,
	}
	var resp unverifiedTransactionsResp
	err = zarinpal.request(zarinpal.APIEndpoint+"unVerified.json", &unverifiedTransactions, &resp)
	if err != nil {
		return
	}

	if resp.Data.Status == 100 {
		statusCode = resp.Data.Status
		authorities = resp.Data.Authorities
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Status))
	}

	return
}

func (zarinpal *Zarinpal) TransactionInquiry(authority string) (inquiry string, statusCode int, err error) {
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}

	transactionInquiry := transactionInquiryReqBody{
		MerchantID: zarinpal.MerchantID,
		Authority:  authority,
	}
	var resp transactionInquiryResp
	err = zarinpal.request(zarinpal.APIEndpoint+"inquiry.json", &transactionInquiry, &resp)
	if err != nil {
		return
	}

	if resp.Data.Status == 100 {
		statusCode = resp.Data.Status
		inquiry = resp.Data.Code
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Status))
	}

	return
}

func (zarinpal *Zarinpal) ReverseTransaction(authority string) (statusCode int, err error) {
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}

	transactionInquiry := transactionInquiryReqBody{
		MerchantID: zarinpal.MerchantID,
		Authority:  authority,
	}
	var resp reverseTransactionResp
	err = zarinpal.request(zarinpal.APIEndpoint+"reverse.json", &transactionInquiry, &resp)
	if err != nil {
		return
	}

	if resp.Data.Status == 100 {
		statusCode = resp.Data.Status
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Status))
	}

	return
}

func (zarinpal *Zarinpal) request(endpoint string, data interface{}, res interface{}) error {
	reqBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		err = errors.New("zarinpal invalid json response")
		return err
	}
	
	return nil
}
