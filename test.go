package main

import (
	"github.com/crystalmine/mapper/bedrock"
)

func main() {
	bedrock.Scan("maps/mcpe-simple/db")
	bedrock.ReadNbtFile("maps/mcpe-simple/level.dat")

}
