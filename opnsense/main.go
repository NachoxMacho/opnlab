package opnsense

import (
	"net/netip"
	"slices"
)

func GetDHCPLeases() ([]DHCPLease, error) {
	return getOPNSenseData[DHCPLease]("/dhcpv4/leases/searchLease")
}

func GetInterfaces() ([]Interface, error) {
	return getOPNSenseData[Interface]("/interfaces/overview/interfacesInfo")
}

func GetNextAvailableIP() (netip.Addr, error) {
	leases, err := GetDHCPLeases()
	if err != nil {
		return netip.Addr{}, err
	}

	startIP := netip.AddrFrom4([4]byte{10, 10, 10, 25})

	for ip := startIP; ip.Compare(netip.Addr{}) != 0; ip = ip.Next() {
		index := slices.IndexFunc(leases, func(d DHCPLease) bool {
			return ip.Compare(d.Address) == 0
		})
		if index == -1 {
			return ip, nil
		}
	}

	return netip.Addr{}, nil
}
