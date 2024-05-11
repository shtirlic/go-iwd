package iwd

import (
	"github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdDaemonIface    = iwdService + ".Daemon"
	callDaemonGetIfno = iwdDaemonIface + ".GetInfo"
)

type Daemon struct {
	Path dbus.ObjectPath // /net/connman/iwd
	iwd  *Iwd
}

type DaemonInfo struct {
	Version                     string // IWD release version
	StateDirectory              string //  Absolute path to the IWD state directory
	NetworkConfigurationEnabled bool   // Whether networkconfiguration is enabled (see iwd(8))
}

func NewDaemon(p dbus.ObjectPath, i *Iwd) (*Daemon, error) {
	return &Daemon{
		Path: p,
		iwd:  i,
	}, nil
}

// Returns basic IWD daemon's status and configuration
// properties.  Their values are global and may be useful
// for D-Bus clients interacting with IWD, not so much
// for the user.  The returned dictionary (a{sv}) maps
// string keys to values of types defined per key.
// Clients should ignore unknown keys.
func (d *Daemon) GetInfo() (*DaemonInfo, error) {
	call, err := d.iwd.CallServiceMethod(d.Path, callDaemonGetIfno)
	if err != nil {
		return nil, err
	}
	var objects utils.DBusMapVariant
	if err := call.Store(&objects); err != nil {
		return nil, err
	}
	var dinfo DaemonInfo
	if err := utils.Transcode(objects, &dinfo); err != nil {
		return nil, err
	}
	return &dinfo, nil
}
