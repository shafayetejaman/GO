package main

type User struct {
	Name string
	Membership
}

func newUser(name string, membershipType string) User {
	// ?
	limit := 100
	if membershipType == "premium" {
		limit = 1000
	}

	return User{Name: name, Membership: Membership{
		Type:             membershipType,
		MessageCharLimit: limit,
	}}
}

type Membership struct {
	Type             string
	MessageCharLimit int
}
