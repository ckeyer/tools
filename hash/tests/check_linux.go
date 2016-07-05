package tests

import (
	"fmt"
	"os/exec"
)

func LocalCheck(method, filename string) (bool, error) {
	var cmd *exec.Cmd
	switch method {
	case "md5":
		cmd = exec.Command("md5sum", "-c", filename)
	case "sha1":
		cmd = exec.Command("sha1sum", "-c", filename)
	case "sha256":
		cmd = exec.Command("sha256sum", "-c", filename)
	case "sha512":
		cmd = exec.Command("sha512sum", "-c", filename)
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
