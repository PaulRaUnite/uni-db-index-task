package assets

import "github.com/gobuffalo/packr/v2"

//go:generate packr2 clean
//go:generate packr2

// SchemaMigrations - defines postgres db migration for table schema
var SchemaMigrations = packr.New("migrations", "./migrations")
