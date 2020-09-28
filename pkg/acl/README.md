# acl
Generic implementation of an Access Control List system for enforcing authorization across different resources. While originally meant for utilization in the [Drago](https://github.com/seashell/drago) project, this ACL system was built with flexibility in mind, and should be straightforward to configure and integrate into other codebases.

## Concepts
* **Resource**: any type of resource that might require authorization. Example: book.

* **Instance**: an uniquely identifiable instance of a resource. Example: "The lord of the rings"

* **Capability**: an operation that can be performed in the context of a resource or instance. In other words, a capability is an allowed operation. From an acessor perspective, it makes sense to talk in terms of capabilities, whereas from a resource perspective, its more usual to use the term operation. 

* **Policy**: a named set of rules for enabling/disabling capabilities on one or more instances.
  
* **Rule**: a combination of a filter for targetting specific resources/instances, as well as a set of capabilities that should be enabled on them.

* **Alias**: a convenience shorthand for representing multiple capabilities with a single string.

* **Secret**: an object in possession of the acessor which can be associated to a set of policies. A secret can be anything, from an opaque tokens to a JWT, as long as it is possible to resolve it.
  
* **ACL**: Object containing a compiled version of the access rules associated with a secret, and that can be used to find out whether a specific operation can be performed on a resource/instance.


## Usage

First one needs to define the resources requiring protection, as well as the available operations/capabilities. Let's assume we're building a backend for an online library, and would like to limit access to our books API.

We start by defining the capabilities:

```go
const (
	capBookRead   = "read"
	capBookWrite  = "write"
	capBookList   = "list"
)
```

Then we create a `acl.Config` object:

```go
var c = &Config{
	Resources: map[string]*Resource{
        "book": &Resource{
			Name: "book",
			Capabilities: map[string]*Capability{
				capBookRead:  {capBookRead, "Allows a book to be read"},
				capBookWrite: {capBookWrite, "Allows a book to be written"},
				capBookList:  {capBookList, "Allows books to be listed"},
			},
			Aliases: map[string]*CapabilityAlias{
				"write": {"read", []string{capBookWrite, capBookRead, capBookList}},
				"read":  {"read", []string{capBookRead, capBookList}},
			},
		},
    },
}
```

Note that besides the capabilities associated to the `book` resource, we have also defined aliases, which are basically shorthands for expressing more than one capability at once. Since our ACL system always expands all aliases prior to processing the actual capabilities, we can have an alias and a capability sharing the same name.

Finally, we need to tell our ACL system how it can resolve policy names and secrets. In order to decouple the ACL system from the persistence layer, we make use of the `acl.Policy` and `acl.Token` interfaces for abstracting these two objects.

Both functions can be set in the configuration struct, as follows:

```go
c.ResolvePolicy = func(policy string) (acl.Policy, error) {
    return someStorage.GetPolicyByName(policy)
}

c.ResolveToken = func(secret string) (acl.Token, error) {
    return someStorage.FindTokenBySecret(secret)
}
```

Now we can initialize our ACL:

```go
// Initialize ACL
Initialize(c)
```

When we receive a request to perform an operation of any kind, say `write` on a specific `book`, we can build a new ACL object based on the secret contained in this request, and query the ACL if the operation is authorized or not:

```go
acl, err := NewACLFromSecret("secret-in-the-request")
if err != nil {
    panic(err)
}

if !acl.IsAuthorized("book", "the-lord-of-the-rings-id", "write") {
    return fmt.Errorf("not authorized")
} 

// Proceed normally with the operation.
...
```





## Credits
Although different in many aspects, this implementation was heavily inspired on the ACL systems built by Hashicorp for Nomad and Consul, as well as on concepts from GCP's ACL system.