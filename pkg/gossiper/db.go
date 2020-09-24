package gossiper

import (
	"fmt"
	. "github.com/Never-M/MyGossip/pkg/types"
	"github.com/syndtr/goleveldb/leveldb"
)

type mydb struct {
	path  string
	db    *leveldb.DB
	batch *leveldb.Batch
}

type pair struct {
	key string
	val string
}

func NewPair(key, val string) pair {
	return pair{key: key, val: val}
}

func Newdb(path string) (int, *mydb, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return DB_CREATE_ERROR, nil, err
	}
	return SUCCEED, &mydb{
		path:  path,
		db:    db,
		batch: new(leveldb.Batch),
	}, nil
}

func (d *mydb) Put(key, val string) (int, error) {
	err := d.db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		return DB_PUT_ERROR, err
	}
	return SUCCEED, nil
}

func (d *mydb) Get(key string) (int, string, error) {
	data, err := d.db.Get([]byte(key), nil)
	if err != nil {
		return DB_GET_ERROR, "", err
	}
	return SUCCEED, string(data), nil
}

func (d *mydb) Delete(key string) (int, error) {
	err := d.db.Delete([]byte(key), nil)
	if err != nil {
		return DB_DELETE_ERROR, err
	}
	return SUCCEED, nil
}

func (d *mydb) Close() {
	d.db.Close()
}

func (d *mydb) ListData() (int, []pair, error) {
	iter := d.db.NewIterator(nil, nil)
	var ans []pair

	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		fmt.Println(key + " " + value)
		ans = append(ans, pair{key: key, val: value})
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return DB_DELETE_ERROR, nil, err
	}
	return SUCCEED, ans, nil
}

func (d *mydb) BatchPut(keyval []pair) {
	for i := 0; i < len(keyval); i++ {
		d.batch.Put([]byte(keyval[i].key), []byte(keyval[i].val))
	}
}

func (d *mydb) BatchDelete(keys []string) {
	for i := 0; i < len(keys); i++ {
		d.batch.Delete([]byte(keys[i]))
	}
}

func (d *mydb) Commit() (int, error) {
	err := d.db.Write(d.batch, nil)
	if err != nil {
		return DB_BATCHWRITE_ERROR, err
	}
	return SUCCEED, nil
}
