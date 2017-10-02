package bedrock

import (
	"fmt"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Level struct {
	Path string
}

func LoadLevel(path string) (*Level, error) {
	return &Level{
		Path: path,
	}, nil
}

func (l *Level) OpenDB() (*DB, error) {
	return OpenDB(filepath.Join(l.Path, "db"))
}

type DB struct {
	db *leveldb.DB
}

func OpenDB(path string) (*DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Put(key, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *DB) Close() error {
	err := db.db.CompactRange(util.Range{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db.db.Close()
}
