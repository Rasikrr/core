package enum

//go:generate enumer -type=LogFormat -text -json -trimprefix LogFormat -transform=snake -output log_format_enumer.go -comment "logger output format"

type LogFormat uint8

const (
	LogFormatText LogFormat = iota
	LogFormatJSON
)
