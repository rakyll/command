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

package command

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/ericaro/compgen"
)

//Commander can register sub commands and:
//
// - print their Usage
//
// - Parse flagset to find out subcommands, and required flags
//
// - Run the matching subcommand
//
// Command implements Cmd.Run([]string) so you can do things like:
//
//      type myCmd struct {
//     	    command.Commander
//      }
//
//		func (c *myCmd) Flags(fs *flag.FlagSet) {
//     	   	c.Commander = command.NewCommander("git", fs)
//     	    c.On("status","Show the working tree status",&GitStatuCmd{}, nil)
//      	...
//		}
//
// Works out of the box.
type Commander interface {
	Cmd             // run a command
	Flagger         //configure flags
	Completer       // configure a terminator
	compgen.Argsgen //ability to complete var args
	// Registers a Cmd for the provided sub-command name. E.g. name is the
	// `status` in `git status`.
	On(name, syntax, description string, command Cmd)
	Path(qname string)
}

type commander struct {
	name string
	cmds map[string]*cmdCont // A map of all of the registered sub-commands.
}

//NewCommander creates a new Commander. The 'name' is the the subcommand name,
// and the fs is the Cmd current flagset.
//
// A NewCommander is better created in the Flag(fs *flag.FlagSet) method.
func NewCommander() Commander { return &commander{cmds: make(map[string]*cmdCont)} }

// implementing the three interface to configure path, flags, and compgens

func (c *commander) Path(qname string)                 { c.name = qname }
func (c *commander) Flags(fs *flag.FlagSet)            { fs.Usage = func() { c.Usage(fs) } }
func (c *commander) Compgens(term *compgen.Terminator) { term.Argsgen(c) }

// Registers a Cmd for the provided sub-command name. E.g. name is the
// `status` in `git status`.
func (c *commander) On(name, syntax, description string, command Cmd) {
	c.cmds[name] = &cmdCont{
		name:        name,
		syntax:      syntax,
		description: description,
		command:     command,
	}
}

// Run the flags and leftover arguments to match them with a
// sub-command. Evaluate all of the global flags and register
// sub-command handlers before calling it. Sub-command handler's
// `Run` will be called if there is a match.
// A usage with flag defaults will be printed if provided arguments
// don't match the configuration.
// Global flags are accessible once Parse executes.
func (c *commander) Run(args []string) {

	cont := c.getContainer(args)
	if cont == nil {
		c.Usage(nil)
		os.Exit(1)
	}

	Launch(cont.command, c.name+" "+cont.name, args)

}
func prepare(cmd Cmd, name string) (*flag.FlagSet, *compgen.Terminator) {

	matchingFlags := flag.NewFlagSet(name, flag.ExitOnError)
	term := compgen.NewTerminator(matchingFlags)

	if cmdr, ok := cmd.(Commander); ok {
		cmdr.Path(name)
	}
	if flagger, ok := cmd.(Flagger); ok {
		flagger.Flags(matchingFlags)
	}
	if cg, ok := cmd.(Completer); ok {
		cg.Compgens(term)
	}
	return matchingFlags, term
}

func Launch(cmd Cmd, name string, args []string) {

	matchingFlags, term := prepare(cmd, name)
	term.Terminate()
	matchingFlags.Parse(args[1:])
	cmd.Run(matchingFlags.Args())
}

func (c *commander) Compgen(args []string, inword bool) (comp []string, err error) {

	// two cases: either we are completing the first arg, or we need to delegate to subcommands
	pos, prefix := compgen.Prefix(args, inword)
	if pos == 0 {
		//completing command
		vals := make([]string, 0, len(c.cmds))
		for k := range c.cmds {
			vals = append(vals, k)
		}
		// use a value based generator
		gen := compgen.ValueGen(vals)
		return gen(prefix), nil
	} else { // we are completing after <command>
		//so we get the get the command
		cont := c.getContainer(args)
		if cont == nil {
			return nil, nil
		}
		_, term := prepare(cont.command, c.name+" "+cont.name)
		//test if the command implements Argsgen
		// if cmdr, ok := cont.command.(compgen.Argsgen); ok {
		// 	//then prepare it
		// 	term.Argsgen(cmdr)
		// }
		return term.Compgen(args, inword)
	}
	return nil, nil
}

// Prints the usage.
func (c *commander) Usage(fs *flag.FlagSet) {
	name := c.name

	fmt.Fprintf(os.Stderr, "Usage: %s <command>\n\n", name)
	fmt.Fprintf(os.Stderr, "where <command> is one of:\n")

	//Sort conts by name
	conts := make([]*cmdCont, 0, len(c.cmds))
	for _, cont := range c.cmds {
		conts = append(conts, cont)
	}
	sort.Stable(byName(conts))

	w := tabwriter.NewWriter(os.Stderr, 8, 8, 2, ' ', 0)
	for _, cont := range conts {
		fmt.Fprintln(w, strings.Join([]string{"", cont.name, cont.syntax, cont.description, ""}, "\t"))
	}
	w.Flush()

	if numOfGlobalFlags(fs) > 0 {
		fmt.Fprintf(os.Stderr, "\navailable flags:\n")
		fs.PrintDefaults()
	}
	fmt.Fprintf(os.Stderr, "\n%s <command> -h for subcommand help\n", name)
}

// getContainer from the cmds
func (c *commander) getContainer(args []string) *cmdCont {
	if len(args) >= 1 {
		return c.cmds[args[0]]
	}
	return nil
}

// Returns the total number of registered flags in a flagset.
func numOfGlobalFlags(fs *flag.FlagSet) (count int) {
	if fs == nil {
		return 0
	}
	fs.VisitAll(func(flag *flag.Flag) {
		count++
	})
	return
}

//byName to sort sub command by their name
type byName []*cmdCont

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].name < a[j].name }
