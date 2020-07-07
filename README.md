# WebBridge
Connecting HTTP servers and clients on disparate networks using WebRTC and blockchain signaling

### Running with dev and debug mode
```
go run -race . -dev -debug
```

### Building
```
go build -i -v -ldflags="-X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(Get-Date -Format "yy.MM.dd")'" github.com/duality-solutions/web-bridge
```

#### Windows NMake
```
nmake /f Makefile
```

#### Linux Make
```
make
```

### Diagrams
![General Diagram](docs/diagram-webbridge-general.png)

![Technical Details Diagram](docs/diagram-webbridge-tech-details.png)

### License
WebBridge: See [LICENSE.md](./LICENSE.md "LICENSE.md") file for copying and use information.
Socks Library: Copyright (c) 2012-present, Samuel Stauffer <samuel@descolada.com>
RPC-Client Library: Copyright (c) 2013-present The btcsuite developers and Copyright (c) 2015-2016 The Decred developers
GoProxy Library: Copyright (c) 2012-present Elazar Leibovich. All rights reserved.
