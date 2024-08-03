package iface

type IUserIdentity interface {
	GetUsername() string
	GetRoles() []string
	HasRole(role string) bool
}
