package gossiper

import (
	"github.com/Never-M/MyGossip/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBAddDelete(t *testing.T) {
	resultCode, mydb, _ := Newdb("/tmp/test")
	assert.Equal(t, types.SUCCEED, resultCode, "Create Database Failed")
	resultCode, _ = mydb.Put("hello", "world")
	assert.Equal(t, types.SUCCEED, resultCode, "Database Put Failed")
	resultCode, val, _ := mydb.Get("hello")
	assert.Equal(t, types.SUCCEED, resultCode, "Can't Get key")
	assert.Equal(t, "world", val, "Database Get Failed")
	resultCode, _ = mydb.Delete("hello")
	assert.Equal(t, types.SUCCEED, resultCode, "Delete Failed")
	resultCode, val, _ = mydb.Get("hello")
	assert.Equal(t, types.DB_GET_ERROR, resultCode, "Get Failed")
	//every test case has a Close, deleting it later
	mydb.Close()
}

func TestListData(t *testing.T) {
	resultCode, mydb, _ := Newdb("/tmp/test1")
	assert.Equal(t, types.SUCCEED, resultCode, "Create Database Failed")
	resultCode, _ = mydb.Put("hello", "world")
	assert.Equal(t, types.SUCCEED, resultCode, "Database Put Failed")
	resultCode, _ = mydb.Put("good", "day")
	assert.Equal(t, types.SUCCEED, resultCode, "Database Put Failed")
	resultCode, res, _ := mydb.ListData()
	assert.Equal(t, types.SUCCEED, resultCode, "Database listData Failed")
	assert.Equal(t, len(res), 2, "ListData Failed")
	//every test case has a Close, deleting it later
	mydb.Close()
}

func TestBatch(t *testing.T) {
	resultCode, mydb, _ := Newdb("/tmp/test2")
	assert.Equal(t, types.SUCCEED, resultCode, "Create Database Failed")
	mydb.BatchPut([]pair{NewPair("hello", "world"), NewPair("good", "morning")})
	mydb.BatchDelete([]string{"hello"})
	resultCode, _ = mydb.Commit()
	assert.Equal(t, types.SUCCEED, resultCode, "Batch Commit Failed")
	//every test case has a Close, deleting it later
	mydb.Close()
}
