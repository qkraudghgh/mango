package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qkraudghgh/mango/manager"
)

const usage = `Usage:
	mango list
		Show all tasks
	mango add
		Add todo task
	mango done [number]
		Check done todo task
	mango undone [number]
		Uncheck done todo task
	mango delete [number]
`

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println(usage)
		os.Exit(-1)
	}

	mango := manager.New()

	mango.AddCommand(addCommand)
	mango.AddCommand(listCommand)
	mango.AddCommand(deleteCommand)
	mango.AddCommand(doneCommand)
	mango.AddCommand(unDoneCommand)

	args := flag.Args()

	if err := mango.Execute(args); err != nil {
		fmt.Println(err)
		fmt.Print(mango.Usage())
	}
}
