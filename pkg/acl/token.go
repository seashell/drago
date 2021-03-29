package acl

// Token ...
type Token interface {
	IsPrivileged() bool
	Policies() []string
}
