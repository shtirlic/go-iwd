package iwd

import (
	"github.com/godbus/dbus/v5"
	"github.com/shtirlic/go-iwd/utils"
)

const (
	iwdStationIface           = iwdService + ".Station"
	iwdStationDiagnosticIface = iwdService + ".StationDiagnostic"

	callStationScan                     = iwdStationIface + ".Scan"
	callStationDisconnect               = iwdStationIface + ".Disconnect"
	callStationGetOrderedNetworks       = iwdStationIface + ".GetOrderedNetworks"
	callStationConnectHiddenNetwork     = iwdStationIface + ".ConnectHiddenNetwork"
	callStationDiagnosticGetDiagnostics = iwdStationDiagnosticIface + ".GetDiagnostics"
)

// Reflects the general network connection state.  One of:
// "connected", "disconnected", "connecting",
// "disconnecting", "roaming"
type ConnectionState string

const (
	ConnectedState     ConnectionState = "connected"
	DisconnectedState  ConnectionState = "disconnected"
	ConnectingState    ConnectionState = "connecting"
	DisconnectingState ConnectionState = "discconnecting"
	RoamingState       ConnectionState = "roaming"
)

type NetworkWithSignal struct {
	*Network
	SignalStrength SignalStrength
}

// Network's maximum signal strength expressed
// in 100 * dBm.  The value is the range of 0
// (strongest signal) to -10000 (weakest signal)
type SignalStrength int16

type Station struct {
	Path             dbus.ObjectPath // /net/connman/iwd/{phy0,phy1,...}/{1,2,...}
	ConnectedNetwork *Network        // [ro] Reflects the object representing the network the device is currently connected to or to which a connection is in progress.
	Scanning         bool            // [ro] Reflects whether the station is currently scanning for networks.
	State            ConnectionState // [ro] Reflects the general network connection state.
	iwd              *Iwd
}

type StationDiagnosticInfo struct {
	AverageRSSI    int    // Average RSSI of currently connected BSS.
	Channel        int    // The WLAN channel number of currently connected BSS.
	ConnectedBss   string // MAC address of currently connected BSS.
	Frequency      int    // Frequency of currently connected BSS.
	PairwiseCipher string // The pairwise cipher chosen for this connection.
	RSSI           int    // The RSSI of the currently connected BSS.
	RxRate         int    // Receive rate in 100kbit/s
	RxBitrate      int    // Receive rate in 100kbit/s
	RxMCS          int    // Receiving MCS index
	RxMode         string // The phy technology being used (802.11n, 802.11ac or 802.11ax).
	Security       string // The chosen security for the connection.
	TxRate         int    // Transmission rate in 100kbit/s
	TxBitrate      int    // Transmission rate in 100kbit/s
	TxMCS          int    // Transmitting MCS index
	TxMode         string // Same meaning as RxMode, just fortransmission.
}

func NewStation(p dbus.ObjectPath, i *Iwd) (*Station, error) {
	objects, err := utils.GetAllProperties(i.conn, iwdService, p, iwdStationIface)
	if err != nil {
		return nil, err
	}
	var cnetwork *Network
	if cnetworkValue := objects["ConnectedNetwork"].Value(); cnetworkValue != nil {
		if cnetwork, err = NewNetwork(cnetworkValue.(dbus.ObjectPath), i); err != nil {
			return nil, err
		}
	}
	return &Station{
		Path:             p,
		ConnectedNetwork: cnetwork,
		Scanning:         objects["Scanning"].Value().(bool),
		State:            ConnectionState(objects["State"].Value().(string)),
		iwd:              i,
	}, nil
}

// Schedule a network scan.
func (s *Station) Scan() error {
	if _, err := s.iwd.CallServiceMethod(s.Path, callStationScan); err != nil {
		return err
	}
	return nil
}

// Disconnect from the network. This also disables
// iwd from trying to autoconnect to any other network
// with this device.
func (s *Station) Disconnect() error {
	if _, err := s.iwd.CallServiceMethod(s.Path, callStationDisconnect); err != nil {
		return err
	}
	return nil
}

// Return the list of networks found in the most recent
// scan, sorted by their user interface importance
// score as calculated by iwd.  If the device is
// currently connected to a network, that network is
// always first on the list, followed by any known
// networks that have been used at least once before,
// followed by any other known networks and any other
// detected networks as the last group.  Within these
// groups the maximum relative signal-strength is the
// main sorting factor.
func (s *Station) GetOrderedNetworks() ([]NetworkWithSignal, error) {
	call, err := s.iwd.CallServiceMethod(s.Path, callStationGetOrderedNetworks)
	if err != nil {
		return nil, err
	}
	var objects utils.DBusArrTupleVariant
	var oNets []NetworkWithSignal
	if err = call.Store(&objects); err != nil {
		return nil, err
	}
	for _, i := range objects {
		p := i[0].Value().(dbus.ObjectPath)
		ss := SignalStrength(i[1].Value().(int16))
		if network, err := NewNetwork(p, s.iwd); err == nil {
			oNets = append(oNets, NetworkWithSignal{network, ss})
		} else {
			return nil, err
		}
	}
	return oNets, nil
}

// Tries to find and connect to a hidden network for the
// first time.  Only hidden networks of type 'psk' and
// 'open' are supported.  WPA-Enterprise hidden networks
// must be provisioned.

// The ssid parameter is used to find the hidden network.
// If no network with the given ssid is found, an
// net.connman.iwd.NotFound error is returned.

// In the unlikely case that both an open and pre-shared
// key hidden network with the given ssid is found an
// net.connman.iwd.ServiceSetOverlap error is returned.

// Once the hidden network is found, the connection will
// proceed as normal.  So the user agent will be asked
// for a passphrase, etc.

// This method should only be called once to provision
// a hidden network.  For all future connections the
// regular Network.Connect() API should be used.
func (s *Station) ConnectHiddenNetwork(ssid string) error {
	if _, err := s.iwd.CallServiceMethod(s.Path, callStationConnectHiddenNetwork,
		ssid); err != nil {
		return err
	}
	return nil
}

// Get all diagnostic information for this interface. The
// diagnostics are contained in a single dictionary. Values
// here are generally low level and not meant for general
// purpose applications which could get by with the
// existing Station interface or values which are volatile
// and change too frequently to be represented as
// properties. The values in the dictionary may come and
// go depending on the state of IWD.
func (s *Station) GetDiagnostics() (*StationDiagnosticInfo, error) {
	call, err := s.iwd.CallServiceMethod(s.Path, callStationDiagnosticGetDiagnostics)
	if err != nil {
		return nil, err
	}
	var objects utils.DBusMapVariant
	if err := call.Store(&objects); err != nil {
		return nil, err
	}
	var diaginfo StationDiagnosticInfo
	if err := utils.Transcode(objects, &diaginfo); err != nil {
		return nil, err
	}
	return &diaginfo, nil
}
