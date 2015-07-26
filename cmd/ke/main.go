package main // import "kego.io/cmd/ke"

import (
	"fmt"
	"os"

	"kego.io/process"
	_ "kego.io/system"
)

func main() {
	set, err := process.InitialiseAutomatic()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := process.KeCommand(set); err != nil {
		fmt.Println(process.FormatError(err))
		os.Exit(1)
	}
	os.Exit(0)
}
