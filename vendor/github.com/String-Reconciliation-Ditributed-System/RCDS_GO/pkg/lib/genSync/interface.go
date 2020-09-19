package genSync

import "github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/set"

type GenSync interface {
	SetFreezeLocal(freezeLocal bool)

	AddElement(elem interface{}) error
	DeleteElement(elem interface{}) error

	SyncClient(ip string, port int) error
	SyncServer(ip string, port int) error

	GetLocalSet() *set.Set
	GetSetAdditions() *set.Set 	// Set the set that is added to the local set.
	GetSentBytes() int
	GetReceivedBytes() int
	GetTotalBytes() int
}
