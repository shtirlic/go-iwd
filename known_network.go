package iwd

import (
	"github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdKnownNetworkIface   = iwdService + ".KnownNetwork"
	callKnownNetworkForget = iwdKnownNetworkIface + ".Forget"
)

type KnownNetwork struct {
	Path              dbus.ObjectPath
	AutoConnect       bool   // [rw]
	Hidden            bool   // [ro]
	LastConnectedTime string // [ro]
	Name              string // [ro]
	Type              string // [ro]
	iwd               *Iwd
}

func NewKnownNetwork(p dbus.ObjectPath, i *Iwd) (*KnownNetwork, error) {
	objects, err := utils.GetAllProperties(i.conn, iwdService, p, iwdKnownNetworkIface)
	if err != nil {
		return nil, err
	}
	return &KnownNetwork{
		Path:              p,
		AutoConnect:       objects["AutoConnect"].Value().(bool),
		Hidden:            objects["Hidden"].Value().(bool),
		LastConnectedTime: objects["LastConnectedTime"].Value().(string),
		Name:              objects["Name"].Value().(string),
		Type:              objects["Type"].Value().(string),
		iwd:               i,
	}, nil
}

// Removes the network from the 'known networks' list and
// removes any associated configuration data.  If the
// network is currently connected, then it is immediately
// disconnected.
func (k *KnownNetwork) Forget() error {
	if _, err := k.iwd.CallServiceMethod(k.Path, callKnownNetworkForget); err != nil {
		return err
	}
	return nil
}
