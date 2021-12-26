package root

import (
	"embed"
)

//go:embed *
var Static embed.FS
