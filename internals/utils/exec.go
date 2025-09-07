package utils

import (
	"os/exec"
)

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
