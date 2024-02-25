package opnsense

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
