package tsdb

import (
	"embed"

	"info/internal/infrastructure/repository/tsdb"
)

//go:embed *.sql
var EmbedMigrations embed.FS

var CurrentRepo *tsdb.Repository
