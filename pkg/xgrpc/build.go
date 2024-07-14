package xgrpc

import (
	"fmt"
)

var (
	XGrpcVersion            = "v0.0.1"
	XGrpcServerBuildVersion = ""
)

func init() {
	fmt.Print("=============\n")
	fmt.Printf("XGrpc Version: %s \n", XGrpcVersion)
	fmt.Printf("XGrpc Service Build At: %s \n", XGrpcServerBuildVersion)
	fmt.Print("=============\n")
}
