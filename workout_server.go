package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "workout"
)

type entry struct {
	id    int
	date  string
	name  string
	count int
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func queryDb(conn net.Conn, db *sql.DB, ent *entry, prnt bool, date string, exercise string) {
	var condition string
	var entries []entry
	if exercise == "all" {
		condition = fmt.Sprintf(" where date = '%s'", date)
	} else {
		condition = fmt.Sprintf(" where date = '%s' and name = '%s'", date, exercise)
	}
	qr := fmt.Sprintf("select * from exercises%s", condition)

	rows, err := db.Query(qr)
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&ent.id, &ent.date, &ent.name, &ent.count)
		checkError(err)

		entries = append(entries, *ent)
	}

	if prnt {
		for _, i := range entries {
			fmt.Fprintf(conn, "%s: %d\n", i.name, i.count)
		}
		if len(entries) == 0 {
			fmt.Fprintln(conn, "There's nothing to show")
		}
		fmt.Fprintln(conn, "<end>")
	}
}

func manageRequest(conn net.Conn, args []string) {
	dbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbconn)
	checkError(err)
	defer db.Close()

	cmd := args[0]
	opt := args[1]
	date := args[2]

	if cmd == "list" {
		var ent entry
		queryDb(conn, db, &ent, true, date, opt)
	} else {
		var ex entry
		n, err := strconv.Atoi(opt[1:])
		checkError(err)
		queryDb(conn, db, &ex, false, "today", opt)
		if ex.id == 0 {
			inp := `insert into "exercises"("date", "name", "count") values($1, $2, $3)`
			_, err := db.Exec(inp, "today", cmd, n)
			checkError(err)
		} else {
			inp := `update "exercises" set "count" = $1 where "id" = $2`
			if string(opt[0]) == "+" {
				_, err := db.Exec(inp, (ex.count + n), ex.id)
				checkError(err)
			} else {
				_, err := db.Exec(inp, (ex.count - n), ex.id)
				checkError(err)
			}
		}
	}
}

func main() {
	li, err := net.Listen("tcp", ":31564")
	checkError(err)
	defer li.Close()

	for {
		conn, err := li.Accept()
		checkError(err)

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	msg, err := bufio.NewReader(conn).ReadString('\n')
	checkError(err)
	msg = strings.ReplaceAll(msg, "\n", "")
	manageRequest(conn, strings.Split(msg, " "))
}
