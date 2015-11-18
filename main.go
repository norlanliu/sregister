package main

import (
	"github.com/norlanliu/sregister/sregistermain"
	"os"
)

func main() {
	sr := sregistermain.NewSRegister(os.Args[1:])
	sr.Run()
}
