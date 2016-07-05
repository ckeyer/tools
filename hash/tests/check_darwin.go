package tests

import (
	"fmt"
	"os/exec"
)

func LocalCheck(method, filename string) (bool, error) {
	var cmd *exec.Cmd
	switch method {
	case "md5":
		cmd = exec.Command("md5", "-c", filename)
	case "sha1":
		cmd = exec.Command("shasum", "-a", "1", "-c", filename)
	case "sha256":
		cmd = exec.Command("shasum", "-a", "256", "-c", filename)
	case "sha512":
		cmd = exec.Command("shasum", "-a", "512", "-c", filename)
	}
	if cmd == nil {
		return false, fmt.Errorf("not support method: %s", method)
	}

	_, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, nil
}
