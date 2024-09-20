package build

import (
	"fmt"
)

var (
	BuildVersion   = ""
	BuildTime      = ""
	BuildCommit    = ""
	BuildGoVersion = ""
)

func init() {
	fmt.Print("-----\n")
	fmt.Printf("XGRPC Branch     : %s \n", BuildVersion)
	fmt.Printf("XGRPC Commit     : %s \n", BuildCommit)
	fmt.Printf("XGRPC BuiltTime  : %s \n", BuildTime)
	fmt.Printf("XGRPC GO Version : %s \n", BuildGoVersion)
	fmt.Print("-----\n")
}
