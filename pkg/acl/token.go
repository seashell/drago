package acl

// Token ...
type Token interface {
	IsManagement() bool
	Policies() []string
}
