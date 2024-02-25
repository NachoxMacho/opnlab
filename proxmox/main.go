package proxmox

import "fmt"

func GetVMs() ([]VM, error) {

	nodes, err := getProxmoxData[[]APINode]("/nodes")
	if err != nil {
		return nil, err
	}

	vms := []VM{}
	for _, node := range nodes {
		nodeVMs, err := getProxmoxData[[]APIVM]("/nodes/" + node.Name + "/qemu")
		if err != nil {
			return nil, err
		}

		for _, vm := range nodeVMs {
			v := VM{
				Stats: vm,
			}
			vmConfig, err := getProxmoxData[APIVMConfig]("/nodes/" + node.Name + "/qemu/" + fmt.Sprintf("%d", vm.VMID) + "/config")
			if err != nil {
				return nil, err
			}

			v.Config = vmConfig

			vms = append(vms, v)
		}
	}
	return vms, nil
}

func GetNodes() ([]APINode, error) {

	nodes, err := getProxmoxData[[]APINode]("/nodes")
	if err != nil {
		return nil, err
	}

	return nodes, err
}

func GetVMByID(id int) (VM, error) {

	nodes, err := getProxmoxData[[]APINode]("/nodes")
	if err != nil {
		return VM{}, err
	}

	for _, node := range nodes {
		nodeVMs, err := getProxmoxData[[]APIVM]("/nodes/" + node.Name + "/qemu")
		if err != nil {
			return VM{}, err
		}

		for _, vm := range nodeVMs {
			if vm.VMID != id {
				continue
			}
			v := VM{
				Stats: vm,
			}
			vmConfig, err := getProxmoxData[APIVMConfig]("/nodes/" + node.Name + "/qemu/" + fmt.Sprintf("%d", vm.VMID) + "/config")
			if err != nil {
				return VM{}, err
			}

			v.Config = vmConfig
			return v, nil
		}
	}
	return VM{}, nil
}
