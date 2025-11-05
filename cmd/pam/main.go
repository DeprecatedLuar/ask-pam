package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/lipgloss"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/eduardofuncao/pam/internal/config"
	"github.com/eduardofuncao/pam/internal/db"
	"github.com/eduardofuncao/pam/internal/editor"
	"github.com/eduardofuncao/pam/internal/spinner"
	"github.com/eduardofuncao/pam/internal/table"
)

func main() {
	cfg, err := config.LoadConfig(config.CfgFile)
	if err != nil {
		log.Fatal("Could not load config file", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("pam create <name> <db-type> <connection-string>")
		fmt.Println("pam switch <db-name>")
		fmt.Println("pam add <query-name> <query>")
		fmt.Println("pam query <query-name>")
		fmt.Println("pam get <db-type> <connection-string> <sql-query>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {

	case "init":
		if len(os.Args) < 5 {
			log.Fatal("Usage: pam create <name> <db-type> <connection-string> <user> <password>")
		}

		conn, err := db.CreateConnection(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			log.Fatalf("Could not create connection interface: %s/%s, %s", os.Args[3], os.Args[2], err)
		}

		err = conn.Open()
		if err != nil {
			log.Fatalf("Could not establish connection to: %s/%s: %s",
				conn.GetDbType(), conn.GetName(), err)
		}
		defer conn.Close()

		err = conn.Ping()
		if err != nil {
			log.Fatalf("Could not communicate with the database: %s/%s, %s", os.Args[3], os.Args[2], err)
		}

		cfg.CurrentConnection = conn.GetName()
		cfg.Connections[cfg.CurrentConnection] = config.ToConnectionYAML(conn)
		cfg.Save()

	case "switch", "use":
		if len(os.Args) < 3 {
			log.Fatal("Usage: pam switch/use <db-name>")
		}

		_, ok := cfg.Connections[os.Args[2]]
		if !ok {
			log.Fatalf("Connection %s does not exist", os.Args[2])
		}
		cfg.CurrentConnection = os.Args[2]

		err := cfg.Save()
		if err != nil {
			log.Fatal("Could not save configuration file")
		}
		fmt.Printf("connected to: %s/%s\n", cfg.Connections[cfg.CurrentConnection].DBType, cfg.CurrentConnection)

	case "add", "save":
		if len(os.Args) < 4 {
			log.Fatal("Usage: pam add <query-name> <query>")
		}

		_, ok := cfg.Connections[cfg.CurrentConnection]
		if !ok {
			cfg.Connections[cfg.CurrentConnection] = config.ConnectionYAML{}
		}
		queries := cfg.Connections[cfg.CurrentConnection].Queries

		queries[os.Args[2]] = db.Query{
			Name: os.Args[2],
			SQL:  os.Args[3],
			Id:   db.GetNextQueryId(queries),
		}
		err := cfg.Save()
		if err != nil {
			log.Fatal("Could not save configuration file")
		}

	case "remove", "delete":
		if len(os.Args) < 3 {
			log.Fatal("Usage: pam remove <query-name>")
		}
		
		conn := cfg.Connections[cfg.CurrentConnection]
		queries := conn.Queries

		query, exists := db.FindQueryWithSelector(queries, os.Args[2])
		if exists{
			delete(conn.Queries, query.Name)
		} else {
			log.Fatalf("Query %s could not be found", os.Args[2])
		}
		err := cfg.Save()
		if err != nil {
			log.Fatal("Could not save configuration file")
		}
		

	case "query", "run":
		if len(os.Args) < 3 {
			log.Fatal("Usage:pam query/run <query-name>")
		}

		editMode := false
		if len(os.Args) > 3 {
			if os.Args[3] == "--edit" || os.Args[3] == "-e" {
				editMode = true
			}
		}

		currConn := config.FromConnectionYaml(cfg.Connections[cfg.CurrentConnection])

		queries := currConn.GetQueries()
		selector := os.Args[2]
		query, found := db.FindQueryWithSelector(queries, selector)
		if !found {
			log.Fatalf("Could not find query with name/id: %v", selector)
		}

		editedQuery, submitted, err := editor.EditQuery(query, editMode)
		if submitted {
			cfg.Connections[cfg.CurrentConnection].Queries[query.Name] = editedQuery
			cfg.Save()
		}

		err = currConn.Open()
		if err != nil {
			log.Fatalf("Could not open the connection to %s/%s: %s", currConn.GetDbType(), currConn.GetName(), err)
		}

		start := time.Now()
		done := make(chan struct{})
		go spinner.Wait(done)

		rows, err := currConn.Query(query.Name)
		if err != nil {
			log.Fatal("Could not complete query: ", err)
		}
		sqlRows, ok := rows.(*sql.Rows)
		if !ok {
			log.Fatal("Query did not return *sql.Rows")
		}
		columns, data, err := db.FormatTableData(sqlRows)

		done <- struct{}{}
		elapsed := time.Since(start)

		if err := table.Render(columns, data, elapsed); err != nil {
			log.Fatalf("Error rendering table: %v", err)
		}

	case "list":
		if len(os.Args) < 3 {
			log.Fatal("Usage:pam list [queries/connections]")
		}

		var objectType string
		if len(os.Args) < 3 {
			objectType = ""
		} else {
			objectType = os.Args[2]
		}

		switch objectType {
		case "connections":
			for name, connection := range cfg.Connections {
				fmt.Printf("◆ %s (%s)\n", name, connection.ConnString)
			}

		case "", "queries":
			titleStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205"))

			for _, query := range cfg.Connections[cfg.CurrentConnection].Queries {
				formatedItem := fmt.Sprintf("\n◆ %d/%s", query.Id, query.Name)
				fmt.Println(titleStyle.Render(formatedItem))
				fmt.Println(editor.HighlightSQL(editor.FormatSQLWithLineBreaks(query.SQL)))
			}
		}

	case "edit":
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		cmd := exec.Command(editor, config.CfgFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to open editor: %v", err)
		}

	case "status":
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("171")).
			Bold(true)
		currConn := cfg.Connections[cfg.CurrentConnection]
		fmt.Println(style.Render("✓ Now using:"), fmt.Sprintf("%s/%s", currConn.DBType, currConn.Name))

	case "history":
		fmt.Println("To be implemented in future releases...")

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
