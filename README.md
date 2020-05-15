# WebBridge
Connecting HTTP within intranets using WebRTC and blockchain signaling

### Running with dev and debug mode
```
go run -race . -dev -debug
```

### Building
```
go build -i -v -ldflags="-X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(Get-Date -Format "yy.MM.dd")'" github.com/duality-solutions/web-bridge
```

### Diagrams
![General Diagram](docs/diagram-webbridge-general.png)

![Technical Details Diagram](docs/diagram-webbridge-tech-details.png)

### License
See [LICENSE.md](./LICENSE.md "LICENSE.md") file for copying and use information.