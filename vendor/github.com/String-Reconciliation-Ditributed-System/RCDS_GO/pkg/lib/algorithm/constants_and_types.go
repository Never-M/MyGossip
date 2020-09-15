package algorithm

type SyncType string

const (
	// string reconciliation
	RCDS SyncType = "RecursiveContentDependentShingling"
	// set reconciliation
	IBLT     SyncType = "InvertibleBloomLookupTable"
	CPI      SyncType = "CharacteristicPolynomialInterpolation"
	InterCPI SyncType = "InteractiveCharacteristicPolynomialInterpolation"
	FullSync SyncType = "FullSync"
	// set difference estimation
	StrataEst SyncType = "StrataEstimation"
)
