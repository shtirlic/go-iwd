package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	callGetManagedObjects = "org.freedesktop.DBus.ObjectManager.GetManagedObjects"
	callPropertiesGetAll  = "org.freedesktop.DBus.Properties.GetAll"
)

type ObjectConstructor[T any, V any] func(dbus.ObjectPath, V) (*T, error)

type DBusRetValues map[dbus.ObjectPath]map[string]map[string]dbus.Variant
type DBusMapVariant map[string]dbus.Variant
type DBusArrTupleVariant [][]dbus.Variant

// Transcode converts a map of dbus.Variant values to JSON and decodes it into an output interface.
func Transcode(in map[string]dbus.Variant, out interface{}) error {
	tmp := make(map[string]interface{})
	for k, v := range in {
		tmp[k] = v.Value()
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(tmp); err != nil {
		return fmt.Errorf("error encoding map to JSON: %w", err)
	}
	if err := json.NewDecoder(buf).Decode(out); err != nil {
		return fmt.Errorf("error decoding JSON to output interface: %w", err)
	}
	return nil
}

func GetObjectsByInterface[T any, V any](conn *dbus.Conn, service string, interfaceName string,
	constructor ObjectConstructor[T, V], i V) ([]*T, error) {

	objects, err := GetManagedObjects(conn, service)
	if err != nil {
		return nil, err
	}
	var result []*T
	for p, v := range objects {
		for r := range v {
			if r == interfaceName {
				obj, err := constructor(p, i)
				if err != nil {
					return nil, err
				}
				result = append(result, obj)
			}
		}
	}
	return result, nil
}

func CallMethod(conn *dbus.Conn, service string, path dbus.ObjectPath, method string, args ...interface{}) (*dbus.Call, error) {
	obj := conn.Object(service, path)
	if call := obj.Call(method, 0, args...); call.Err == nil {
		return call, nil
	} else {
		return nil, call.Err
	}
}

func GetManagedObjects(conn *dbus.Conn, service string) (DBusRetValues, error) {
	call, err := CallMethod(conn, service, "/", callGetManagedObjects)
	if err != nil {
		return nil, err
	}
	var objects DBusRetValues
	if err = call.Store(&objects); err != nil {
		return nil, err
	}
	return objects, nil
}

func GetAllProperties(conn *dbus.Conn, service string, path dbus.ObjectPath, iface string) (DBusMapVariant, error) {
	call, err := CallMethod(conn, service, path, callPropertiesGetAll, iface)
	if err != nil {
		return nil, err
	}
	var objects DBusMapVariant
	if err = call.Store(&objects); err != nil {
		return nil, err
	}
	return objects, nil
}
