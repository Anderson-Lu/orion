package build

import (
	"fmt"
)

var (
	BuildVersion   = ""
	BuildTime      = ""
	BuildCommit    = ""
	BuildGoVersion = ""
	UitVersion     = "dev0.0.2"
)

func PrintVerbose() {
	fmt.Printf("Git Branch     : %s \n", BuildVersion)
	fmt.Printf("Git Commit     : %s \n", BuildCommit)
	fmt.Printf("Built Time     : %s \n", BuildTime)
	fmt.Printf("Go Version     : %s \n", BuildGoVersion)
	fmt.Printf("Uit Version  : %s \n", UitVersion)
}
