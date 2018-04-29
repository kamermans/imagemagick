package test_resources

import (
	"math/rand"
	"os"
	"os/exec"
)

type MockRun struct {
	cmd  string
	args []string
}

func (run *MockRun) GetCommand() string {
	return run.cmd
}

func (run *MockRun) GetArgs() []string {
	return run.args
}

type MockExec struct {
	helperName string
	runs       []*MockRun
}

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

func (mock *MockExec) RunCount() int {
	return len(mock.runs)
}

func (mock *MockExec) LastRun() *MockRun {
	if mock.RunCount() == 0 {
		return nil
	}

	return mock.runs[len(mock.runs)-1]
}

func (mock *MockExec) Runs() []*MockRun {
	return mock.runs
}

func (mock *MockExec) Command(cmd string, args ...string) *exec.Cmd {
	mock.addRun(cmd, args)

	realArgs := []string{"-test.run=" + mock.helperName, "--", cmd}
	realArgs = append(realArgs, args...)
	realCmd := exec.Command(os.Args[0], realArgs...)
	realCmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return realCmd
}
