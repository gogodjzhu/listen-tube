package auth

// Authenticator authenticates the user
type Authenticator interface {
	// Authenticate authenticates the user
	Authenticate(username, password string) (bool, error)
}

// Authorize checks if the user is authorized to perform the action
type Authorizer interface {
	// Authorize checks if the user is authorized to perform the action
	Authorize(username, action string) (bool, error)
}
