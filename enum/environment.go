package enum

//go:generate enumer -type=Environment -text -json -trimprefix Environment -transform=snake -output environment_enumer.go -comment "app environment"
type Environment uint8

const (
	EnvironmentDev Environment = iota
	EnvironmentProd
	EnvironmentStage
)
