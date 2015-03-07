[![Build Status](https://travis-ci.org/ericaro/command.png?branch=master)](https://travis-ci.org/ericaro/command) [![GoDoc](https://godoc.org/github.com/ericaro/command?status.svg)](https://godoc.org/github.com/ericaro/command)

`command` is a tiny package that helps you to add cli subcommands to your Go program.

This work is a derivative of [rakyll's](https://github.com/rakyll/command) `command` library as an attempt to add:
- **Modularity**: flags, recursion, completion are all optionals
- **Recursivity**: commands can have subcommands and so on 
- **Completion**: executable are  native bash completion commands.

## Usage

~~~ sh
go get github.com/ericaro/command
~~~

### The Simplest commands

~~~ go

    package main

    import "github.com/ericaro/command"
    import "fmt"

    type HelloCommand struct{}

    func (cmd *HelloCommand) Run(args []string) { fmt.Printf("hello %s\n", args) }
    func main() {
      // register hello as a subcommand
      command.On("hello", "<name>", "prints hello <name>", &HelloCommand{})
      command.Run()
    }

~~~

Fully functional Helloworld (with autocompletion builtin)



### Adding flags

Simply make your command implement the `command.Flagger` interface

~~~ go

type VersionCommand struct{
  flagVerbose *bool
}

func (cmd *VersionCommand) Flags(fs *flag.FlagSet) {
  // define subcommand's flags
  cmd.flagVerbose = fs.Bool("v", false, "provides verbose output")
}

// everything else is unchanged

~~~

### Autocompletion

See [compgen package](https://github.com/ericaro/compgen) for more details.

When using `command`, executables are built with completion capabilities. 

Therefore you just need to register them as their own completion command:

~~~ bash
$ complete -C cmd cmd
~~~

You can copy this statement in a file into `/etc/bash_completion.d/` to make it persistent.

By default completion works with
- subcommands
- flags names
- flags default value

It is possible though to finely tune it. Simple make your command implement the `command.Completer` interface.

~~~ go

    type VersionCommand struct{}
    
    func (cmd *VersionCommand) Compgens(term *compgen.Terminator) {
      term.Flag("d", compgen.CompgenCmd("directory") )
    }
    
    // everything else is unchanged

~~~

Now the "-d" flag will be completed with a local directories.

See [compgen package](https://github.com/ericaro/compgen) for more details.

## Recursive Commands

~~~ go

`command` exposes a Commander interface that can be configured to support any subcommands

    type RemoteCommand struct{
      command.Commander
    }

    // create and add a "commander"
    remote := command.New()
    command.On("remote", "<command>", "remote subcommands", remote)
    // and configure it
    remote.On("add", "<url>", "add a remote by url", adderCmd{})
    remote.On("remove", "<url>", "add a remote by url", removeCmd{})
~~~

~~~ bash
$ cmd remote add http://github.com/ericaro/command
~~~

## License

Copyright 2013 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
