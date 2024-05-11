package iwd

import (
	"github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdDeviceIface = iwdService + ".Device"
)

type Device struct {
	Path    dbus.ObjectPath // /net/connman/iwd/{phy0,phy1,...}/{1,2,...}
	Adapter *Adapter        // [ro] Adapter the device belongs to
	Address string          // [ro] Interface's hardware address in the XX:XX:XX:XX:XX:XX format
	Mode    DeviceMode      // [rw] Use to set the device mode
	Name    string          // [ro] Device's interface name
	Powered bool            // [rw]
	iwd     *Iwd
}

func NewDevice(p dbus.ObjectPath, i *Iwd) (*Device, error) {
	objects, err := utils.GetAllProperties(i.conn, iwdService, p, iwdDeviceIface)
	if err != nil {
		return nil, err
	}
	var adapter *Adapter
	if adapterValue := objects["Adapter"].Value(); adapterValue != nil {
		if adapter, err = NewAdapter(adapterValue.(dbus.ObjectPath), i); err != nil {
			return nil, err
		}
	}
	return &Device{
		Path:    p,
		Adapter: adapter,
		Address: objects["Address"].Value().(string),
		Mode:    DeviceMode(objects["Mode"].Value().(string)),
		Name:    objects["Name"].Value().(string),
		Powered: objects["Powered"].Value().(bool),
		iwd:     i,
	}, nil
}
