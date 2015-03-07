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
// The top-level functions On and  Run are wrapper around this var
var CommandLine = New()

// Cmd represents a sub command
type Cmd interface {
	Run(args []string)
}

//Flagger is the interface that defines the Flag method that allow to configure flags
type Flagger interface {
	Flags(*flag.FlagSet)
}

//Completer is the interface that defines the Compgens method that allow to configure a Terminator
type Completer interface {
	Compgens(*compgen.Terminator)
}

// Commander can register sub commands and:
//
// - Configure flag.FlagSet Usage function
//
// - Configure compgen.Terminator to complete command line with subcommands
//
// - Run the matching subcommand
//
type Commander interface {
	Cmd
	Flagger         //configure flags
	Completer       // configure a terminator
	compgen.Argsgen //ability to complete var args
	On(name, syntax, description string, command Cmd)
	Path(qname string)
}

//New creates a new Commander.
func New() Commander { return &commander{cmds: make(map[string]*cmdCont)} }

// Registers a Cmd for the provided sub-command name
//
// name is the command name: like 'status' in 'git status'
//
// syntax is the usual command syntax description like :
//
//      git [--version] [--help] [-C <path>] [-c <name>=<value>] <command> [<args>]
//
// description is a short line describing the subcommand.
//
func On(name, syntax, description string, command Cmd) {
	CommandLine.On(name, syntax, description, command)
}

// Runs the default commander.
func Run() {
	Launch(CommandLine, os.Args[0], os.Args)
}

//Launch a standalone Cmd.
//
// name: the command qualified name to be displayed in the Usage ("git status")
//
// args: the command args (args[0] shall contain the local command name 'status' for instance)
//
func Launch(cmd Cmd, name string, args []string) {
	matchingFlags, term := prepare(cmd, name)
	term.Terminate()
	matchingFlags.Parse(args[1:])
	cmd.Run(matchingFlags.Args())
}
