package iwd

import (
	"github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdAdapterIface = iwdService + ".Adapter"
)

type DeviceMode string

const (
	AdHocDeviceMode   DeviceMode = "ad-hoc"
	StationDeviceMode DeviceMode = "station"
	APDeviceMode      DeviceMode = "ap"
)

type Adapter struct {
	Path           dbus.ObjectPath // /net/connman/iwd/{phy0,phy1,...}
	Name           string          // [ro] Contains the name of the adapter
	Model          string          // [ro] Contains the model name of the adapter, if available
	Vendor         string          // [ro] Contains the vendor name of the adapter, if available
	Powered        bool            // [rw]
	SupportedModes []DeviceMode    // [ro] Contains the supported modes for this adapter's devices
	iwd            *Iwd
}

func NewAdapter(p dbus.ObjectPath, i *Iwd) (*Adapter, error) {
	objects, err := utils.GetAllProperties(i.conn, iwdService, p, iwdAdapterIface)
	if err != nil {
		return nil, err
	}
	var modes []DeviceMode
	for _, m := range objects["SupportedModes"].Value().([]string) {
		modes = append(modes, DeviceMode(m))
	}
	return &Adapter{
		Path:           p,
		Model:          objects["Model"].Value().(string),
		Name:           objects["Name"].Value().(string),
		Powered:        objects["Powered"].Value().(bool),
		SupportedModes: modes,
		Vendor:         objects["Vendor"].Value().(string),
		iwd:            i,
	}, nil
}
