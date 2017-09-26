package main

import (
	"github.com/crystalmine/mapper/bedrock"
)

func main() {
	bedrock.Scan("d:/mcpe/world/db")
	bedrock.ReadNbtFile("d:/mcpe/world/level.dat")

}
