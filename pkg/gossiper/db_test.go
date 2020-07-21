package gossiper

import (
	"github.com/Never-M/MyGossip/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBAddDelete(t *testing.T) {
	resultCode, mydb := Newdb("/tmp/test")
	assert.Equal(t, types.SUCCEED, resultCode, "Create Database Failed")
	resultCode = mydb.Add("hello", "world")
	assert.Equal(t, types.SUCCEED, resultCode, "Database Add Failed")
	resultCode, val := mydb.Obtain("hello")
	assert.Equal(t, types.SUCCEED, resultCode, "Can't Obtain key")
	assert.Equal(t, "world", val, "Database Obtain Failed")
	resultCode = mydb.Remove("hello")
	assert.Equal(t, types.SUCCEED, resultCode, "Remove Failed")
	resultCode, val = mydb.Obtain("hello")
	assert.Equal(t, types.DB_GET_ERROR, resultCode, "Remove Failed")
}

func TestListData(t *testing.T) {
	resultCode, mydb := Newdb("/tmp/test1")
	assert.Equal(t, types.SUCCEED, resultCode, "Create Database Failed")
	resultCode = mydb.Add("hello", "world")
	assert.Equal(t, types.SUCCEED, resultCode, "Database Add Failed")
	resultCode = mydb.Add("good", "day")
	assert.Equal(t, types.SUCCEED, resultCode, "Database Add Failed")
	resultCode, res := mydb.ListData()
	assert.Equal(t, types.SUCCEED, resultCode, "Database listData Failed")
	assert.Equal(t, len(res), 2, "ListData Failed")
}

func TestBatch(t *testing.T) {
	resultCode, mydb := Newdb("/tmp/test2")
	assert.Equal(t, types.SUCCEED, resultCode, "Create Database Failed")
	mydb.mybatch.Set([]string{"hello", "good"}, []string{"world", "day"})
	mydb.mybatch.Remove([]string{"hello"})
	resultCode = mydb.WriteDown(mydb.mybatch.batch)
	assert.Equal(t, types.SUCCEED, resultCode, "Batch Write Failed")
}
