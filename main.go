package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jedib0t/go-pretty/table"
	"github.com/mattn/go-sqlite3"
	"github.com/packethost/packngo"
	"github.com/patrickdevivo/emsql/pkg/emqlite"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	flag.Parse()

	authToken := os.Getenv("PACKET_AUTH_TOKEN")

	if authToken == "" {
		handleErr(errors.New("please supply an auth token via the PACKET_AUTH_TOKEN env var"))
	}

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	client := packngo.NewClientWithAuth("", authToken, httpClient)

	sql.Register("emsql", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("devices", emqlite.NewDevicesModule(client))
		},
	})
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("emsql", ":memory:")
	handleErr(err)
	defer db.Close()

	query := flag.Arg(0)

	if query == "" {
		handleErr(errors.New("must supply a query"))
	}

	rows, err := db.Query(query)
	handleErr(err)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	columns, err := rows.Columns()
	handleErr(err)

	h := make(table.Row, len(columns))
	for c, col := range columns {
		h[c] = col
	}
	t.AppendHeader(h)

	pointers := make([]interface{}, len(columns))
	container := make([]sql.NullString, len(columns))

	for i := range pointers {
		pointers[i] = &container[i]
	}
	for rows.Next() {
		err := rows.Scan(pointers...)
		handleErr(err)

		r := make(table.Row, len(columns))
		for i, c := range container {
			if c.Valid {
				r[i] = c.String
			} else {
				r[i] = ""
			}
		}

		t.AppendRow(r)
	}

	width, _, err := terminal.GetSize(0)
	handleErr(err)

	t.SetAllowedRowLength(width)
	t.Render()
}
