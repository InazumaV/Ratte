package handler

import (
	"os"
	"os/exec"
)

func (h *Handler) execHookCmd(cmd string, ares ...string) error {
	c := exec.Command(cmd, ares...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
