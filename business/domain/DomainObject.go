// Package domain defines domain object used in Tenant service
package domain

// Tenant defines how a tenant should look like
type Tenant struct {
	SecretKey string
}

// Application defines how a application should look like
type Application struct {
	Name string
}
