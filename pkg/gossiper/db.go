package gossiper

import (
	. "github.com/Never-M/MyGossip/pkg/types"
	"github.com/syndtr/goleveldb/leveldb"
)

type mydb struct {
	path  string
	db    *leveldb.DB
	batch *batch
}

type pair struct {
	key string
	val string
}

func Newdb(path string) (int, error, *mydb) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return DB_CREATE_ERROR, err, nil
	}
	return SUCCEED, nil, &mydb{
		path:  path,
		db:    db,
		batch: &batch{},
	}
}

func (d *mydb) Add(key, val string) (int, error) {
	err := d.db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		return DB_PUT_ERROR, err
	}
	return SUCCEED, nil
}

func (d *mydb) Obtain(key string) (int, string, error) {
	data, err := d.db.Get([]byte(key), nil)
	if err != nil {
		return DB_GET_ERROR, "", err
	}
	return SUCCEED, string(data), nil
}

func (d *mydb) Remove(key string) (int, error) {
	err := d.db.Delete([]byte(key), nil)
	if err != nil {
		return DB_DELETE_ERROR, err
	}
	return SUCCEED, nil
}

func (d *mydb) ListData() (int, []pair, error) {
	iter := d.db.NewIterator(nil, nil)
	var ans []pair

	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		ans = append(ans, pair{key: key, val: value})
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return DB_DELETE_ERROR, nil, err
	}
	return SUCCEED, ans, nil
}

type batch struct {
	leveldb.Batch
}

func (b *batch) Set(keyval []pair) {
	for i := 0; i < len(keyval); i++ {
		b.Put([]byte(keyval[i].key), []byte(keyval[i].val))
	}
}

func (b *batch) Remove(keys []string) {
	for i := 0; i < len(keys); i++ {
		b.Delete([]byte(keys[i]))
	}
}

func (d *mydb) Commit(mybatch *leveldb.Batch) (int, error) {
	err := d.db.Write(mybatch, nil)
	if err != nil {
		return DB_BATCHWRITE_ERROR, err
	}
	return SUCCEED, nil
}
