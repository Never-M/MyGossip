package gossiper

import (
	. "github.com/Never-M/MyGossip/pkg/types"
	"github.com/syndtr/goleveldb/leveldb"
	)

type mydb struct {
	path	string
	db		*leveldb.DB
}

func newdb(path string) (int, *mydb) {
	newdb := &mydb{}
	newdb.path = path
	var err error
	newdb.db, err = leveldb.OpenFile("path/to/db", nil)
	if err != nil {
		return DB_CREATE_ERROR, nil
	}
	return SUCCEED, newdb
}

func put()  {

}

func get()  {

}