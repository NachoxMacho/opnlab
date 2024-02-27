package opnsense

import (
	"net"
	"net/netip"
)

type DHCPLease struct {
	Address        string `json:"address"`
	Starts         string `json:"starts"`
	Ends           string `json:"ends"`
	CLTT           string `json:"cltt"`
	Binding        string `json:"binding"`
	ClientHostname string `json:"client-hostname"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	Description    string `json:"descr"`
	MACAddress     string `json:"mac"`
	Hostname       string `json:"hostname"`
	State          string `json:"state"`
	Man            string `json:"man"`
	If             string `json:"if"`
	IfDescription  string `json:"if_descr"`
}

func (l *DHCPLease) GetIP() (netip.Addr, error) {
	return netip.ParseAddr(l.Address)
}

func (l *DHCPLease) GetMACAddress() (net.HardwareAddr, error) {
	return net.ParseMAC(l.MACAddress)
}

type Interface struct {
	Flags          []string             `json:"flags"`
	Capabilities   []string             `json:"capabilities"`
	Options        []string             `json:"options"`
	MACAddress     string               `json:"macaddr"`
	SupportedMedia []string             `json:"supported_media"`
	IsPhysical     bool                 `json:"is_physical"`
	Device         string               `json:"device"`
	MTU            string               `json:"mtu"`
	Media          string               `json:"media"`
	MediaRaw       string               `json:"media_raw"`
	Status         string               `json:"status"`
	Routes         []string             `json:"routes"`
	Config         InterfaceConfig      `json:"config"`
	Identifier     string               `json:"identifier"`
	Description    string               `json:"description"`
	Enabled        bool                 `json:"enabled"`
	LinkType       string               `json:"link_type"`
	IPv4           []InterfaceIPAddress `json:"ipv4"`
	IPv6           []InterfaceIPAddress `json:"ipv6"`
	Gateways       []string             `json:"gateways"`
}

type InterfaceIPAddress struct {
	IPAddress string `json:"ipaddr"`
}

type InterfaceConfig struct {
	Enabled              string `json:"enable"`
	Interface            string `json:"if"`
	IPv4                 string `json:"ipaddr"`
	IPv6                 string `json:"ipaddr6"`
	Gateway              string `json:"gateway"`
	BlockPrivateNetworks string `json:"blockpriv"`
	BlockBogonNetworks   string `json:"blockbogons"`
	Media                string `json:"media"`
	MediaOption          string `json:"mediaopt"`
	DHCPv6               string `json:"dhcp6-ia-pd-len"`
	Identifier           string `json:"identifier"`
}
