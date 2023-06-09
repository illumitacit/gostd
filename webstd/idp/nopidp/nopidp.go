package nopidp

import (
	"context"

	"github.com/fensak-io/gostd/webstd/idp"
)

// NOPIdP is an IdP service provider that does nothing (no-op). This is most useful for testing.
type NOPIdP struct{}

// Make sure NOPIdP struct adheres to the idp.Service interface.
var _ idp.Service = (*NOPIdP)(nil)

func (s NOPIdP) RemoveUser(ctx context.Context, userID string) error {
	return nil
}

func (s NOPIdP) AddUser(ctx context.Context, profile idp.UserProfile) (string, error) {
	return "", nil
}

func (s NOPIdP) UpdateUser(ctx context.Context, profile idp.UserProfile) error {
	return nil
}

func (s NOPIdP) GetLogoutURL(ctx context.Context) (string, error) {
	return "", nil
}

func (s NOPIdP) ResendInviteEmail(ctx context.Context, userID string) error {
	return nil
}
