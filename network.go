package iwd

import (
	dbus "github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdNetworkIface    = iwdService + ".Network"
	callNetworkConnect = iwdNetworkIface + ".Connect"
)

type NetworkType string

const (
	OpenNetworkType NetworkType = "open"
	WEPNetworkType  NetworkType = "wep"
	PSKNetworkType  NetworkType = "psk"
	EAPNetworkType  NetworkType = "8021x"
)

type NetworkInterface interface {
}

type Network struct {
	Path         dbus.ObjectPath // /net/connman/iwd/{phy0,phy1,...}/{1,2,...}/Xxx
	Name         string          // [ro] Network SSID
	Device       *Device         // [ro]
	Connected    bool            // [ro]
	KnownNetwork *KnownNetwork   // [ro] KnownNetwork object corresponding to this Network
	Type         NetworkType     // [ro] Contains the type of the network
	iwd          *Iwd
}

func NewNetwork(p dbus.ObjectPath, i *Iwd) (*Network, error) {
	objects, err := utils.GetAllProperties(i.conn, iwdService, p, iwdNetworkIface)
	if err != nil {
		return nil, err
	}
	var knetwork *KnownNetwork
	if knetworkValue := objects["KnownNetwork"].Value(); knetworkValue != nil {
		if knetwork, err = NewKnownNetwork(knetworkValue.(dbus.ObjectPath), i); err != nil {
			return nil, err
		}
	}
	var device *Device
	if deviceValue := objects["Device"].Value(); deviceValue != nil {
		if device, err = NewDevice(deviceValue.(dbus.ObjectPath), i); err != nil {
			return nil, err
		}
	}
	return &Network{
		Path:         p,
		Connected:    objects["Connected"].Value().(bool),
		Device:       device,
		KnownNetwork: knetwork,
		Name:         objects["Name"].Value().(string),
		Type:         NetworkType(objects["Type"].Value().(string)),
		iwd:          i,
	}, nil
}

func (n *Network) Connect() error {
	if _, err := n.iwd.CallServiceMethod(n.Path, callNetworkConnect); err != nil {
		return err
	}
	return nil
}
