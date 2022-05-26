# sona

A [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) implementation with Log, Payload, Authentication interceptor.


usage:
```
package main

import (
        "flag"
        "os"

        pb "service/gen/pb"
        "github.com/xiviu123/sona"
)

var(
  authServerEndpoint = flag.String("auth-server-endpoint", "localhost:9090", "Authenticate server endpoint")

  apiServerEndpoint = flag.String("api-server-endpoint", "localhost:9091", "API server endpoint")
)

func main() {
        server := sona.NewGateway()
        server.AddServiceHandle(authServerEndpoint, pb.RegisterAuthenticationServiceHandlerFromEndpoint)
        server.AddServiceHandle(apiServerEndpoint, pb.RegisterApiServiceHandlerFromEndpoint)
        if err := server.Start(":8080"); err != nil {
                panic(err)
        }

}

```
