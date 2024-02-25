package opnsense

func GetDHCPLeases() ([]DHCPLease, error) {
	return getOPNSenseData[DHCPLease]("/dhcpv4/leases/searchLease")
}
