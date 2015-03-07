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

	"github.com/ericaro/compgen"
)

//CommandLine is the replacement for all variables see commander for implementation

// CommandLine is the default Commander.
// The top-level functions such as On, Usage, Parse and so on are wrappers for the
// methods of CommandLine.
var CommandLine = NewCommander()

// Cmd represents a sub command, the simplest subcommands on have to implement this interface
type Cmd interface {
	Run(args []string)
}

//Flagger is the interface that defines the Flag methods that allow to configure flags
type Flagger interface {
	Flags(*flag.FlagSet)
}

//Completer is the interface that defines the Compgens methods that allow to configure a Terminator
type Completer interface {
	Compgens(*compgen.Terminator)
}

// simple command container
type cmdCont struct {
	name, syntax, description string
	command                   Cmd
}

// Registers a Cmd for the provided sub-command name. E.g. name is the
// `status` in `git status`.
func On(name, syntax, description string, command Cmd) {
	CommandLine.On(name, syntax, description, command)
}

// Runs the subcommand's runnable. If there is no subcommand
// registered, it silently returns.
func Run() {
	Launch(CommandLine, os.Args[0], os.Args)
}
