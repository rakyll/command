// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package command allows you to define subcommands
// for your command line interfaces. It extends the flag package
// to provide flag support for subcommands.
package command

import (
	"flag"
	"fmt"
	"os"
)

// A map of all of the registered sub-commands.
var cmds map[string]*cmdCont = make(map[string]*cmdCont)

// Matching subcommand.
var matchingCmd *cmdCont

// Arguments to call subcommand's runnable.
var args []string

// Flag to determine whether help is
// asked for subcommand or not
var flagHelp *bool

// Cmd represents a sub command, allowing to define subcommand
// flags and runnable to run once arguments match the subcommand
// requirements.
type Cmd interface {
	Flags(*flag.FlagSet) *flag.FlagSet
	Run(args []string)
}

type cmdCont struct {
	name         string
	desc         string
	command      Cmd
	requiredArgs []string
}

// Registers a Cmd for the provided sub-command name. E.g. name is the
// `status` in `git status`.
func On(name, description string, command Cmd, reqArgs []string) {
	cmds[name] = &cmdCont{
		name:         name,
		desc:         description,
		command:      command,
		requiredArgs: reqArgs,
	}
}

// Prints the usage.
func Usage() {
	program := os.Args[0]
	if len(cmds) == 0 {
		// no subcommands
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", program)
		flag.PrintDefaults()
		return
	}

	fmt.Fprintf(os.Stderr, "Usage: %s <command>\n\n", program)
	fmt.Fprintf(os.Stderr, "where <command> is one of:\n")
	for name, cont := range cmds {
		fmt.Fprintf(os.Stderr, "  %-15s %s\n", name, cont.desc)
	}

	if numOfGlobalFlags() > 0 {
		fmt.Fprintf(os.Stderr, "\navailable flags:\n")
		flag.PrintDefaults()
	}
	fmt.Fprintf(os.Stderr, "\n%s <command> -h for subcommand help\n", program)
}

func subcommandUsage(cont *cmdCont) {
	fmt.Fprintf(os.Stderr, "Usage of %s %s:\n", os.Args[0], cont.name)
	// should only output sub command flags, ignore h flag.
	fs := matchingCmd.command.Flags(flag.NewFlagSet(cont.name, flag.ContinueOnError))
	fs.PrintDefaults()
	if len(cont.requiredArgs) > 0 {
		fmt.Fprintf(os.Stderr, "\nArguments:\n\n")
		for _, a := range cont.requiredArgs {
			fmt.Fprintf(os.Stderr, "  %s\n", a)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
}

// Parses the flags and leftover arguments to match them with a
// sub-command. Evaluate all of the global flags and register
// sub-command handlers before calling it. Sub-command handler's
// `Run` will be called if there is a match.
// A usage with flag defaults will be printed if provided arguments
// don't match the configuration.
// Global flags are accessible once Parse executes.
func Parse() {
	flag.Parse()
	// if there are no subcommands registered,
	// return immediately
	if len(cmds) < 1 {
		return
	}

	flag.Usage = Usage
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	name := flag.Arg(0)
	if cont, ok := cmds[name]; ok {
		fs := cont.command.Flags(flag.NewFlagSet(name, flag.ExitOnError))
		flagHelp = fs.Bool("h", false, "")
		fs.Parse(flag.Args()[1:])
		args = fs.Args()
		if len(args) < len(cont.requiredArgs) {
			*flagHelp = true
		}
		matchingCmd = cont
	} else {
		flag.Usage()
		os.Exit(1)
	}
}

// Runs the subcommand's runnable. If there is no subcommand
// registered, it silently returns.
func Run() {
	if matchingCmd != nil {
		if *flagHelp {
			subcommandUsage(matchingCmd)
			return
		}
		matchingCmd.command.Run(args)
	}
}

// Parses flags and run's matching subcommand's runnable.
func ParseAndRun() {
	Parse()
	Run()
}

// Returns the total number of globally registered flags.
func numOfGlobalFlags() (count int) {
	flag.VisitAll(func(flag *flag.Flag) {
		count++
	})
	return
}
