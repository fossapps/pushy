### Pushy.me strongly typed SDK for golang
<p align="center"><img src="./.github/gopher.png" width="200" /></p>

[![Go Report Card](https://goreportcard.com/badge/github.com/cyberhck/pushy)](https://goreportcard.com/report/github.com/cyberhck/pushy)[![Build Status](https://travis-ci.org/cyberhck/pushy.svg?branch=master)](https://travis-ci.org/cyberhck/pushy)

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
