package idp

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/alexedwards/scs/v2"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"

	"github.com/fensak-io/gostd/primitiveptr"
	"github.com/fensak-io/gostd/webstd"
	"github.com/fensak-io/gostd/webstd/chistd"
	"github.com/fensak-io/gostd/webstd/idp"
)

type AADB2C struct {
	clt          *msgraphsdk.GraphServiceClient
	appURL       url.URL
	tenantDomain string
	logger       *zap.SugaredLogger

	// Store references to OIDC auth and session manager for obtaining a valid ID token. These are used for logout.
	auth    *webstd.Authenticator
	sessMgr *scs.SessionManager
}

// Make sure AADB2C struct adheres to the idp.Service interface.
var _ idp.Service = (*AADB2C)(nil)

func NewAADB2C(
	logger *zap.Logger,
	appURL url.URL,
	idpCfg *webstd.IdP,
	oidc *webstd.OIDCProvider,
	auth *webstd.Authenticator,
	sm *scs.SessionManager,
) (*AADB2C, error) {
	sugar := logger.Sugar()

	cred, err := azidentity.NewClientSecretCredential(
		idpCfg.TenantID,
		oidc.ClientID,
		oidc.ClientSecret,
		nil,
	)
	if err != nil {
		return nil, err
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(
		cred,
		[]string{"https://graph.microsoft.com/.default"},
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the tenant domain based on the ID
	dsRaw, err := client.Domains().Get(context.Background(), nil)
	if err != nil {
		odataerr, isODataErr := err.(*odataerrors.ODataError)
		if isODataErr {
			msg := odataerr.GetError().GetMessage()
			sugar.Errorf("Error getting domain information from IdP: %s", primitiveptr.StringDeref(msg))
		}
		return nil, err
	}
	ds := dsRaw.GetValue()
	d := findMainTenantDomain(ds)
	if d == "" {
		doms := make([]string, 0, len(ds))
		for _, dom := range ds {
			doms = append(doms, primitiveptr.StringDeref(dom.GetId()))
		}
		err := fmt.Errorf("Could not find main tenant domain from list of domains: %v", doms)
		sugar.Error(err.Error())
		return nil, err
	}

	return &AADB2C{
		clt:          client,
		appURL:       appURL,
		auth:         auth,
		sessMgr:      sm,
		tenantDomain: d,
		logger:       sugar,
	}, nil
}

func (a AADB2C) RemoveUser(ctx context.Context, userID string) error {
	err := a.clt.UsersById(userID).Delete(ctx, nil)
	if err != nil {
		odataerr, isODataErr := err.(*odataerrors.ODataError)
		if isODataErr {
			msg := odataerr.GetError().GetMessage()
			a.logger.Errorf("Error creating user in IdP: %s", primitiveptr.StringDeref(msg))
		}
	}
	return err
}

func (a AADB2C) AddUser(ctx context.Context, profile idp.UserProfile) (string, error) {
	user := a.newUserFromProfile(profile)

	// Set default values to make the account creation successful
	accountEnabled := true
	user.SetAccountEnabled(&accountEnabled)
	user.SetCreationType(primitiveptr.String("LocalAccount"))

	// Set a random password and force the user to change it. Note that the user will go through the forgot password flow
	// on initial login, so this is not strictly necessary, but we set it to anyway avoid potential login loopholes.
	genp, err := password.Generate(20, rand.Intn(5)+1, rand.Intn(5)+1, false, false)
	if err != nil {
		return "", err
	}
	passwordProfile := graphmodels.NewPasswordProfile()
	forceChangePasswordNextSignIn := true
	passwordProfile.SetForceChangePasswordNextSignIn(&forceChangePasswordNextSignIn)
	passwordProfile.SetPassword(&genp)
	user.SetPasswordProfile(passwordProfile)

	result, err := a.clt.Users().Post(ctx, user, nil)
	if err != nil {
		odataerr, isODataErr := err.(*odataerrors.ODataError)
		if isODataErr {
			msg := odataerr.GetError().GetMessage()
			a.logger.Errorf("Error creating user in IdP: %s", primitiveptr.StringDeref(msg))
		}
		return "", err
	}
	return primitiveptr.StringDeref(result.GetId()), nil
}

func (a AADB2C) UpdateUser(ctx context.Context, profile idp.UserProfile) error {
	user := a.newUserFromProfile(profile)

	_, err := a.clt.UsersById(profile.GetID()).Patch(ctx, user, nil)
	if err != nil {
		odataerr, isODataErr := err.(*odataerrors.ODataError)
		if isODataErr {
			msg := odataerr.GetError().GetMessage()
			a.logger.Errorf("Error creating user in IdP: %s", primitiveptr.StringDeref(msg))
		}
	}
	return err
}

func (a AADB2C) GetLogoutURL(ctx context.Context) (string, error) {
	idpLogoutURL := a.auth.LogoutURL()
	if idpLogoutURL == "" {
		return "", fmt.Errorf("Missing expected end_session_endpoint in OIDC discovery claims")
	}

	idpLogoutURLParsed, err := url.Parse(idpLogoutURL)
	if err != nil {
		return "", err
	}

	refreshToken := a.sessMgr.GetString(ctx, chistd.RefreshTokenSessionKey)
	rawIDToken, _, err := a.auth.RefreshIDToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	appURLCopy := a.appURL
	appURLCopy.Path = chistd.OIDCLogoutPath
	qp := url.Values{}
	// Refer to the following document for info on these parameters:
	// https://learn.microsoft.com/en-us/azure/active-directory-b2c/session-behavior?pivots=b2c-user-flow#sign-out
	qp.Set("post_logout_redirect_uri", appURLCopy.String())
	qp.Set("id_token_hint", rawIDToken)
	idpLogoutURLParsed.RawQuery = qp.Encode()

	return idpLogoutURLParsed.String(), nil
}

func findMainTenantDomain(ds []graphmodels.Domainable) string {
	for _, d := range ds {
		did := primitiveptr.StringDeref(d.GetId())
		if strings.HasSuffix(did, "onmicrosoft.com") {
			return did
		}
	}
	return ""
}

func (a AADB2C) newUserFromProfile(profile idp.UserProfile) *graphmodels.User {
	requestBody := graphmodels.NewUser()

	fn := profile.GetFirstName()
	ln := profile.GetLastName()
	requestBody.SetGivenName(&fn)
	requestBody.SetSurname(&ln)
	displayName := fn + " " + ln
	requestBody.SetDisplayName(&displayName)

	localId := graphmodels.NewObjectIdentity()
	localId.SetIssuer(&a.tenantDomain)
	localId.SetSignInType(primitiveptr.String("emailAddress"))
	localId.SetIssuerAssignedId(
		primitiveptr.String(profile.GetPrimaryEmailAddress()),
	)
	requestBody.SetIdentities([]graphmodels.ObjectIdentityable{localId})

	return requestBody
}
