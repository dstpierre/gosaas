package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// EmailProvider defines possible implementation of email providers.
type EmailProvider string

const (
	// EmailProviderSES used with Amazon SES email service.
	EmailProviderSES EmailProvider = "amazonses"
)

// Configuration defines important settings used across the library.
type Configuration struct {
	EmailFrom     string        `json:"emailFrom"`
	EmailFromName string        `json:"emailFromName"`
	EmailProvider EmailProvider `json:"emailProvider"`

	StripeKey string `json:"stripeKey"`

	SignUpTemplate            string `json:"signupTemplate"`
	SignUpSendEmailValidation bool   `json:"sendEmailValidation"`
	SignUpSuccessRedirect     string `json:"signupSuccessRedirect"`
	SignUpErrorRedirect       string `json:"signupErrorRedirect"`
	SignInTemplate            string `json:"signinTemplate"`
	SignInSuccessRedirect     string `json:"signinSuccessRedirect"`
	SignInErrorRedirect       string `json:"signinErrorRedirect"`
}

// Current holds the current configuration
var Current Configuration

func init() {
	b, err := ioutil.ReadFile("./gosaas.json")
	if err != nil {
		return
	}

	if err := json.Unmarshal(b, &Current); err != nil {
		log.Println("error parsing your gosaas.json config file", err)
	}
}

// Configure sets the proper values for various important aspects of the library
// controlling the behavior of important process like emails and signups.
//
// You may create a "gosaas.json" config file that will be automatically loaded
// at startup.
func Configure(conf Configuration) {
	Current = conf
}
