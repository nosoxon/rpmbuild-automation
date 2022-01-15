package createrepo

import (
	"bytes"
	"os/exec"
)

type result struct {
	Output string `json:"output"`
	Error  string `json:"error"`
	Status int	  `json:"status"`
}

func Run(args ...string) (*result, error) {
	var output, errput bytes.Buffer
	cmd := exec.Command("createrepo", args...)
	cmd.Stdout = &output
	cmd.Stderr = &errput
	err := cmd.Run()

	return &result{
		Output: output.String(),
		Error: errput.String(),
		Status: cmd.ProcessState.ExitCode(),
	}, err
}