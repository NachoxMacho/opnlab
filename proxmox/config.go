package proxmox

import (
	"errors"
	"os"
)

var proxmoxURL string
var proxmoxToken string

func InitalizeConfig() error {

	proxmoxToken = os.Getenv("PVE_TOKEN")
	if proxmoxToken == "" {
		return errors.New("missing env variable `PVE_TOKEN`")
	}

	proxmoxURL = os.Getenv("PVE_BASE_URL")
	if proxmoxURL == "" {
		return errors.New("missing env variable `PVE_BASE_URL`")
	}

	return nil
}
