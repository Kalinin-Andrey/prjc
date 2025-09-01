package tsdb

import (
	"embed"

	"course/internal/infrastructure/repository/tsdb"
)

//go:embed *.sql
var EmbedMigrations embed.FS

var CurrentRepo *tsdb.Repository
