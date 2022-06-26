package envcheck

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/sys/execabs"
)

type task struct {
	name, hint string
	test       func() (string, error)
}

var tasks = []*task{
	{"Go compiler", "Checking Go is installed and up to date", func() (string, error) {
		path, err := execabs.LookPath("go")
		if err != nil {
			return "", err
		}

		cmd := execabs.Command(path, "version")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(ret)), nil
	}},
	{"C compiler", "Checking C is installed for system code", func() (string, error) {
		_, err := execabs.LookPath("gcc")
		if err == nil {
			return "gcc found", nil
		}

		_, err = execabs.LookPath("clang")
		if err == nil {
			return "clang found", nil
		}

		return "", errors.New("could not find gcc or clang compilers")
	}},
	{"Go bin PATH", "Verify that the PATH environment is set", func() (string, error) {
		cmd := execabs.Command("go", "env", "GOPATH")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		goPath := strings.TrimSpace(string(ret))

		allPath := os.Getenv("PATH")
		if !strings.Contains(allPath, goPath) {
			return "", errors.New("PATH missing " + goPath)
		}
		return "confirmed", nil
	}},
	//{"Dependencies installed", "Various libraries required are present", func() (string, error) {
	//	return "", errors.New("no dependencies found")
	//}},
	{"Fyne helper", "Checking Fyne tool is installed for packaging", func() (string, error) {
		path, err := execabs.LookPath("fyne")
		if err != nil {
			return "", err
		}

		cmd := execabs.Command(path, "--version")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(ret)), nil
	}},
}
