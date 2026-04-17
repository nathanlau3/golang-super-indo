package postgres

import _ "embed"

//go:embed migrations/20240115100000_create_products_table.up.sql
var MigrationSQL string
