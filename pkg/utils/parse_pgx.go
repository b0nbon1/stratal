package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(id string) (pgtype.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID string: %w", err)
	}

	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	return pgUUID, nil
}

func ParseText(text string) pgtype.Text {
	return pgtype.Text{
		String: text,
		Valid: true,
	}
}
