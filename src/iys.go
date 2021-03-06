package iys

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type ConsentSource string

type ConsentStatus string

type ConsentType string

type RecipientType string

type Config struct {
	BaseURL   string
	UserCode  string
	BrandCode string
	Username  string
	Password  string
}

type API struct {
	Config Config

	Client struct {
		Username  interface{} `json:"username,omitempty"`
		Password  interface{} `json:"password,omitempty"`
		GrantType interface{} `json:"grant_type,omitempty"`
	}

	Authentication struct {
		AccessToken      interface{} `json:"access_token,omitempty"`
		RefreshToken     interface{} `json:"refresh_token,omitempty"`
		TokenType        interface{} `json:"token_type,omitempty"`
		ExpiresIn        interface{} `json:"expires_in,omitempty"`
		RefreshExpiresIn interface{} `json:"refresh_expires_in,omitempty"`
	}
}

type Request struct {
	Recipient     interface{}   `json:"recipient,omitempty"`
	RecipientType RecipientType `json:"recipientType,omitempty"`
	ConsentType   ConsentType   `json:"type,omitempty"`
	ConsentSource ConsentSource `json:"source,omitempty"`
	ConsentStatus ConsentStatus `json:"status,omitempty"`
	ConsentDate   interface{}   `json:"consentDate,omitempty"`
}

type Response struct {
	TransactionID interface{} `json:"transactionId,omitempty"`
	CreationDate  interface{} `json:"creationDate,omitempty"`
}

const (
	Physical    ConsentSource = "HS_FIZIKSEL_ORTAM"
	Sign        ConsentSource = "HS_ISLAK_IMZA"
	Web         ConsentSource = "HS_WEB"
	Message     ConsentSource = "HS_MESAJ"
	CallCenter  ConsentSource = "HS_CAGRI_MERKEZI"
	Mobile      ConsentSource = "HS_MOBIL"
	Email       ConsentSource = "HS_EPOSTA"
	Emedia      ConsentSource = "HS_EORTAM"
	SocialMedia ConsentSource = "HS_SOSYAL_MEDYA"
	Event       ConsentSource = "HS_ETKINLIK"
	Year2015    ConsentSource = "HS_2015"
	Atm         ConsentSource = "HS_ATM"
	Decision    ConsentSource = "HS_KARAR"
)

const (
	Accept ConsentStatus = "ONAY"
	Reject ConsentStatus = "RED"
)

const (
	Call ConsentType = "ARAMA"
	Sms  ConsentType = "MESAJ"
	Mail ConsentType = "EPOSTA"
)

const (
	Individual RecipientType = "BIREYSEL"
	Merchant   RecipientType = "TACIR"
)

func (api *API) Authorize() bool {
	tokenurl := api.Config.BaseURL + "/oauth2/token"
	api.Client.Username = api.Config.Username
	api.Client.Password = api.Config.Password
	api.Client.GrantType = "password"
	post, _ := json.Marshal(api.Client)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", tokenurl, bytes.NewReader(post))
	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	if res.StatusCode != 200 {
		var description interface{}
		decoder.Decode(&description)
		pretty, _ := json.MarshalIndent(description, " ", "\t")
		log.Println(string(pretty))
		return false
	}
	decoder.Decode(&api.Authentication)
	return true
}

func (api *API) CreateConsent(request *Request) (response *Response) {
	response = new(Response)
	apiurl := api.Config.BaseURL + "/sps/" + api.Config.UserCode + "/brands/" + api.Config.BrandCode + "/consents"
	post, _ := json.Marshal(request)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(post))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.Authentication.AccessToken.(string))
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	if res.StatusCode != 200 {
		var description interface{}
		decoder.Decode(&description)
		pretty, _ := json.MarshalIndent(description, " ", "\t")
		log.Println(string(pretty))
		return response
	}
	decoder.Decode(&response)
	pretty, _ := json.MarshalIndent(response, " ", "\t")
	log.Println(string(pretty))
	return response
}
