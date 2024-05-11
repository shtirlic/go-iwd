package iwd

import (
	"github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdService = "net.connman.iwd"
	iwdObjPath = "/net/connman/iwd"
)

type Iwd struct {
	conn *dbus.Conn
}

func NewIwd() (*Iwd, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return NewIwdWithConn(conn), nil
}

func NewIwdWithConn(conn *dbus.Conn) *Iwd {
	return &Iwd{conn: conn}
}

func (i *Iwd) Close() error {
	return i.conn.Close()
}

func (i *Iwd) Stations() ([]*Station, error) {
	return utils.GetObjectsByInterface(i.conn, iwdService, iwdStationIface, NewStation, i)
}

func (i *Iwd) Networks() ([]*Network, error) {
	return utils.GetObjectsByInterface(i.conn, iwdService, iwdNetworkIface, NewNetwork, i)
}

func (i *Iwd) Adapters() ([]*Adapter, error) {
	return utils.GetObjectsByInterface(i.conn, iwdService, iwdAdapterIface, NewAdapter, i)
}

func (i *Iwd) Devices() ([]*Device, error) {
	return utils.GetObjectsByInterface(i.conn, iwdService, iwdDeviceIface, NewDevice, i)
}

func (i *Iwd) KnownNetworks() ([]*KnownNetwork, error) {
	return utils.GetObjectsByInterface(i.conn, iwdService, iwdKnownNetworkIface, NewKnownNetwork, i)
}

func (i *Iwd) daemons() ([]*Daemon, error) {
	return utils.GetObjectsByInterface(i.conn, iwdService, iwdDaemonIface, NewDaemon, i)
}

func (i *Iwd) Daemon() (*Daemon, error) {
	if d, err := i.daemons(); err == nil {
		return d[0], nil
	} else {
		return nil, err
	}
}

func (i *Iwd) WSC() ([]*WSC, error) {
	stations, err := i.Stations()
	var wscs []*WSC
	if err != nil {
		return nil, err
	}
	for _, sta := range stations {
		if wsc, err := NewWSC(sta.Path, i); err == nil {
			wscs = append(wscs, wsc)
		} else {
			return nil, err
		}
	}
	return wscs, nil
}

func (i *Iwd) CallServiceMethod(path dbus.ObjectPath, method string, args ...interface{}) (*dbus.Call, error) {
	return utils.CallMethod(i.conn, iwdService, path, method, args...)
}
