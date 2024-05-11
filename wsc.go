package iwd

import (
	"github.com/godbus/dbus/v5"
)

const (
	iwdWSCIface = iwdService + ".SimpleConfiguration"

	callWSCPushButton  = iwdWSCIface + ".PushButton"
	callWSCGeneratePin = iwdWSCIface + ".GeneratePin"
	callWSCStartPin    = iwdWSCIface + ".StartPin"
	callWSCCancel      = iwdStationIface + ".Cancel"
)

type WSC struct {
	Path dbus.ObjectPath // /net/connman/iwd/{phy0,phy1,...}/{1,2,...}
	iwd  *Iwd
}

func NewWSC(p dbus.ObjectPath, i *Iwd) (*WSC, error) {
	return &WSC{
		Path: p,
		iwd:  i,
	}, nil
}

// Start WSC (formerly known as WPS, Wi-Fi Protected
// Setup) configuration in PushButton mode or trigger a
// connection to a specific P2P peer.  The usage will
// depend on which object this interface is found on.

// In the first use case any connected networks on the
// device will be disconnected and scanning will commence
// to find the access point in PushButton mode.  If
// multiple access points are found, then a
// SessionOverlap error will be returned.

// This method returns once the configuration has been
// completed and the network or the P2P peer has been
// successfully connected.
func (w *WSC) PushButton() error {
	if _, err := w.iwd.CallServiceMethod(w.Path, callWSCPushButton); err != nil {
		return err
	}
	return nil
}

// Generates a random 8 digit PIN with an included check
// digit suitable for use by most user interfaces.
func (w *WSC) GeneratePin() (string, error) {
	call, err := w.iwd.CallServiceMethod(w.Path, callWSCGeneratePin)
	if err != nil {
		return "", err
	}
	var obj dbus.Variant
	if err := call.Store(&obj); err != nil {
		return "", err
	}
	return obj.Value().(string), nil
}

// Start WSC or connect to a specific P2P peer in PIN
// mode.  If iwd's WSC configuration indicates that the
// device does not support a display, a static PIN from
// the main.conf configuration file is used.  Contents
// of pin are ignored in this case.

// Otherwise, the pin provided will be utilized.  This
// can be an automatically generated PIN that contains a
// check digit, or a user-specified PIN.  The
// GeneratePin() method can be used a generate a random
// 8 digit PIN with an included check digit.

// This method returns once the configuration has been
// completed and the network or the P2P peer has been
// successfully connected.
func (w *WSC) StartPin(pin string) error {
	if _, err := w.iwd.CallServiceMethod(w.Path, callWSCStartPin, pin); err != nil {
		return err
	}
	return nil
}

// Aborts any ongoing WSC operations or a P2P connection.
// If no operation is ongoing, net.connman.iwd.NotAvailable
// is returned.
func (w *WSC) Cancel() error {
	if _, err := w.iwd.CallServiceMethod(w.Path, callWSCCancel); err != nil {
		return err
	}
	return nil
}
