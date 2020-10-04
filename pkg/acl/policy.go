package acl

// Policy represents a named set of rules for allowing
// operations on specific resources.
type Policy interface {
	Name() string
	Rules() []Rule
}

// Rule is used to allow operations on specific resources. In
// addition to the name of the target resource type and the allowed
// operations, a rule also specifies an optional glob pattern for
// targeting specific instances of the resource.
type Rule interface {
	// Resource targeted by this rule.
	Resource() string
	// Pattern used to target specific instances
	// of the target resource, if applicable.
	Pattern() string
	// Capabilities contains the actions allowed on
	// instances of a resource matching this rule.
	Capabilities() []string
}
