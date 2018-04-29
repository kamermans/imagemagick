package test

import (
	"os"
	"os/exec"
)

// MockRun collects the command and arguments when Command() is run
type MockRun struct {
	cmd  string
	args []string
}

// Command returns the command that was run
func (run *MockRun) Command() string {
	return run.cmd
}

// Args returns the arguments that were passed to the commnad
func (run *MockRun) Args() []string {
	return run.args
}

// MockExec is a mock helper for `exec.Cmd`
type MockExec struct {
	helperName string
	runs       []*MockRun
}

// NewMockExec creates a new MockExec instance.  The `helperName` argument is the
// name of a function in the scope of your test that will be called in a separate
// thread by exec.Cmd (through go test).  It must be a valid test method, and it's
// stdOut, stdErr and exit code will be available on the real `exec.Cmd` object that
// is returned by MockExec.Command()
//
// The helper function should begin with this code to prevent it from being run
// as a normal test:
//
// 	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
// 		return
// 	}
//
// By convention, helper functions are named TestHelper<test>, where `<test>` is the
// name of the test that is using this helper.
func NewMockExec(helperName string) *MockExec {
	return &MockExec{
		helperName: helperName,
	}
}

func (mock *MockExec) addRun(cmd string, args []string) {
	run := &MockRun{
		cmd:  cmd,
		args: args,
	}
	mock.runs = append(mock.runs, run)
}

// RunCount returns the number of time a Command() was created
func (mock *MockExec) RunCount() int {
	return len(mock.runs)
}

// LastRun returns the most recent MockRun, or nil if it was never run
func (mock *MockExec) LastRun() *MockRun {
	if mock.RunCount() == 0 {
		return nil
	}

	return mock.runs[len(mock.runs)-1]
}

// Runs returns all the MockRun objects
func (mock *MockExec) Runs() []*MockRun {
	return mock.runs
}

// Command is a drop-in replacement for exec.Command and returns a valid *exec.Cmd
// that will be tracked and mocked, sending the requested command to a test method
// of your choosing (see NewMockExec).
func (mock *MockExec) Command(cmd string, args ...string) *exec.Cmd {
	mock.addRun(cmd, args)

	realArgs := []string{"-test.run=" + mock.helperName, "--", cmd}
	realArgs = append(realArgs, args...)
	realCmd := exec.Command(os.Args[0], realArgs...)
	realCmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return realCmd
}
