package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/qkraudghgh/mango/printer"
)

// Todo task structure
type Todo struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsCheck   int       `json:"is_check"`
}

// PrintTodos is draw todo tasks to terminal
func PrintTodos(todos []Todo) {
	if len(todos) > 0 {
		bold := color.New(color.Bold).SprintFunc()

		// TODO: sort unDone, Done
		for _, todo := range todos {
			fmt.Printf("\n %d. %s %s\n", todo.ID, checkStamp(todo.IsCheck), bold(todo.Content))
		}

		fmt.Println()
	} else {
		fmt.Println("There is not Todo")
	}
}

// return different check symbol depend on OS
func checkStamp(isCheck int) string {
	var symbol string

	green := color.New(color.FgGreen).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()

	// check OS
	if runtime.GOOS == "windows" {
		if isCheck == 1 {
			symbol = green(printer.DoneSignW)
		} else {
			symbol = red(printer.UndoneSignW)
		}
	} else {
		if isCheck == 1 {
			symbol = green(printer.DoneSign)
		} else {
			symbol = red(printer.UndoneSign)
		}
	}

	return symbol
}
