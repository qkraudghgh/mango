package manager

import (
	"bytes"
	"fmt"
)

// MangoBucket is global bucket name
const MangoBucket = "todos"

// Manager structure
type Manager struct {
	commands map[string]Command
}

// Command structure
type Command struct {
	Name  string
	Usage string
	Run   func([]string) error
}

// New function make new manager
func New() *Manager {
	return &Manager{
		commands: make(map[string]Command),
	}
}

// Usage creates usage message string of all available commands.
func (m *Manager) Usage() string {
	buf := bytes.NewBufferString("\n")

	for _, c := range m.commands {
		fmt.Fprintln(buf, c.Usage)
	}

	return buf.String()
}

// AddCommand method is add New Command
func (m *Manager) AddCommand(cmd Command) {
	m.commands[cmd.Name] = cmd
}

// Execute parses the command line arguments and
// runs the 'Run' function of command with that parsed arguments.
func (m *Manager) Execute(args []string) error {
	var cmdName string
	var cmdArgs []string

	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	cmdName = args[0]

	cmd, ok := m.commands[cmdName]
	if !ok {
		return fmt.Errorf("%s is not defined", cmdName)
	}

	if err := cmd.Run(cmdArgs); err != nil {
		return err
	}

	return nil
}
