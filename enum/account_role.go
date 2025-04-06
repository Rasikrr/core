package enum

//go:generate enumer -type=AccountRole -json -trimprefix AccountRole -transform=snake -output account_role_enumer.go
type AccountRole uint8

const (
	AccountRoleAdmin AccountRole = iota
	AccountRoleUser
)

func (a AccountRole) OneOf(roles ...AccountRole) bool {
	for _, role := range roles {
		if a == role {
			return true
		}
	}
	return false
}
