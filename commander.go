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
//		func (c *myCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
//     	   	c.Commander = command.NewCommander("git", fs)
//     	    c.On("status","Show the working tree status",&GitStatuCmd{}, nil)
//      	...
//		}
//
// Works out of the box.
type Commander interface {
	//Usage print the usage for this commander
	Usage()
	// Registers a Cmd for the provided sub-command name. E.g. name is the
	// `status` in `git status`.
	On(name, description string, command Cmd, requiredFlags []string)
	// Parses the flags and leftover arguments to match them with a
	// sub-command. Evaluate all of the global flags and register
	// sub-command handlers before calling it. Sub-command handler's
	// `Run` will be called if there is a match.
	// A usage with flag defaults will be printed if provided arguments
	// don't match the configuration.
	// Global flags are accessible once Parse executes.
	Parse()
	// Runs the subcommand's runnable. If there is no subcommand
	// registered, it silently returns.
	Run(args []string)
	//Set Command Name as it should appear in the doc.
	SetName(name string)
}

// Execute a Cmd as a main command
//
//For example:
//
//    Exec(&myCmd{}, os.Args)
//
// If you have write a Commander as a Cmd object (to be used recursively) and you want to
// use it as a main, this method is for you.
func Exec(cmd Cmd, args []string) {
	// init the cmd's flagset
	fs := cmd.Flags(flag.NewFlagSet(args[0], flag.ExitOnError))
	//and use it to parse the args
	fs.Parse(args[1:])

	// recursive case: cmd is also a commander,
	if cmdr, ok := cmd.(Commander); ok {
		//pass on the name
		cmdr.SetName(args[0])
		// let the commander parse the args too, (to find out the actual subcommand)
		cmdr.Parse()
	}
	// in any case, run it
	cmd.Run(fs.Args())
}

type commander struct {
	//a commander is made of:
	// - its own name (defined from outside using setName)
	// - the command's flagset (for self options)
	// - the map of registered subcommands
	// - the subcommand matched during "parse" stage
	// - the args (build during parse stage) to be passed to the subcommand

	name  string
	flags *flag.FlagSet // this command flags

	// A map of all of the registered sub-commands.
	cmds map[string]*cmdCont

	// Matching subcommand.
	matchingCmd   *cmdCont
	matchingFlags *flag.FlagSet

	// Flag to determine whether help is
	// asked for subcommand or not
	flagHelp *bool
}

//NewCommander creates a new Commander. The 'name' is the the subcommand name,
// and the fs is the Cmd current flagset.
//
// A NewCommander is better created in the Flag(fs *flag.FlagSet) method.
func NewCommander(name string, fs *flag.FlagSet) Commander {

	c := &commander{
		name:  name,
		cmds:  make(map[string]*cmdCont),
		flags: fs,
	}
	//a command record the usage to be declared
	fs.Usage = c.Usage
	return c
}

func (c *commander) SetName(name string) { c.name = name }

// Registers a Cmd for the provided sub-command name. E.g. name is the
// `status` in `git status`.
func (c *commander) On(name, description string, command Cmd, requiredFlags []string) {
	c.cmds[name] = &cmdCont{
		name:          name,
		desc:          description,
		command:       command,
		requiredFlags: requiredFlags,
	}
}

