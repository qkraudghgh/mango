package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qkraudghgh/mango/manager"
)

const usage = `
NAME:
   mango - very simple todo app in your terminal

USAGE:
   mango command [arguments...]

VERSION:
   1.0.0

AUTHOR:
  myoungho.pak - <qkraudghgh@gmail.com>

COMMANDS:
   list - Show your todos
   add "your todo" - Add your todo
   done [number] - Check done your todo
   undone [number] - Uncheck done your todo
   delete [number] - Delete your todo
   clear - Clear your todo
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
	mango.AddCommand(clearCommand)

	args := flag.Args()

	if err := mango.Execute(args); err != nil {
		fmt.Println(err)
		fmt.Print(mango.Usage())
	}
}
