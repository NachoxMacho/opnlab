package opnsense

import (
	"errors"
	"os"
)

var opnsenseURL string
var opnsenseKey string
var opnsenseSecret string

func InitalizeConfig() error {

	opnsenseKey = os.Getenv("OPNSENSE_API_KEY")
	if opnsenseKey == "" {
		return errors.New("missing env variable `OPNSENSE_API_KEY`")
	}
	opnsenseSecret = os.Getenv("OPNSENSE_API_SECRET")
	if opnsenseKey == "" {
		return errors.New("missing env variable `OPNSENSE_API_SECRET`")
	}

	opnsenseURL = os.Getenv("OPNSENSE_BASE_URL")
	if opnsenseURL == "" {
		return errors.New("missing env variable `OPNSENSE_BASE_URL`")
	}

	return nil
}
