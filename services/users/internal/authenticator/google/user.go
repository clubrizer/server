package google

// A User represents a Google user.
type User struct {
	Issuer     string
	ID         string
	GivenName  string
	FamilyName string
	Email      string
	Picture    string
}
