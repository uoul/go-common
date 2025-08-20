package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
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
type JwksAuthenticator[T any] struct {
	jwksUri         string
	jwkSet          jwk.Set
	lastJwksRefresh time.Time

	requestTimeout      time.Duration
	jwksRefreshInterval time.Duration
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
func NewJwksAuthenticator[T any](jwksUri string) IAuthenticator[T] {
	return &JwksAuthenticator[T]{
		jwksUri:         jwksUri,
		jwkSet:          nil,
		lastJwksRefresh: time.Now(),

		requestTimeout:      10 * time.Second,
		jwksRefreshInterval: 600 * time.Second,
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
func (a *JwksAuthenticator[T]) GetIdentityFromAuthorizationHeader(httpHeader http.Header) (T, error) {
	authHeader, found := httpHeader[AUTH_HEADER]
	if !found {
		return *new(T), fmt.Errorf("failed to get authentication header")
	}
	if rawToken, found := strings.CutPrefix(authHeader[0], "Bearer "); found {
		identity, err := a.GetIdentityFromAccessToken(rawToken)
		if err != nil {
			return *new(T), err
		}
		return identity, nil
	} else {
		return *new(T), fmt.Errorf("invalid authorization header - header has to start with \"Bearer \"")
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
func (a *JwksAuthenticator[T]) GetIdentityFromAccessToken(token string) (T, error) {
	// Get JsonWebKeySet
	jwkSet, err := a.getJwkSet()
	if err != nil {
		return *new(T), err
	}
	// Check if token is valid
	_, err = jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(jwkSet),
		jwt.WithVerify(true),
		jwt.WithValidate(true),
	)
	if err != nil {
		return *new(T), err
	}
	// Parse payload
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return *new(T), fmt.Errorf("invalid token - token must contain 3 parts split by '.'")
	}
	// Decode base64 payload
	claimsStr, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return *new(T), fmt.Errorf("failed to decode claims - %v", err)
	}
	var claims T
	err = json.Unmarshal(claimsStr, &claims)
	if err != nil {
		return *new(T), fmt.Errorf("failed to parse claims - %v", err)
	}
	return claims, nil
}

// -------------------------------------------------------------------
// Private helper methods/functions
// -------------------------------------------------------------------
func (a *JwksAuthenticator[T]) getJwkSet() (jwk.Set, error) {
	if a.jwkSet == nil || time.Since(a.lastJwksRefresh) > a.jwksRefreshInterval {
		ctx, cancel := context.WithTimeout(context.Background(), a.requestTimeout)
		defer cancel()
		set, err := jwk.Fetch(ctx, a.jwksUri)
		if err != nil {
			return nil, err
		}
		a.jwkSet = set
	}
	return a.jwkSet, nil
}
