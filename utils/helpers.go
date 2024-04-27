package utils

import (
	"context"
	"errors"

	"github.com/SEC-Jobstreet/backend-employer-service/models"
)

const (
	AuthorizationPayloadKey = "authorization_payload"
)

func GetCurrentUser(ctx context.Context) (models.AuthClaim, error) {
	userCtx := ctx.Value(AuthorizationPayloadKey)
	currentUser, ok := userCtx.(models.AuthClaim)
	if !ok {
		return models.AuthClaim{}, errors.New("cannot parse to type User")
	}

	return currentUser, nil
}
