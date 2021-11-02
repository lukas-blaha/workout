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

var workouts = []string{"pullups", "pushups", "dips", "squats"}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "workout"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func findEx(ex string) bool {
	for _, i := range workouts {
		if ex == i {
			return true
		}
	}
	return false
}

func queryDb(conn net.Conn, db *sql.DB, ent *entry, prnt bool, date string, exercise string) {
	var condition string
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

		if prnt {
			fmt.Fprintf(conn, "%s: %d\n", ent.name, ent.count)
		}
	}
}

type entry struct {
	id    int
	date  string
	name  string
	count int
}

func manageRequest(conn net.Conn, args []string) {
	dbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbconn)
	checkError(err)
	defer db.Close()

	if len(args) >= 2 {
		opt := strings.ToLower(args[0])
		opt = strings.ReplaceAll(opt, "-", "")

		if opt == "list" {
			if len(args) >= 2 {
				exercise := strings.ToLower(args[1])
				exercise = strings.ReplaceAll(exercise, "-", "")

				if findEx(exercise) || exercise == "all" {
					date := "today"
					if len(args) > 2 {
						date = args[2]
					}
					var ent entry
					queryDb(conn, db, &ent, true, date, exercise)
				}
			}
		} else if findEx(opt) {
			if len(args) >= 2 {
				count := args[1]
				if string(count[0]) == "-" || string(count[0]) == "+" {
					var ex entry
					queryDb(conn, db, &ex, false, "today", opt)
					if ex.id == 0 {
						inp := `insert into "exercises"("date", "name", "count") values($1, $2, $3)`
						_, err := db.Exec(inp, "today", opt, count)
						checkError(err)
					} else {
						n, err := strconv.Atoi(count[1:])
						checkError(err)
						inp := `update "exercises" set "count" = $1 where "id" = $2`
						if string(count[0]) == "+" {
							_, err := db.Exec(inp, (ex.count + n), ex.id)
							checkError(err)
						} else {
							_, err := db.Exec(inp, (ex.count - n), ex.id)
							checkError(err)
						}
					}

				} else {
					fmt.Fprintf(conn, "You have to use +/- new value\n%s\n",
						"Ex.: myworkout pushups +20")
				}
			}
		} else {
			fmt.Fprintln(conn, "incorrect option")
		}
	} else {
		fmt.Fprintln(conn, "You have to specify at least one option")
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
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		manageRequest(conn, strings.Split(scanner.Text(), " "))
		break
	}
}
