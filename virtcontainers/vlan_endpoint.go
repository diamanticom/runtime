package virtcontainers

import (
	"fmt"
	"github.com/containernetworking/plugins/pkg/ns"
	persistapi "github.com/kata-containers/runtime/virtcontainers/persist/api"
)

// VlanEndpoint gathers a network pair and its properties.
type VlanEndpoint struct {
	NetPair            NetworkInterfacePair
	EndpointProperties NetworkInfo
	EndpointType       EndpointType
	// todo: vlan id, how does this differ from pciaddress? how is it used?
	VlanId  int
	PCIAddr string
}

func createVlanNetworkEndpoint(idx int, ifName string, interworkingModel NetInterworkingModel) (*VlanEndpoint, error) {
	if idx < 0 {
		return &VlanEndpoint{}, fmt.Errorf("invalid network endpoint index: %d", idx)
	}
	netPair, err := createNetworkInterfacePair(idx, ifName, interworkingModel)
	if err != nil {
		return nil, err
	}

	endpoint := &VlanEndpoint{
		// TODO This is too specific. We may need to create multiple
		// end point types here and then decide how to connect them
		// at the time of hypervisor attach and not here
		NetPair:      netPair,
		EndpointType: VlanEndpointType,
	}
	if ifName != "" {
		endpoint.NetPair.VirtIface.Name = ifName
	}

	return endpoint, nil

}

// Properties returns properties for the vlan interface in the network pair.
func (endpoint *VlanEndpoint) Properties() NetworkInfo {
	return endpoint.EndpointProperties
}

// Name returns name of the vlan interface in the network pair.
func (endpoint *VlanEndpoint) Name() string {
	return endpoint.NetPair.VirtIface.Name
}

// HardwareAddr returns the mac address that is assigned to the tap interface
// in th network pair.
func (endpoint *VlanEndpoint) HardwareAddr() string {
	return endpoint.NetPair.TAPIface.HardAddr
}

// Type identifies the endpoint as a vlan endpoint.
func (endpoint *VlanEndpoint) Type() EndpointType {
	return endpoint.EndpointType
}

// PciAddr returns the PCI address of the endpoint.
func (endpoint *VlanEndpoint) PciAddr() string {
	return endpoint.PCIAddr
}

// SetPciAddr sets the PCI address of the endpoint.
func (endpoint *VlanEndpoint) SetPciAddr(pciAddr string) {
	endpoint.PCIAddr = pciAddr
}

// NetworkPair returns the network pair of the endpoint.
func (endpoint *VlanEndpoint) NetworkPair() *NetworkInterfacePair {
	return &endpoint.NetPair
}

// SetProperties sets the properties for the endpoint.
func (endpoint *VlanEndpoint) SetProperties(properties NetworkInfo) {
	endpoint.EndpointProperties = properties
}

// Attach for vlan endpoint bridges the network pair and adds the
// tap interface of the network pair to the hypervisor.
func (endpoint *VlanEndpoint) Attach(h hypervisor) error {
	if err := xConnectVMNetwork(endpoint, h); err != nil {
		networkLogger().WithError(err).Error("Error bridging vlan ep")
		return err
	}

	rc := h.addDevice(endpoint, netDev)
	return rc
}

// Detach for the vlan endpoint tears down the tap and bridge
// created for the vlan interface.
func (endpoint *VlanEndpoint) Detach(netNsCreated bool, netNsPath string) error {
	// The network namespace would have been deleted at this point
	// if it has not been created by virtcontainers.
	if !netNsCreated {
		return nil
	}

	return doNetNS(netNsPath, func(_ ns.NetNS) error {
		return xDisconnectVMNetwork(endpoint)
	})
}

// HotAttach for the vlan endpoint uses hot plug device
func (endpoint *VlanEndpoint) HotAttach(h hypervisor) error {
	return fmt.Errorf("hot attach for vlan endpoint is not yet implemented")
}

// HotDetach for the vlan endpoint uses hot pull device
func (endpoint *VlanEndpoint) HotDetach(h hypervisor, netNsCreated bool, netNsPath string) error {
	return fmt.Errorf("hot detach for vlan endpoint is not yet implemented")
}

func (endpoint *VlanEndpoint) save() persistapi.NetworkEndpoint {
	netpair := saveNetIfPair(&endpoint.NetPair)

	return persistapi.NetworkEndpoint{
		Type: string(endpoint.Type()),
		Vlan: &persistapi.VlanEndpoint{
			NetPair: *netpair,
		},
	}
}

func (endpoint *VlanEndpoint) load(s persistapi.NetworkEndpoint) {
	endpoint.EndpointType = VlanEndpointType

	if s.Vlan != nil {
		netpair := loadNetIfPair(&s.Vlan.NetPair)
		endpoint.NetPair = *netpair
	}
}
