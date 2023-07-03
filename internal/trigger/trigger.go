package trigger

import (
	"errors"
	"os/exec"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

func Execute(triggers []string, oldip string, newip string) error {

	var exec_error bool = false

	for _, trigger := range triggers {
		logger.Infof("Executing trigger: %s %s %s", trigger, oldip, newip)

		cmd := exec.Command(trigger, oldip, newip)
		out, err := cmd.Output()
		if err != nil {
			logger.Errorln("Error while executing trigger:", err)
			exec_error = true
		}
		logger.Debugf("Sucessfully executed trigger: %s", out)
	}

	if exec_error {
		return errors.New("exec error")
	} else {
		return nil
	}
}
