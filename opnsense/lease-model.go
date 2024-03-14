package opnsense

import (
	"net/netip"
)

type DHCPLease struct {
	Address              netip.Addr `json:"address"`
	Starts               string     `json:"starts"`
	Ends                 string     `json:"ends"`
	CLTT                 int        `json:"cltt"`
	Binding              string     `json:"binding"`
	ClientHostname       string     `json:"client-hostname"`
	Type                 string     `json:"type"`
	Status               string     `json:"status"`
	Description          string     `json:"descr"`
	MAC                  string     `json:"mac"`
	Hostname             string     `json:"hostname"`
	State                string     `json:"state"`
	Man                  string     `json:"man"`
	Interface            string     `json:"if"`
	InterfaceDescription string     `json:"if_descr"`
}
