package api

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	session "github.com/chaostreff-flensburg/moc/context"
	"github.com/chaostreff-flensburg/moc/router"
)

// add token informations to context
func (api *API) withToken(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()

	bearerToken, err := extractBearerToken(r)
	if err != nil {
		return ctx, nil
	}

	if api.config.OperatorToken == "" && bearerToken != api.config.OperatorToken {
		return ctx, nil
	}

	ctx = session.SetOperator(ctx)
	return ctx, nil
}

// authRequired check if user has a valid jwt token or is operator
func authRequired(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()

	isOperator := ctx.Value("isOperator")

	if isOperator == nil {
		return nil, router.UnauthorizedError("Authorization error")
	}

	return ctx, nil
}

// extractBearerToken from Request
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	bearerRegexp := regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("Bad authentication header")
	}

	return matches[1], nil
}
