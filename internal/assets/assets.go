// Package assets provides embedded binary data for the fairy library.
// NOTE: By using go:embed, we bundle the datamined game JSON files located in the
// `data/` directory directly into the compiled Go binary. This ensures that users
// of the library don't have to manage external data files or perform runtime downloads
// to calculate stats.
package assets

import "embed"

//go:embed data/*.json
var DataFS embed.FS
