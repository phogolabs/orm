package routine

import "embed"

//go:embed *.sql
var Statement embed.FS
