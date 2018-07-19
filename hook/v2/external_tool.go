package v2

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

// runExternalTools iterates over external tools to validate commit
func (cfg *Configuration) runExternalTools() bool {
	result := true
	for _, val := range cfg.Externals {
		cmd := exec.Command(val.Command[0], val.Command[1:]...)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			fmt.Println(fmt.Sprintf("error: unable to start [%s]: %s", val.Name, err.Error()))
			result = false
		} else {
			if err := cmd.Wait(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					fmt.Println("output:")
					fmt.Println()
					fmt.Println(stdout.String())

					fmt.Println("error:")
					fmt.Println()
					fmt.Println(stderr.String())
					if _, ok := exitError.Sys().(syscall.WaitStatus); ok {
						if val.Severity == ErrorSeverity {
							result = false
						}
						fmt.Println(fmt.Sprintf("%s: execution of [%s] failed", val.Severity, val.Name))
					}
				} else {
					fmt.Println(fmt.Sprintf("error: unable to start [%s]: %s", val.Name, err.Error()))
					result = false
				}
			}
		}
	}
	return result
}