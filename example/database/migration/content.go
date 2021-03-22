package migration

import "embed"

//go:embed *.sql
var Schema embed.FS
