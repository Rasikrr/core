package enum

//go:generate enumer -type=IsoLevel -text -json -trimprefix IsoLevel -transform=snake -output tx_iso_level_enumer.go -comment "transactions isolation levels"

type IsoLevel uint8

const (
	IsoLevelReadCommited IsoLevel = iota
	IsoLevelRepeatableRead
	IsoLevelSerializable
)
