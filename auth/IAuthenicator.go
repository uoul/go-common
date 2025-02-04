package auth

import "net/http"

type IAuthenticator[T IUserIdentity] interface {
	GetIdentity(httpHeader http.Header) (T, error)
	GetIdentityOfAccessToken(token string) (T, error)
}
