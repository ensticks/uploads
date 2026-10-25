// Package migrations embeds all SQL migration files so that the binary can
// apply them at startup without needing access to the source tree at runtime.
package migrations

import "embed"

// FS holds all *.sql files from this directory.
//
//go:embed *.sql
var FS embed.FS
