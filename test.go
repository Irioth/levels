package main

import (
	"github.com/crystalmine/levels/anvil"
	"github.com/crystalmine/levels/bedrock"
)

func main() {

	if false {
		bedrock.Scan("maps/mcpe-simple/db")
		bedrock.ReadNbtFile("maps/mcpe-simple/level.dat")
	} else {
		anvil.Scan("maps/anvil-kathal")
	}

}
