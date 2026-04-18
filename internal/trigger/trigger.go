package trigger

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

func Execute(triggers []string, oldip string, newip string) error {

	var exec_error bool = false

	for _, trigger := range triggers {
		logger.Infof("Executing trigger: %s %s %s", trigger, oldip, newip)

		if trigger == "" {
			logger.Warnf("Skipping empty trigger")
			exec_error = true
			continue
		}

		command := fmt.Sprintf("%s %s %s", trigger, strconv.Quote(oldip), strconv.Quote(newip))
		cmd := exec.Command("sh", "-c", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Errorln("Error while executing trigger:", err)
			logger.Debugf("Trigger output: %s", string(out))
			exec_error = true
			continue
		}
		logger.Debugf("Successfully executed trigger: %s", string(out))
	}

	if exec_error {
		return errors.New("exec error")
	} else {
		return nil
	}
}
