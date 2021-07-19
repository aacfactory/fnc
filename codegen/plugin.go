package codegen

import (
	"fmt"
	"os"
	"os/exec"
)

type GenPlugin struct {
	Command string
	Debug   bool
}

func (p *GenPlugin) Execute(arg string) (err error) {

	cmd := exec.Command(p.Command, "--debug", fmt.Sprintf("%v", p.Debug), "-s", arg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	runErr := cmd.Start()
	if runErr != nil {
		err = fmt.Errorf("fnc call plugin %s failed, %v", p.Command, runErr)
		return
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		err = fmt.Errorf("fnc wait plugin %s executing failed, %v", p.Command, waitErr)
		return
	}

	return
}
