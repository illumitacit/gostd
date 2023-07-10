package zitadel

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/alexedwards/scs/v2"
	"github.com/zitadel/oidc/pkg/oidc"
	"go.uber.org/zap"

	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/fensak-io/gostd/webstd"
	"github.com/fensak-io/gostd/webstd/chistd"
	"github.com/fensak-io/gostd/webstd/idp"
)

type Zitadel struct {
	logger *zap.SugaredLogger
	appURL url.URL

	c *management.Client

	// Store references to OIDC auth and session manager for obtaining a valid ID token. These are used for logout.
	auth    *webstd.Authenticator
	sessMgr *scs.SessionManager
}

// Make sure Zitadel struct adheres to the idp.Service interface.
var _ idp.Service = (*Zitadel)(nil)

func NewZitadel(
	logger *zap.SugaredLogger,
	appURL url.URL,
	idpCfg *webstd.IdP,
	auth *webstd.Authenticator,
	sm *scs.SessionManager,
) (*Zitadel, error) {
	data, err := base64.StdEncoding.DecodeString(idpCfg.Zitadel.JWTKeyBase64)
	if err != nil {
		return nil, err
	}

	issuer := fmt.Sprintf("https://%s.zitadel.cloud", idpCfg.Zitadel.InstanceName)
	api := fmt.Sprintf("%s.zitadel.cloud:443", idpCfg.Zitadel.InstanceName)
	client, err := management.NewClient(
		issuer,
		api,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData(data)),
	)
	if err != nil {
		return nil, err
	}

	return &Zitadel{
		logger:  logger,
		appURL:  appURL,
		c:       client,
		auth:    auth,
		sessMgr: sm,
	}, nil
}

func (z Zitadel) RemoveUser(ctx context.Context, userID string) error {
	req := &pb.RemoveUserRequest{
		Id: userID,
	}
	if _, err := z.c.RemoveUser(ctx, req); err != nil {
		z.logger.Errorf("Error removing user from IdP: %s", err)
		return err
	}
	return nil
}

func (z Zitadel) AddUser(ctx context.Context, profile idp.UserProfile) (string, error) {
	req := &pb.AddHumanUserRequest{
		UserName: profile.GetPrimaryEmailAddress(),
		Profile: &pb.AddHumanUserRequest_Profile{
			FirstName: profile.GetFirstName(),
			LastName:  profile.GetLastName(),
		},
		Email: &pb.AddHumanUserRequest_Email{
			Email:           profile.GetPrimaryEmailAddress(),
			IsEmailVerified: false,
		},
	}
	resp, err := z.c.AddHumanUser(ctx, req)
	if err != nil {
		z.logger.Errorf("Error creating user in IdP: %s", err)
		return "", err
	}
	return resp.UserId, nil
}

func (z Zitadel) UpdateUser(ctx context.Context, profile idp.UserProfile) error {
	req := &pb.UpdateUserNameRequest{
		UserId:   profile.GetID(),
		UserName: profile.GetPrimaryEmailAddress(),
	}
	if _, err := z.c.UpdateUserName(ctx, req); err != nil {
		z.logger.Errorf("Error updating username for user %s in IdP: %s", req.UserId, err)
		return err
	}

	profUpdReq := &pb.UpdateHumanProfileRequest{
		UserId:    profile.GetID(),
		FirstName: profile.GetFirstName(),
		LastName:  profile.GetLastName(),
	}
	if _, err := z.c.UpdateHumanProfile(ctx, profUpdReq); err != nil {
		z.logger.Errorf("Error updating profile for user %s in IdP: %s", req.UserId, err)
		return err
	}

	return nil
}

func (z Zitadel) GetLogoutURL(ctx context.Context) (string, error) {
	idpLogoutURL := z.auth.LogoutURL()
	if idpLogoutURL == "" {
		return "", fmt.Errorf("Missing expected end_session_endpoint in OIDC discovery claims")
	}

	idpLogoutURLParsed, err := url.Parse(idpLogoutURL)
	if err != nil {
		return "", err
	}

	refreshToken := z.sessMgr.GetString(ctx, chistd.RefreshTokenSessionKey)
	rawIDToken, _, newToken, err := z.auth.RefreshIDToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	z.sessMgr.Put(ctx, chistd.AccessTokenSessionKey, newToken.AccessToken)
	z.sessMgr.Put(ctx, chistd.AccessTokenExpirySessionKey, newToken.Expiry)
	if newToken.RefreshToken != "" {
		z.sessMgr.Put(ctx, chistd.RefreshTokenSessionKey, newToken.RefreshToken)
	}

	appURLCopy := z.appURL
	appURLCopy.Path = chistd.OIDCLogoutPath
	qp := url.Values{}
	// Refer to the following document for info on these parameters:
	// https://zitadel.com/docs/guides/integrate/logout
	qp.Set("post_logout_redirect_uri", appURLCopy.String())
	qp.Set("id_token_hint", rawIDToken)
	idpLogoutURLParsed.RawQuery = qp.Encode()

	return idpLogoutURLParsed.String(), nil
}

func (z Zitadel) ResendInviteEmail(ctx context.Context, userID string) error {
	emailReq := &pb.GetHumanEmailRequest{UserId: userID}
	resp, err := z.c.GetHumanEmail(ctx, emailReq)
	if err != nil {
		z.logger.Errorf("Error looking up user %s email in IdP: %s", userID, err)
		return err
	}

	emailObj := resp.GetEmail()
	if emailObj == nil {
		z.logger.Errorf("Error looking up user %s email in IdP: email is empty", userID)
		return err
	}

	req := &pb.ResendHumanInitializationRequest{
		UserId: userID,
		Email:  emailObj.Email,
	}
	if _, err := z.c.ResendHumanInitialization(ctx, req); err != nil {
		z.logger.Errorf("Error resending initialization email for user %s: %s", userID, err)
		return err
	}
	return nil
}
