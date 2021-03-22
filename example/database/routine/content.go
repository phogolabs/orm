package routine

import "embed"

//go:embed *.sql
var Query embed.FS
