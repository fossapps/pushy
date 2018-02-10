### Pushy.me strongly typed SDK for golang
<p align="center"><img src="./.github/gopher.png" width="200" /></p>

Visit godoc for documentation

## Installation
```
go get github.com/cyberhck/pushy
```

Usage:
```go
sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Second))
res, _, _ := sdk.DeviceInfo("DEVICE_ID")
log.Println(res)
```
