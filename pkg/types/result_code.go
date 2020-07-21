package types

const (
	SUCCEED             = 0
	FAILED              = -1
	GOSSIPER_PEER_EXIST = 1

	HEARTBEAT_REQUEST_ERROR  = 11
	HEARTBEAT_RESPONSE_ERROR = 12

	DB_CREATE_ERROR     = 21
	DB_PUT_ERROR        = 22
	DB_GET_ERROR        = 23
	DB_DELETE_ERROR     = 24
	DB_BATCHWRITE_ERROR = 25
)