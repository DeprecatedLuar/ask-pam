// Package engines blank-imports every database engine package so their
// init() functions run and self-register with internal/db's registry.
// Import this package once at startup (from main) before calling
// db.CreateConnection.
package engines

import (
	_ "github.com/eduardofuncao/squix/internal/db/clickhouse"
	_ "github.com/eduardofuncao/squix/internal/db/duckdb"
	_ "github.com/eduardofuncao/squix/internal/db/firebird"
	_ "github.com/eduardofuncao/squix/internal/db/mysql"
	_ "github.com/eduardofuncao/squix/internal/db/oracle"
	_ "github.com/eduardofuncao/squix/internal/db/postgres"
	_ "github.com/eduardofuncao/squix/internal/db/snowflake"
	_ "github.com/eduardofuncao/squix/internal/db/sqlite"
	_ "github.com/eduardofuncao/squix/internal/db/sqlserver"
)
