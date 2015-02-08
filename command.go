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
	"os"
)

//CommandLine is the replacement for all variables see commander for implementation

// CommandLine is the default commander.
// The top-level functions such as On, Usage, Parse and so on are wrappers for the
// methods of CommandLine.
var CommandLine = NewCommander(os.Args[0], flag.CommandLine)

// Cmd represents a sub command, allowing to define subcommand
// flags and runnable to run once arguments match the subcommand
// requirements.
type Cmd interface {
	Flags(*flag.FlagSet) *flag.FlagSet
	Run(args []string)
}

type cmdCont struct {
	name          string
	desc          string
	command       Cmd
	requiredFlags []string
}

// Registers a Cmd for the provided sub-command name. E.g. name is the
// `status` in `git status`.
func On(name, description string, command Cmd, requiredFlags []string) {
	CommandLine.On(name, description, command, requiredFlags)
}

// Prints the usage.
func Usage() {
	CommandLine.Usage()
}

// Parses the flags and leftover arguments to match them with a
// sub-command. Evaluate all of the global flags and register
// sub-command handlers before calling it. Sub-command handler's
// `Run` will be called if there is a match.
// A usage with flag defaults will be printed if provided arguments
// don't match the configuration.
// Global flags are accessible once Parse executes.
func Parse() {
	flag.Parse() // the recursive definition of Commander requires that the command flag is parsed before calling Parse()
	CommandLine.Parse()
}

// Runs the subcommand's runnable. If there is no subcommand
// registered, it silently returns.
func Run() { CommandLine.Run(nil) }

// Parses flags and run's matching subcommand's runnable.
func ParseAndRun() {
	Parse()
	Run()
}
