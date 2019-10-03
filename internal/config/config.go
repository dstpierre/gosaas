package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dstpierre/gosaas/data"
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
	Plans []data.BillingPlan `json:"plans"`

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

// LoadFromFile loads the ./gosaas.json file as the default library configuration
func LoadFromFile() error {
	b, err := ioutil.ReadFile("./gosaas.json")
	if err != nil { 
		return err
	}

	if err := json.Unmarshal(b, &Current); err != nil {
		return fmt.Errorf("error parsing your gosaas.json config file: %v", err)
	}
	return nil
}

// Configure sets the proper values for various important aspects of the library
// controlling the behavior of important process like emails and signups.
//
// You may create a "gosaas.json" config file that will be automatically loaded
// at startup.
func Configure(conf Configuration) {
	Current = conf
}
