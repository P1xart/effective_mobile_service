package schema

import "embed"

//go:embed migrations/*.sql
var DB embed.FS
