package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/eduardofuncao/pam/internal/config"
	"github.com/eduardofuncao/pam/internal/db"
	"github.com/eduardofuncao/pam/internal/spinner"
	"github.com/eduardofuncao/pam/internal/table"
)

func Explore(cfg *config.Config) {
	if len(os.Args) < 3 {
		fmt.Println("No table specified. Available tables:")
		ListTables(cfg)
		return
	}

	tableName := os.Args[2]
	limit := 1000

	if len(os.Args) > 3 {
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--limit" || os.Args[i] == "-l" {
				if i+1 < len(os.Args) {
					parsedLimit, err := strconv.Atoi(os.Args[i+1])
					if err != nil {
						log.Fatalf("Invalid limit value: %s", os.Args[i+1])
					}
					limit = parsedLimit
					break
				}
			}
		}
	}

	currConn := config.FromConnectionYaml(cfg.Connections[cfg.CurrentConnection])

	if err := currConn.Open(); err != nil {
		log.Fatalf("Could not open the connection to %s/%s: %s", currConn.GetDbType(), currConn.GetName(), err)
	}

	querySQL := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, limit)

	start := time.Now()
	done := make(chan struct{})
	go spinner.Wait(done)

	rows, err := currConn.QueryDirect(querySQL)
	if err != nil {
		done <- struct{}{}
		log.Fatalf("Could not query table '%s': %v", tableName, err)
	}

	sqlRows, ok := rows.(*sql.Rows)
	if !ok {
		done <- struct{}{}
		log.Fatal("Query did not return *sql.Rows")
	}

	tableData, err := db.BuildTableData(sqlRows, querySQL, currConn)
	if err != nil {
		done <- struct{}{}
		log.Fatalf("Error building table data: %v", err)
	}

	done <- struct{}{}
	elapsed := time.Since(start)

	if err := table.Render(tableData, elapsed); err != nil {
		log.Fatalf("Error rendering table: %v", err)
	}
}
