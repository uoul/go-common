package auth

import "net/http"

type IAuthenticator[T any] interface {
	GetIdentityFromAuthorizationHeader(httpHeader http.Header) (T, error)
	GetIdentityFromAccessToken(token string) (T, error)
}
