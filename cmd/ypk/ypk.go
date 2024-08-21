// Yew package manager
package ypk

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/petersalex27/ypk/config"
)

// global variables
var (
	// install command
	installCmd        = flag.NewFlagSet("install", flag.ExitOnError)
	installCmdName    = installCmd.String("name", "", "Package name")
	installCmdVersion = installCmd.String("version", "", "Package version")
	installCmdSource  = installCmd.String("source", "", "Package source")
	installCmdForce   = installCmd.Bool("force", false, "Force install")
	installCmdNoDeps  = installCmd.Bool("no-deps", false, "Do not install dependencies")
	// remove command
	removeCmd        = flag.NewFlagSet("remove", flag.ExitOnError)
	removeCmdName    = removeCmd.String("name", "", "Package name")
	removeCmdVersion = removeCmd.String("version", "", "Package version")
	removeCmdForce   = removeCmd.Bool("force", false, "Force remove")
	// search command
	searchCmd        = flag.NewFlagSet("search", flag.ExitOnError)
	searchCmdName    = searchCmd.String("name", "", "Package name")
	searchCmdVersion = searchCmd.String("version", "", "Package version")
	// update command
	updateCmd = flag.NewFlagSet("update", flag.ExitOnError)
	// upgrade command
	upgradeCmd        = flag.NewFlagSet("upgrade", flag.ExitOnError)
	upgradeCmdName    = upgradeCmd.String("name", "", "Package name")
	upgradeCmdVersion = upgradeCmd.String("version", "", "Package version")
	upgradeCmdSource  = upgradeCmd.String("source", "", "Package source")
	upgradeCmdForce   = upgradeCmd.Bool("force", false, "Force upgrade")
	upgradeCmdNoDeps  = upgradeCmd.Bool("no-deps", false, "Do not upgrade dependencies")
)

type formatter struct {
	padding [2]int // left and right padding at 0 and 1 respectively
}

// printCmds function to print commands
func (f formatter) printCmds(cmds [][2]string) {
	padLeft := strings.Repeat(" ", f.padding[0])
	// max amount of padding required for an empty first column
	padRight := strings.Repeat(" ", f.padding[1])
	for _, cmd := range cmds {
		paddingLength := f.padding[1] - len(cmd[0])
		command := padLeft + cmd[0] + padRight[:paddingLength]
		commandDesc := cmd[1]
		fmt.Printf("%s %s\n", command, commandDesc)
	}
}

func makeFormatter(commands [][2]string) formatter {
	formatter := formatter{padding: [2]int{2, math.MinInt}}
	for _, cmd := range commands {
		if x := len(cmd[0]); x > formatter.padding[1] {
			formatter.padding[1] = x
		}
	}
	return formatter
}

type commandNoDesc struct {
	name string `json:"name"`
}

type command struct {
	commandNoDesc
	desc string `json:"desc"`
	alts []commandNoDesc `json:"alts"`
}

// initial capacity for commands slice
const initCmdCap int = 10

// set of commands, mutex for async
type commandSet struct {
	cmds []*command
	*sync.Mutex
}

func makeCommandSet() commandSet {
	return commandSet{
		cmds:  make([]*command, 0, initCmdCap),
		Mutex: &sync.Mutex{},
	}
}

// global command set
var commands commandSet = makeCommandSet()

func makeCommand() (cmd command, commit func(cmd command)) {
	commit = func(cmd command) {
		commands.Lock()
		defer commands.Unlock()

		input := &cmd
		commands.cmds = append(commands.cmds, input)
	}
	return
}

func readCommands(conf config.Config) {
	f, err := os.Open(conf.CommandJsonPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	jsonParser := json.NewDecoder(f)

	for jsonParser.More() {
		cmd, commit := makeCommand()

		err = jsonParser.Decode(&cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		commit(cmd)
	}
}

// init function to set custom usage message
func init() {
	conf := config.GetConfig()

	// read commands from json
	readCommands()

	// commands := [][2]string{
	// 	{"install", "Install a package"},
	// 	{"remove", "Remove a package"},
	// 	{"search", "Search for a package"},
	// 	{"update", "Update package list"},
	// 	{"upgrade", "Upgrade packages"},
	// }
	fmttr := makeFormatter(commands)

	flag.Usage = func() {
		println("Usage: ypk [command] [options]")
		println("Commands:")
		println(" install    Install a package")
		println(" list       List installed packages")
		println(" remove     Remove a package")
		println(" search     Search for a package")
		println(" update     Update package list")
		println("Options:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
}
