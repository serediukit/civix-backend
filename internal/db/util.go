package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("db: no rows in result set")

func SB() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
