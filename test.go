package main

import (
	"os"

	"github.com/crystalmine/levels/anvil"
	"github.com/crystalmine/levels/bedrock"
	"github.com/crystalmine/levels/convert"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if true {
		if true {
			if err := convert.PC2PE("maps/anvil-kathal", "maps/mcpe-conv"); err != nil {
				panic(err)
			}
		} else {
			bedrock.Scan("maps/mcpe-convx/db")
			// bedrock.Scan("maps/mcpe-simple/db")
			// bedrock.ReadNbtFile("maps/mcpe-simple/level.dat")
		}
	} else {
		// anvil.Scan("maps/anvil-kathal")
		l, err := anvil.LoadLevel("maps/anvil-kathal")
		if err != nil {
			panic(err)
		}
		_ = l
		// pretty.Println(l)
		r, err := l.LoadRegion(l.Regions[0])
		if err != nil {
			panic(err)
		}
		_ = r
		// pretty.Println(r)
	}

}
