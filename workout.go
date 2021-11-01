package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

func queryDb(db *sql.DB, ent *entry, prnt bool, date string, exercise string) {
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
			fmt.Printf("%s: %d\n", ent.name, ent.count)
		}
	}
}

type entry struct {
	id    int
	date  string
	name  string
	count int
}

func main() {
	dbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbconn)
	checkError(err)
	defer db.Close()

	if len(os.Args) > 2 {
		opt := strings.ToLower(os.Args[1])
		opt = strings.ReplaceAll(opt, "-", "")

		if opt == "list" {
			if len(os.Args) >= 3 {
				exercise := strings.ToLower(os.Args[2])
				exercise = strings.ReplaceAll(exercise, "-", "")

				if findEx(exercise) || exercise == "all" {
					date := "today"
					if len(os.Args) > 3 {
						date = os.Args[3]
					}
					var ent entry
					queryDb(db, &ent, true, date, exercise)
				}
			}
		} else if findEx(opt) {
			if len(os.Args) > 2 {
				count := os.Args[2]
				if string(count[0]) == "-" || string(count[0]) == "+" {
					var ex entry
					queryDb(db, &ex, false, "today", opt)
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
					fmt.Printf("You have to use +/- new value\n%s\n",
						"Ex.: myworkout pushups +20")
				}
			}
		} else {
			fmt.Println("incorrect option")
		}
	} else {
		fmt.Println("You have to specify at least one option")
	}
}