// Parses the flags and leftover arguments to match them with a
// sub-command. Evaluate all of the global flags and register
// sub-command handlers before calling it. Sub-command handler's
// `Run` will be called if there is a match.
// A usage with flag defaults will be printed if provided arguments
// don't match the configuration.
// Global flags are accessible once Parse executes.
func (c *commander) Parse() {

	// if there are no subcommands registered,
	// return immediately
	if len(c.cmds) < 1 {
		return
	}

	if c.flags.NArg() < 1 {
		c.flags.Usage()
		os.Exit(1)
	}

	name := c.flags.Arg(0)

	if cont, ok := c.cmds[name]; ok { //this is an existing command

		// Init it
		c.matchingFlags = cont.command.Flags(flag.NewFlagSet(name, flag.ExitOnError))
		// always append a -h option to print "help"
		c.flagHelp = c.matchingFlags.Bool("h", false, "")
		c.matchingFlags.Parse(c.flags.Args()[1:])
		c.matchingCmd = cont

		// recursive case: if it's also a commander (it has subcommands)
		if cmdr, ok := cont.command.(Commander); ok {
			cmdr.SetName(c.name + " " + name)
			//we had to split setName, and Parse because:

			// checking for required flags might target subCommandUsage (that needs the name to be set)
			// but this check need to be made before parsing the subcommand
		}

		// Check for required flags.
		flagMap := make(map[string]bool)
		for _, flagName := range cont.requiredFlags {
			flagMap[flagName] = true
		}
		c.matchingFlags.Visit(func(f *flag.Flag) {
			delete(flagMap, f.Name)
		})
		if len(flagMap) > 0 { // missed a required flag
			c.subcommandUsage(c.matchingCmd)
			os.Exit(1)
		}

		// ok this one is good. Now, if this "Cmd" is also a Commander go on parsing
		if cmdr, ok := cont.command.(Commander); ok {
			cmdr.Parse() // fs has been parsed
		}

	} else {
		c.flags.Usage()
		os.Exit(1)
	}
}

//Implement the Cmd .Run method so that a Commander is almsot a Cmd too.
// the args is ignored.
func (c *commander) Run(args []string) {
	c.Exec()
}

// Exec the subcommand's runnable. If there is no subcommand
// registered, it silently returns.
//args is totally ignored and kept to be
func (c *commander) Exec() {
	if c.matchingCmd != nil {
		if *c.flagHelp {
			c.subcommandUsage(c.matchingCmd)
			return
		}
		c.matchingCmd.command.Run(c.matchingFlags.Args())
	}
}

// Prints the usage.
func (c *commander) Usage() {
	name := c.name
	if len(c.cmds) == 0 {
		// no subcommands
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
		c.flags.PrintDefaults()
		return
	}

	fmt.Fprintf(os.Stderr, "Usage: %s <command>\n\n", name)
	fmt.Fprintf(os.Stderr, "where <command> is one of:\n")

	//Sort conts by name
	conts := make([]*cmdCont, 0, len(c.cmds))
	for _, cont := range c.cmds {
		conts = append(conts, cont)
	}
	sort.Stable(byName(conts))

	for _, cont := range conts {
		fmt.Fprintf(os.Stderr, "  %-15s %s\n", cont.name, cont.desc)
	}

	if numOfGlobalFlags(c.flags) > 0 {
		fmt.Fprintf(os.Stderr, "\navailable flags:\n")
		c.flags.PrintDefaults()
	}
	fmt.Fprintf(os.Stderr, "\n%s <command> -h for subcommand help\n", name)
}

func (c *commander) subcommandUsage(cont *cmdCont) {
	name := c.name
	fmt.Fprintf(os.Stderr, "  %s %-15s %s\n", name, cont.name, cont.desc)

	// should only output sub command flags, ignore h flag.
	fs := cont.command.Flags(flag.NewFlagSet(cont.name, flag.ContinueOnError))
	if cmdr, ok := cont.command.(Commander); ok {
		cmdr.SetName(c.name + " " + cont.name)
	}
	if fs.Usage != nil { // if the cmd has defined a usage, use it.
		fs.Usage()
		return
	}
	// else
	if numOfGlobalFlags(fs) > 0 {
		fmt.Fprintf(os.Stderr, "Usage %s:\n", cont.name)
		fs.PrintDefaults()
	}
	if len(cont.requiredFlags) > 0 {
		fmt.Fprintf(os.Stderr, "\nrequired flags:\n")
		fmt.Fprintf(os.Stderr, "  %s\n\n", strings.Join(cont.requiredFlags, ", "))
	}
}

// Returns the total number of registered flags in a flagset.
func numOfGlobalFlags(fs *flag.FlagSet) (count int) {
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
