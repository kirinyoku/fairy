// Package assets provides embedded binary data for the fairy library.
package assets

import "embed"

// DataFS contains the embedded JSON metadata files.
//
//go:embed data/*.json
var DataFS embed.FS
