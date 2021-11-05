package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	cmd  string
	opt  string
	date string
)

var workouts = []string{"pullups", "pushups", "chinups", "dips", "squats", "plank", "ttb"}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func printHelp() {
	fmt.Println(`Usage: myworkout COMMAND OPTION
COMMANDS:
	help				Print help message
	list all|EXERCISE [DATE]	List exercise(s) for the day (default today)
	EXERCISE +/-COUNT		Add the count to the total sum for today
EXERCISES:
	pullups
	pushups
	chinups
	dips
	squats
	plank
	ttb (toes to bar)
	`)
}

func findEx(ex string) bool {
	for _, i := range workouts {
		if ex == i {
			return true
		}
	}
	return false
}

func getInput() {
	fail := true
	if len(os.Args) >= 3 && len(os.Args) < 5 {
		cmd = strings.ToLower(os.Args[1])
		opt = strings.ToLower(os.Args[2])
		date = "today"
		if len(os.Args) == 4 {
			date = os.Args[3]
		}

		if cmd == "list" {
			if opt == "all" || findEx(opt) {
				fail = false
			} else {
				fmt.Println("Unknown argument...")
			}
		} else if findEx(cmd) {
			if string(opt[0]) == "+" || string(opt[0]) == "-" {
				fail = false
			} else {
				fmt.Println("You have to use +/-COUNT")
			}
		} else {
			fmt.Println("Unknown command...")
		}
	} else if len(os.Args) > 4 {
		fmt.Println("Too many arguments...")
	} else if len(os.Args) < 3 {
		fmt.Println("You have to specify at least two arguments...")
	} else if os.Args[1] == "help" {
	} else {
		fmt.Println("Unknown argument...")
	}

	if fail {
		printHelp()
		os.Exit(0)
	}
}

func main() {
	getInput()

	conn, err := net.Dial("tcp", "localhost:31564")
	checkError(err)
	defer conn.Close()

	fmt.Fprintf(conn, "%s %s %s\n", cmd, opt, date)

	if cmd == "list" {
		msg := bufio.NewScanner(conn)
		for msg.Scan() {
			if msg.Text() == "<end>" {
				break
			}
			fmt.Println(msg.Text())
		}
	}
}
