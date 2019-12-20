package utils

// This code adapted from https://github.com/dpapathanasiou/go-recaptcha/blob/master/recaptcha.go
// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness
//
// Edit the recaptchaPrivateKey constant before building and using

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type Recaptcha struct {
	secret string
}

const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"

// check uses the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func check(remoteip, response string, secretKey string) (r RecaptchaResponse, err error) {
	resp, err := http.PostForm(recaptchaServerName,
		url.Values{"secret": {secretKey}, "remoteip": {remoteip}, "response": {response}})
	if err != nil {
		log.Printf("Post error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read error: could not read body: %s \n", err)
		return
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Read error: got invalid JSON: %s \n", err)
		return
	}
	return
}

// Confirm is the public interface function.
// It calls check, which the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func confirm(remoteip, response string, secretKey string) (result bool, err error) {
	resp, err := check(remoteip, response, secretKey)
	result = resp.Success
	return
}

// Init allows the webserver or code evaluating the reCaptcha form input to set the
// reCaptcha private key (string) value, which will be different for every domain.
func NewRecaptha(key string) *Recaptcha {
	return &Recaptcha{
		secret: key,
	}
}

// VerifyCaptch verifying recaptcha
func (recap *Recaptcha) VerifyCaptch(recaptchaResponse string, remoteIpAddress string) (bool, error) {

	result, err := confirm(remoteIpAddress, recaptchaResponse, recap.secret)
	if err != nil {
		log.Println("recaptcha server error", err)
	}
	return result, err

}
