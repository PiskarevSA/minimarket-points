package handlers

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/PiskarevSA/minimarket-points/pkg/middlewares/jwtauth"
)

var getJwtFromContext = func(ctx context.Context, op string) (*jwt.Token, bool) {
	token, ok := jwtauth.JwtCtxKey.ValueOk(ctx)
	if !ok {
		log.Error().
			Str("op", op).
			Str("layer", "handler").
			Msg("failed to get jwt from ctx")
	}

	return token, ok
}

var getUserIdFromJwt = func(token *jwt.Token, op string) (uuid.UUID, error) {
	sub, err := token.Claims.GetSubject()
	if err != nil {
		log.Error().
			Str("op", op).
			Str("layer", "handler").
			Msg("failed to get user id from jwt")
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		log.Error().
			Err(err).
			Str("op", op).
			Str("layer", "handler").
			Msg("failed to parse sub to user id")
	}

	return userId, err
}
