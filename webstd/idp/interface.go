package idp

import (
	"context"
)

// Service represents the identity provider.
type Service interface {
	// RemoveUser will remove a user from the identity provider so that the user can no longer access the service.
	RemoveUser(ctx context.Context, userID string) error

	// AddUser will create a new user on the identity provider. This should return the associated external ID in the
	// identity provider.
	AddUser(ctx context.Context, profile UserProfile) (string, error)

	// UpdateUser will update the user profile information in the identity provider.
	UpdateUser(ctx context.Context, profile UserProfile) error

	// GetLogoutURL will return the session logout URL for the IdP.
	GetLogoutURL(ctx context.Context) (string, error)

	// ResendInviteEmail will resend the invite email for the user from the IdP. This is useful if the activation link is
	// expired on the IdP.
	ResendInviteEmail(ctx context.Context, userID string) error
}

// UserProfile represents the profile information of a user.
type UserProfile interface {
	// GetID should return the object ID of the user to uniquely identify the user on the identity provider.
	GetID() string

	// GetFirstName should return the first name of the user.
	GetFirstName() string

	// GetLastName should return the last name of the user.
	GetLastName() string

	// GetPrimaryEmailAddress should return the primary email address of the user.
	GetPrimaryEmailAddress() string
}
