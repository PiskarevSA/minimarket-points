package postgresql

import (
	"context"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/google/uuid"
)

func (s *PostgreSql) GetBalance(
	ctx context.Context,
	userId uuid.UUID,
) (objects.Balance, error) {
	return objects.NilBalance, nil
}
