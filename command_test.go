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
	"os"
	"testing"
)

// Tests if subcommand runs if it's provided as a part of arguments.
func TestCommand(t *testing.T) {
	resetForTesting("command1")

	c1 := &testCmd1{}
	On("command1", "", "", c1)
	Run()
	if !c1.run {
		t.Error("command 'command1' was expected to run, but it didn't")
	}
	if *c1.flag1 {
		t.Errorf("flag1 should be set to default: expected false, found %v", *c1.flag1)
	}
}

// Tests if subcommand runs and subcommand flags are set.
func TestCommandFlags(t *testing.T) {
	resetForTesting("command1", "-flag1=true")

	c1 := &testCmd1{}
	On("command1", "", "", c1)
	Run()
	if !c1.run {
		t.Error("command 'command1' was expected to run, but it didn't")
	}
	if !*c1.flag1 {
		t.Errorf("flag1 should be set: expected true, found %v", *c1.flag1)
	}
}

// Tests if correct subcommand runs if multiple subcommands
// are registered.
func TestMultiCommands(t *testing.T) {
	resetForTesting("command2")

	c1 := &testCmd1{}
	c2 := &testCmd2{}
	On("command1", "", "", c1)
	On("command2", "", "", c2)
	Run()
	if c1.run {
		t.Error("command 'command1' was not expected to run, but it did")
	}
	if !c2.run {
		t.Error("command 'command2' was expected to run, but it didn't")
	}
}

// Resets os.Args and the default flag set.
func resetForTesting(args ...string) {

	os.Args = append([]string{"cmd"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	CommandLine = New()
}

// testCmd1 is a test sub command.
type testCmd1 struct {
	flag1 *bool

	run bool
}

// Defines flags for the sub command.
func (cmd *testCmd1) Flags(fs *flag.FlagSet) {
	cmd.flag1 = fs.Bool("flag1", false, "Description about flag1")
}

// Sets the run flag.
func (cmd *testCmd1) Run(args []string) {
	cmd.run = true
}

// testCmd2 is a test sub command.
type testCmd2 struct {
	flag2 *bool

	run bool
}

// Defines flags for the sub command.
func (cmd *testCmd2) Flags(fs *flag.FlagSet) {
	cmd.flag2 = fs.Bool("flag2", false, "Description about flag2")
}

// Sets the run flag.
func (cmd *testCmd2) Run(args []string) {
	cmd.run = true
}
