# GO OpenID RP (Relying Party)

### Project Summary

GO OpenID RP is a Go repository contains exmaple of code that works as an RP connecting to TAMA via OpenID Connect 1.0 Protocol

The example code uses [go-oidc] for OpenID Connect 1.0 library

### Installation
```bashp
go get ./...
```

### Run
```bashp
go run server.go
```

### Usage
To try using OpenID Connect, please go to
```bashp
http://localhost:3000/login
```

[go-oidc]: https://github.com/coreos/go-oidc 
