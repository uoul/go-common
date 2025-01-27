package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	authError "github.com/uoul/go-common/auth/error"
	"github.com/uoul/go-common/auth/iface"
)

// -------------------------------------------------------------------
// Constant's
// -------------------------------------------------------------------
const (
	AUTH_HEADER = "Authorization"
)

// -------------------------------------------------------------------
// Typedefinitions
// -------------------------------------------------------------------
type KeyCloakAuthenticator[T iface.IUserIdentity] struct {
	jwksUri string
	jwkSet  jwk.Set
}

// -------------------------------------------------------------------
// Public Methods/Functions
// -------------------------------------------------------------------

// This function creates a new instance of Oidc
//
// IN:
//   - jwksUri: URL of jwks Endpoint (for validating given token)
//
// OUT:
//   - IAuthenticator: new instance of IAuthenticator(means in this case Oidc)
func NewKeyCloakAuthenticator[T iface.IUserIdentity](jwksUri string) iface.IAuthenticator[T] {
	return &KeyCloakAuthenticator[T]{
		jwksUri: jwksUri,
		jwkSet:  nil,
	}
}

// This Method extracts a authorization-header of a given context and maps the
// result to an object of IUserIdentity.
//
// IN:
//   - httpHeader: The http request header, which should include the authorization header
//
// OUT:
//   - IUserIdentity: UserIdentity of type T
//   - error: if any error occures, the error out will report the issue
func (authenticator *KeyCloakAuthenticator[T]) GetIdentity(httpHeader http.Header) (T, error) {
	authHeader, found := httpHeader[AUTH_HEADER]
	if !found {
		return *new(T), authError.NewAuthenticationError("failed to get authentication header")
	}
	if rawToken, found := strings.CutPrefix(authHeader[0], "Bearer "); found {
		accessToken, err := jwt.Parse(rawToken, authenticator.keyFunc)
		if err != nil {
			return *new(T), authError.NewAuthenticationError("invalid token - %v", err)
		}
		claims := accessToken.Claims.(jwt.MapClaims)
		j, err := json.Marshal(claims)
		if err != nil {
			return *new(T), authError.NewAuthenticationError("failed to marshal token-claims - %v", err)
		}
		var customClaims T
		err = json.Unmarshal(j, &customClaims)
		if err != nil {
			return *new(T), authError.NewAuthenticationError("failed to unmarshal token-claims - %v", err)
		}
		return customClaims, nil
	} else {
		return *new(T), authError.NewAuthenticationError("invalid authorization header - header has to start with \"Bearer \"")
	}
}

// This Method maps the given token to an object of IUserIdentity
//
// IN:
//   - token: jwt-token
//
// OUT:
//   - IUserIdentity: UserIdentity of type T
//   - error: if any error occures, the error out will report the issue
func (authenticator *KeyCloakAuthenticator[T]) GetIdentityOfAccessToken(token string) (T, error) {
	accessToken, err := jwt.Parse(token, authenticator.keyFunc)
	if err != nil {
		return *new(T), authError.NewAuthenticationError("invalid token - %v", err)
	}
	claims := accessToken.Claims.(jwt.MapClaims)
	j, err := json.Marshal(claims)
	if err != nil {
		return *new(T), authError.NewAuthenticationError("failed to marshal token-claims - %v", err)
	}
	var customClaims T
	err = json.Unmarshal(j, &customClaims)
	if err != nil {
		return *new(T), authError.NewAuthenticationError("failed to unmarshal token-claims - %v", err)
	}
	return customClaims, nil
}

// -------------------------------------------------------------------
// Private helper methods/functions
// -------------------------------------------------------------------
func (authenticator *KeyCloakAuthenticator[T]) keyFunc(t *jwt.Token) (any, error) {
	set, err := authenticator.getJwkSet()
	if err != nil {
		return nil, err
	}
	keyID, ok := t.Header["kid"].(string)
	if !ok {
		return nil, authError.NewAuthenticationError("expecting JWT header to have string kid")
	}
	if key, found := set.LookupKeyID(keyID); found {
		var pubkey any
		err := key.Raw(&pubkey)
		if err != nil {
			return nil, err
		}
		return pubkey, nil
	}
	return nil, authError.NewAuthenticationError("unable to find key %q", keyID)
}

func (authenticator *KeyCloakAuthenticator[T]) getJwkSet() (jwk.Set, error) {
	if authenticator.jwkSet == nil {
		set, err := jwk.Fetch(context.Background(), authenticator.jwksUri)
		if err != nil {
			return nil, err
		}
		authenticator.jwkSet = set
	}
	return authenticator.jwkSet, nil
}
