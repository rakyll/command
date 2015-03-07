[![Build Status](https://travis-ci.org/ericaro/command.png?branch=master)](https://travis-ci.org/ericaro/command) [![GoDoc](https://godoc.org/github.com/ericaro/command?status.svg)](https://godoc.org/github.com/ericaro/command)

This library is fully `go gettable`.

command is a tiny package that helps you to add cli subcommands to your Go program with no effort, and prints a pretty guide if needed.


This work is a derivative of [rakyll's](https://github.com/rakyll/command) command library.

Mainly to make it:
- **Modular**: flags, completion, recursion are all optionals
- **Recursive**: commands can have subcommands and so on 
- **autocompletion**: mode compatible with bash completion

## Usage

get go an `go get`

~~~ sh
go get github.com/ericaro/command
~~~


### Simplest commands

~~~ go

    import "github.com/ericaro/command"
     
     type VersionCommand struct{}
     
     func (cmd *VersionCommand) Run(args []string) {
       // implement the main body of the subcommand here
       // arguments are found in args
     }
     
     
     // register version as a subcommand
     command.On("version", "", prints the version", &VersionCommand{})
     command.On("command1","[-option] <arguments>", "some description about command1")   
     command.On("command2","[-option] <arguments>", "some description about command2")
     // ...
     command.Run()
~~~

That's it. It works.

### Adding autocompletion

See [compgen package](https://github.com/ericaro/compgen) for more details.

When using `command` executable are builtin with completion capability. So you just need to register them as their own completion command:

~~~ bash
$ complete -C cmd cmd
~~~

You can copy this statement in a file into `/etc/bash_completion.d/` to make it persistent.

By default completion works with
- subcommands
- flags names
- flags default value

It is possible though to configure it (see below)

### Adding flags

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

### Hacking autocomplete

Commands come with a default support for completion (see above)

It is possible to hack in:

~~~ go

    type VersionCommand struct{}
    
    func (cmd *VersionCommand) Compgens(term *compgen.Terminator) {
      term.Flag("d", compgen.CompgenCmd("directory") )
    }
    
    // everything else is unchanged

~~~

Now the "-d" flag will be auto-completed with local directories.

See [compgen package](https://github.com/ericaro/compgen) for more details.


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
