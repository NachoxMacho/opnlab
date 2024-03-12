package opnsense

import (
	"encoding/json"
	"net"
	"net/netip"
)

type DHCPLease struct {
	Address              netip.Addr       `json:"address"`
	Starts               string           `json:"starts"`
	Ends                 string           `json:"ends"`
	CLTT                 int              `json:"cltt"`
	Binding              string           `json:"binding"`
	ClientHostname       string           `json:"client-hostname"`
	Type                 string           `json:"type"`
	Status               string           `json:"status"`
	Description          string           `json:"descr"`
	MAC                  net.HardwareAddr `json:"mac"`
	Hostname             string           `json:"hostname"`
	State                string           `json:"state"`
	Man                  string           `json:"man"`
	Interface            string           `json:"if"`
	InterfaceDescription string           `json:"if_descr"`
}

func (l *DHCPLease) UnmarshalJSON(b []byte) error {

	type lease struct {
		Address        string `json:"address"`
		Starts         string `json:"starts"`
		Ends           string `json:"ends"`
		CLTT           int    `json:"cltt"`
		Binding        string `json:"binding"`
		ClientHostname string `json:"client-hostname"`
		Type           string `json:"type"`
		Status         string `json:"status"`
		Description    string `json:"descr"`
		MAC            string `json:"mac"`
		Hostname       string `json:"hostname"`
		State          string `json:"state"`
		Man            string `json:"man"`
		If             string `json:"if"`
		IfDescription  string `json:"if_descr"`
	}

	tmpLease := lease{}
	err := json.Unmarshal(b, &tmpLease)
	if err != nil {
		return err
	}

	MAC, err := net.ParseMAC(tmpLease.MAC)
	if err != nil {
		return err
	}

	Address, err := netip.ParseAddr(tmpLease.Address)
	if err != nil {
		return err
	}

	l.MAC = MAC
	l.Address = Address
	l.Starts = tmpLease.Starts
	l.Ends = tmpLease.Ends
	l.CLTT = tmpLease.CLTT
	l.Binding = tmpLease.Binding
	l.ClientHostname = tmpLease.ClientHostname
	l.Type = tmpLease.Type
	l.Status = tmpLease.Status
	l.Description = tmpLease.Description
	l.Hostname = tmpLease.Hostname
	l.State = tmpLease.State
	l.Man = tmpLease.Man
	l.Interface = tmpLease.If
	l.InterfaceDescription = tmpLease.IfDescription

	return nil
}
