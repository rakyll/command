# command

[![Build Status](https://travis-ci.org/rakyll/command.png?branch=master)](https://travis-ci.org/rakyll/command)

command is a tiny package that helps you to add cli subcommands to your Go program.

## Usage

In order to start, go get this repository:

~~~ sh
go get github.com/rakyll/command
~~~

This package allows you to use flags package as you used to do, and provides additional parsing for subcommands and subcommand flags.

~~~ go
import "github.com/rakyll/command"

// register any global flags
var flagExecPath = flag.String("exec-path", "", "a custom path to executable")

type VersionCommand struct{
	flagVerbose *bool
}

func (cmd *VersionCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
	// define subcommand's flags
	cmd.flagVerbose = fs.Bool("v", false, "provides verbose output")
	return fs
}

func (cmd *VersionCommand) Run(args []string) {
	// implement the main body of the subcommand here
}

// register version as a subcommand
command.On("version", "prints the version", &VersionCommand{})
command.On("command1", "some description about command1", ...)
command.On("command2", "some description about command2", ...)
command.Parse()
// ...
command.Run()
~~~

The program above will handle the registered commands and invoke the matching command's `Run` or print subcommand help if `-h` is set.

~~~
$ program -exec-path=/home/user/bin/someexec version -v=true
~~~

will out the version of the program in a verbose way, and will set the exec path to the provided path. If arguments doesn't match any subcommand or illegal arguments are provided, it will print this beautiful usage guide:

~~~
$ program
Usage: program <command>

where <command> is one of:
  version   prints the version
  command1  some description about command1
  command2  some description about command2

available flags:
  -exec-path="": a custom path to executable

program <command> -h for subcommand help
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
