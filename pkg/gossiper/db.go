package gossiper

import (
	. "github.com/Never-M/MyGossip/pkg/types"
	"github.com/syndtr/goleveldb/leveldb"
)

type mydb struct {
	path    string
	db      *leveldb.DB
	mybatch *mybatch
}

type pair struct {
	key string
	val string
}

func Newdb(path string) (int, *mydb) {
	newdb := &mydb{}
	newdb.path = path
	var err error
	newdb.db, err = leveldb.OpenFile(newdb.path, nil)
	newdb.mybatch = NewBatch()
	if err != nil {
		return DB_CREATE_ERROR, nil
	}
	return SUCCEED, newdb
}

func (d *mydb) Add(key, val string) int {
	err := d.db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		return DB_PUT_ERROR
	}
	return SUCCEED
}

func (d *mydb) Obtain(key string) (int, string) {
	data, err := d.db.Get([]byte(key), nil)
	if err != nil {
		return DB_GET_ERROR, ""
	}
	return SUCCEED, string(data)
}

func (d *mydb) Remove(key string) int {
	err := d.db.Delete([]byte(key), nil)
	if err != nil {
		return DB_DELETE_ERROR
	}
	return SUCCEED
}

func (d *mydb) ListData() (int, []pair) {
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
		return DB_DELETE_ERROR, nil
	}
	return SUCCEED, ans
}

type mybatch struct {
	batch *leveldb.Batch
}

func NewBatch() *mybatch {
	return &mybatch{
		batch: new(leveldb.Batch),
	}
}

func (mb *mybatch) Set(keys, vals []string) {
	for i := 0; i < len(keys); i++ {
		mb.batch.Put([]byte(keys[i]), []byte(vals[i]))
	}
}

func (mb *mybatch) Remove(keys []string) {
	for i := 0; i < len(keys); i++ {
		mb.batch.Delete([]byte(keys[i]))
	}
}

func (d *mydb) WriteDown(mybatch *leveldb.Batch) int {
	err := d.db.Write(mybatch, nil)
	if err != nil {
		return DB_BATCHWRITE_ERROR
	}
	return SUCCEED
}
