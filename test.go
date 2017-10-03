package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/crystalmine/levels/anvil"
	"github.com/crystalmine/levels/bedrock"
	"github.com/crystalmine/levels/convert"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// printmap()

	// conv()

	// return

	if true {
		if true {
			if err := convert.PC2PE("maps/anvil-kathal", "maps/mcpe-conv"); err != nil {
				panic(err)
			}
		} else {
			cp := bedrock.Scan("maps/mcpe-conv/db")
			draw(cp)
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

func conv() {
	pe, err := bedrock.LoadLevel("maps/mcpe-conv")
	if err != nil {
		panic(err)
	}

	db, err := pe.OpenDB()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	for x := 0; x < 10; x++ {
		for z := 0; z < 10; z++ {
			if (x+z)&1 == 1 {
				if err := convert.FillChunk(db, bedrock.ChunkPos{x, z}); err != nil {
					panic(err)
				}
			}
		}
	}
}

func printmap() {
	l, err := anvil.LoadLevel("maps/anvil-kathal")
	if err != nil {
		panic(err)
	}
	cp := make(map[anvil.ChunkPos]int)
	for _, ri := range l.Regions {
		r, err := l.LoadRegion(ri)
		if err != nil {
			panic(err)
		}
		log.Info().Interface("pos", r.RegionPos).Msgf("region loaded")
		for pos, c := range r.Chunks {
			cp[pos] = len(c.Data.Sections)
		}
	}
	draw(cp)

}

func draw(cp map[anvil.ChunkPos]int) {
	var minx, minz = math.MaxInt32, math.MaxInt32
	var maxx, maxz = math.MinInt32, math.MinInt32
	for pos := range cp {
		if minx > pos.X {
			minx = pos.X
		}
		if minz > pos.Z {
			minz = pos.Z
		}
		if maxx < pos.X {
			maxx = pos.X
		}
		if maxz < pos.Z {
			maxz = pos.Z
		}
	}

	log.Info().Int("minx", minx).Int("minz", minz).Msg("minimum")
	log.Info().Int("maxx", maxx).Int("maxz", maxz).Msg("maximum")

	dx, dz := maxx-minx+1, maxz-minz+1

	p := make([]int, dx*dz)
	for pos, v := range cp {
		p[(pos.X-minx)*dz+(pos.Z-minz)] = v
	}

	for x := 0; x < dx; x++ {
		fmt.Printf("%3d ", x+minx)
		for z := 0; z < dz; z++ {
			if p[x*dz+z] == 0 {
				fmt.Print(".")
			} else {
				fmt.Print(strconv.FormatInt(int64(p[x*dz+z]-1), 16))
			}
		}
		fmt.Println()
	}
}
