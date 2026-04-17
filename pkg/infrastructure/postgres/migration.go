package postgres

import _ "embed"

//go:embed migrations/20240115100000_create_products_table.up.sql
var MigrationSQL string

//go:embed migrations/20240116100000_create_users_table.up.sql
var UserMigrationSQL string
