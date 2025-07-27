package pyinit

import "embed"

//go:embed templates/*
var EmbeddedTemplates embed.FS
