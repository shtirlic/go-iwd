# go-iwd Go binding for iwd

__go-iwd__ is a Go D-BUS binding for iNet Wireless Daemon [iwd](https://iwd.wiki.kernel.org/) API.

> [!NOTE]
> Please keep in mind that project is under active development.

## Getting started

### Prerequisites

go-iwd requires [Go](https://go.dev/) version [1.22](https://go.dev/doc/devel/release#go1.22.0) or above.

### Getting go-iwd

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```sh
import "github.com/shtirlic/go-iwd"
```

Alternatively, use `go get`:

```sh
go get -u github.com/shtirlic/go-iwd
```

## Usage

```go
func main() {
  i, err := iwd.NewIwd()
  if err != nil {
    panic(err)
  }
  defer i.Close()

  // Get iwd Stations
  stations, err := i.Stations()
  if err != nil {
    panic(err)
  }

  // Disconnect from network
  for _, station := range stations {
    station.Disconnect()
  }

  // Connect to Known Network
  for _, station := range stations {
    if station.State == iwd.ConnectedState {
      return
    }
    nets, _ := station.GetOrderedNetworks()
    for _, net := range nets {
      net.Connect()
    }
  }

  // Output iwd Station diagnostic info
  fmt.Println(station[0].GetDiagnostics())
}
```

#### Examples

A number of ready-to-run examples demonstrating various use cases of go-iwd are available in the [go-iwd examples](https://github.com/shtirlic/go-iwd/tree/main/examples) dir.

## Features
- [x] Minimal dependencies
- [x] Easy API access
- [ ] Full API support + Experimental iwd API
- [ ] TUI Client
- [ ] API Tests

### IWD API
- [x] Adapter
- [x] Daemon
- [x] Device
- [x] KnowNetwork
- [x] Network
- [x] Station
- [x] Station Diagnostic
- [x] WSC
- [ ] Access Point
- [ ] Adhoc
- [ ] Agent
- [ ] Device Provisioning
- [ ] RadioManager
- [ ] RuleManager
- [ ] P2P (peer, service)
- [ ] Station Debug
- [ ] IWD specific error handling
- [ ] D-BUS Signals

## iwd Architecture

![iwd Architecture](https://iwd.wiki.kernel.org/_media/wiki/iwd-architecture.png)

## Links
- [iwd Project page](https://iwd.wiki.kernel.org/)
- [iwd Git repo](https://git.kernel.org/pub/scm/network/wireless/iwd.git)
- [iwd ArchWiki](https://wiki.archlinux.org/title/iwd)
