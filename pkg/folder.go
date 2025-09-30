package pkg

import (
	"errors"
	"os"
	"path"
)

func CreateLogsServicesFolderStructure(logsDir, serviceName string) (*os.File, *os.File, error) {
	err := os.MkdirAll(path.Join(logsDir, serviceName), 0750)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, nil, err
		}
	}
	stdoutFile := path.Join(logsDir, serviceName, "stdout.log")
	fstdout, err := os.Open(stdoutFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fstdout, err = os.Create(stdoutFile)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}
	stderrFile := path.Join(logsDir, serviceName, "stderr.log")
	fstderr, err := os.Open(stderrFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fstderr, err = os.Create(stderrFile)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}
	return fstdout, fstderr, nil
}
