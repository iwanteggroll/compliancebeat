package scripting

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/elastic/beats/libbeat/logp"
)

type Script interface {
	Execute() ([]byte, error)
	Path() string
	Name() string
	Params() []string
	// SetName(string)
	// SetPath(string)
}

type PowershellScript struct {
	ScriptPath   string
	ScriptName   string
	ScriptParams []string
}

type PythonScript struct {
	ScriptPath   string
	ScriptName   string
	ScriptParams []string
}

func (pos *PowershellScript) Name() string {
	return pos.ScriptName
}

// func (pos *PowershellScript) SetName(name string) {
// 	pos.scriptName = name
// }

func (pos *PowershellScript) Path() string {
	return pos.ScriptPath
}

func (pos *PowershellScript) Params() []string {
	return pos.ScriptParams
}

// func (pos *PowershellScript) SetPath(path string) {
// 	pos.scriptPath = path
// }

func (pos *PowershellScript) Execute() (stdout []byte, err error) {
	_, err = exec.LookPath("powershell")

	if err != nil {
		logp.Err("Command cannot be found in path: %v", err)
		return
	}

	scPath := AbsScriptPath(pos)
	args := append([]string{"-NoLogo", "-File", scPath}, pos.ScriptParams...)
	//logp.NewLogger("powershellblock")
	//fmt.Println("AM I PRINTING?!")
	//logp.Debug("powershellblock", "Powershell running %s", pos.ScriptName)
	//logp.L().Infow("Test in execute", "powershellblock", logp.Info)
	cmd := exec.Command("powershell", args...)

	var waitStatus syscall.WaitStatus

	stdout, err = cmd.CombinedOutput()
	if cmd.ProcessState == nil {
		return
	}

	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)

	logp.Debug("examplebeat", "Command Returned: %q, exit code %d", stdout, waitStatus.ExitStatus())

	return stdout, err
}

func (pys *PythonScript) Name() string {
	return pys.ScriptName
}

// func (pys *PythonScript) SetName(name string) {
// 	pys.ScriptName = name
// }

func (pys *PythonScript) Path() string {
	return pys.ScriptPath
}

func (pys *PythonScript) Params() []string {
	return pys.ScriptParams
}

// func (pys *PythonScript) SetPath(path string) {
// 	pys.ScriptPath = path
// }

func (pys *PythonScript) Execute() (stdout []byte, err error) {
	// _, err = exec.LookPath("python3")

	// if err != nil {
	// 	logp.Err("Command cannot be found in path: %v", err)
	// 	return
	// }

	scPath := AbsScriptPath(pys)
	// for some reason cannot log
	logp.NewLogger("pythonblock")
	logp.Debug("pythonblock", "Python3 running %s", pys.ScriptName)

	var cmd *exec.Cmd

	args := append([]string{scPath}, pys.ScriptParams...)
	if runtime.GOOS == "windows" {
		_, err = exec.LookPath("python")

		if err != nil {
			logp.Err("Command cannot be found in path: %v", err)
			return
		}

		cmd = exec.Command("python", args...)
	} else {
		cmd = exec.Command(scPath, pys.ScriptParams...)
	}

	var waitStatus syscall.WaitStatus

	stdout, err = cmd.CombinedOutput()
	if cmd.ProcessState == nil {
		return
	}

	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)

	logp.L().Debugw("monitoring", "Command Returned: %q, exit code %d", stdout, waitStatus.ExitStatus())

	return stdout, err
}

func AbsScriptPath(sc Script) string {
	ex, err := os.Executable()
	if err != nil {
		logp.Err("Cannot find executable")
	}

	exPath := filepath.Dir(ex)

	scPath := filepath.Join(exPath, sc.Path(), sc.Name())

	return scPath
}
