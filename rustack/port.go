package rustack

import (
	"fmt"
)

type Port struct {
	manager           *Manager
	ID                string              `json:"id"`
	IpAddress         *string             `json:"ip_address,omitempty"`
	Network           *Network            `json:"network"`
	FirewallTemplates []*FirewallTemplate `json:"fw_templates,omitempty"`
	Connected         *Connected          `json:"connected"`
	Locked            bool                `json:"locked"`
}

type Connected struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Vdc  *Vdc   `json:"vdc"`
}

func NewPort(network *Network, firewallTemplates []*FirewallTemplate, ipAddress *string) Port {
	p := Port{Network: network, FirewallTemplates: firewallTemplates, IpAddress: ipAddress}
	return p
}

func (v *Vdc) GetPorts(extraArgs ...Arguments) (ports []*Port, err error) {
	args := Arguments{
		"vdc": v.ID,
	}

	args.merge(extraArgs)

	path := "v1/port"
	err = v.manager.GetItems(path, args, &ports)
	for i := range ports {
		ports[i].manager = v.manager
		ports[i].Network.manager = v.manager
	}
	return
}

func (p *Port) UpdateFirewall(firewallTemplates []*FirewallTemplate) error {
	path := fmt.Sprintf("v1/port/%s", p.ID)

	var fwTemplates = make([]*string, 0)
	for _, fwTemplate := range firewallTemplates {
		fwTemplates = append(fwTemplates, &fwTemplate.ID)
	}

	args := &struct {
		FwTemplates []*string `json:"fw_templates"`
	}{
		FwTemplates: fwTemplates,
	}

	err := p.manager.Put(path, args, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *Port) Delete() error {
	path := fmt.Sprintf("v1/port/%s", p.ID)
	return p.manager.Delete(path, Defaults(), p)
}

func (r *Router) CreatePort(port *Port, toConnect interface{}) (err error) {
	args := &struct {
		manager           *Manager
		ID                string              `json:"id"`
		IpAddress         *string             `json:"ip_address,omitempty"`
		Network           string              `json:"network"`
		Router            string              `json:"router,omitempty"`
		Vm                string              `json:"vm,omitempty"`
		Lbaas             string              `json:"lbaas,omitempty"`
		FirewallTemplates []*FirewallTemplate `json:"fw_templates,omitempty"`
	}{
		ID:                port.ID,
		IpAddress:         port.IpAddress,
		Network:           port.Network.ID,
		FirewallTemplates: port.FirewallTemplates,
	}
	switch v := toConnect.(type) {
	case *Router:
		args.Router = v.ID
	case *Vm:
		args.Vm = v.ID
		// TODO: Create lbaas
		// case Lbaas:
		// 	args.Lbaas = toConnect.(Lbaas).ID
	default:
		return fmt.Errorf("ERROR. Unknown type: %s", v)
	}
	err = r.manager.Post("v1/port", args, &port)
	return
}

func (m *Manager) GetPort(id string) (port *Port, err error) {
	path := fmt.Sprintf("v1/port/%s", id)
	err = m.Get(path, Defaults(), &port)
	if err != nil {
		return
	}
	port.manager = m
	return
}

func (p Port) WaitLock() (err error) {
	path := fmt.Sprintf("v1/port/%s", p.ID)
	return loopWaitLock(p.manager, path)
}
