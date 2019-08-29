package virtcontainers

import "fmt"

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
	return nil, fmt.Errorf("vlan net endpoint is not implemented")
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
func (endpoint *VlanEndpoint) Attach(hypervisor) error {
	return fmt.Errorf("attach for vlan endpoint is not yet implemented")
}

// Detach for the vlan endpoint tears down the tap and bridge
// created for the vlan interface.
func (endpoint *VlanEndpoint) Detach(netNsCreated bool, netNsPath string) error {
	return fmt.Errorf("detach for vlan endpoint is not yet implemented")
}

// HotAttach for the vlan endpoint uses hot plug device
func (endpoint *VlanEndpoint) HotAttach(h hypervisor) error {
	return fmt.Errorf("hot attach for vlan endpoint is not yet implemented")
}

// HotDetach for the vlan endpoint uses hot pull device
func (endpoint *VlanEndpoint) HotDetach(h hypervisor, netNsCreated bool, netNsPath string) error {
	return fmt.Errorf("hot detach for vlan endpoint is not yet implemented")
}
