package types

const (
	// ModuleName defines the module name
	ModuleName = "icq"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_icq"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	PendingICQInstanceKey      = "PendingICQInstance-value-"
	PendingICQInstanceCountKey = "PendingICQInstance-count-"

	PeriodicLastDataPointIdKey      = "PeriodicLastDataPointId-value-"
	PeriodicLastDataPointIdCountKey = "PeriodicLastDataPointId-count-"

	DataPointKey = "DataPoint-value-"

	ICQTimeoutsKey = "ICQTimeouts-value-"

	PeriodicICQKey      = "PeriodicICQ-value-"
	PeriodicICQCountKey = "PeriodicICQ-count-"
)
