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
var cmds map[string]Cmd = make(map[string]Cmd)

// Cmd represents a sub command, allowing to define subcommand
// flags and runnable to run once arguments match the subcommand
// requirements.
type Cmd interface {
	Flags(*flag.FlagSet) *flag.FlagSet
	Run(args []string)
}

// Registers a Cmd for the provided sub-command name. E.g. name is the
// `status` in `git status`.
func On(name string, command Cmd) {
	cmds[name] = command
}

// Prints the usage.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	for name, cmd := range cmds {
		fmt.Fprintf(os.Stderr, "\n  %s\n", name)
		// should only output sub command flags
		fs := cmd.Flags(flag.NewFlagSet(name, flag.ContinueOnError))
		fs.PrintDefaults()
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
	if cmd, ok := cmds[name]; ok {
		fs := cmd.Flags(flag.NewFlagSet(name, flag.ExitOnError))
		args := flag.Args()[1:]
		fs.Parse(args)
		cmd.Run(args)
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
