package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func Start(args ...string) (p *os.Process, err error) {
	if args[0], err = exec.LookPath(args[0]); err == nil {

		var procAttr os.ProcAttr
		// procAttr.Files = []*os.File{
		// 	os.Stdin,
		// 	os.Stdout,
		// 	os.Stderr,
		// }

		sys := syscall.SysProcAttr{
			Setsid: true,
		}

		cwd, _ := os.Getwd()

		procAttr = os.ProcAttr{
			Dir:   cwd,
			Env:   os.Environ(),
			Sys:   &sys,
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		}

		// procAttr.Sys.Ptrace = true

		p, err := os.StartProcess(args[0], args, &procAttr)
		if err == nil {
			return p, nil
		}
	}

	return nil, err
}

func WookPath() ([]string, []string) {
	path := os.Getenv("PATH")
	var dirs []string
	var exes []string
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		dirs = append(dirs, dir)
		dirent, err := os.ReadDir(dir)
		if err != nil {
		}
		for x := range dirent {
			exe := dirent[x].Name()
			exes = append(exes, exe)
		}

	}

	return dirs, exes
}

func main() {
	_, exes := WookPath()

	cmd := "fzf --multi"
	c := exec.Command("sh", "-c", cmd)
	p, err := c.StdinPipe()

	if err != nil {
		fmt.Println(err)
	}

	c.Stderr = os.Stderr
	defer p.Close()

	go func(o io.WriteCloser) {
		for ex := range exes {
			// fmt.Println(exes[ex])
			// b := []byte(exes[ex] + "\n")
			b := []byte(exes[ex])
			p.Write(b)
			p.Write([]byte("\n"))
		}
	}(p)

	out, err := c.Output()
	if err != nil {
		fmt.Println(err)
	}

	executable := string(out)
	executable = strings.ReplaceAll(executable, "\n", "")
	fmt.Println(executable)

	if proc, err := Start(executable); err == nil {
		fmt.Println("started process: ", proc.Pid)
		// proc.Wait()
	}
}
