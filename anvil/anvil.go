package anvil

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/crystalmine/levels/nbt"
	"github.com/kr/pretty"
)

func Scan(fname string) {
	scanFile(fname + "/level.dat")
	scanFile(fname + "/data/villages.dat")
	scanFile(fname + "/players/Player.dat")
	scanFile(fname + "/players/ScarmuzziBoy33.dat")
}

func scanFile(fname string) {
	data, err := readGzipped(fname)
	if err != nil {
		panic(err)
	}

	name, res, err := nbt.Decode(data, false)
	if err != nil {
		panic(err)
	}
	pretty.Println(name, res)

}

func readGzipped(fname string) ([]byte, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return raw, nil

}
