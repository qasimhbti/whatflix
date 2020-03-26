package envutils

import "github.com/pkg/errors"

const (
	// Testing is the testing environment
	Testing = "testing"
	// Development is the development environment
	Development = "development"
	// Staging is the staging environment
	Staging = "staging"
	// Production is the production environment
	Production = "production"
)

// Check checks if an environment value is valid
func Check(env string) error {
	switch env {
	case Testing, Development, Staging, Production:
		return nil
	default:
		return errors.Errorf("invalid: %s", env)
	}

}
