// Copyright © 2016 Eric Larson <eric@ionrock.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package buildpack

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LoggableCommand(dir string, c ...string) error {
	cmd := exec.Command(c[0], c[1:]...)
	cmd.Dir = dir
	o, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdout pipe: ", err)
	}

	e, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderr pipe: ", err)
	}

	stdout := bufio.NewScanner(o)
	stderr := bufio.NewScanner(e)
	go func() {
		for stdout.Scan() {
			log.Info(stdout.Text())
		}
	}()

	go func() {
		for stderr.Scan() {
			log.Info(stderr.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal("Error starting cmd: ", err)
	}

	return cmd.Wait()
}

func Workspace() string {
	dir, err := ioutil.TempDir("workspaces", "")
	if err != nil {
		log.Fatal("Error creating workspace: ", err)
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		log.Fatal("Unable to find the abspath of the ws: ", err)
	}

	return dir
}

func Checkout(repo string, ws string) error {
	log.Infof("Running: git clone %s %s", repo, ws)
	return LoggableCommand(ws, "git", "clone", repo, ws)
}

type Cmd struct {
	command string
}

func (c *Cmd) Do(dir string) error {
	parts := strings.Split(c.command, " ")
	return LoggableCommand(dir, parts...)
}

type Buildpack struct {
	Name      string
	Dir       string
	Files     []os.FileInfo
	Bootstrap Cmd
	Build     Cmd
	Test      Cmd
	Run       Cmd
}

func (bp *Buildpack) Install(ws string) error {
	for _, fn := range bp.Files {
		src := filepath.Join(bp.Dir, fn.Name())
		dest := filepath.Join(ws, path.Base(fn.Name()))
		log.Infof("Installing: %s to %s", src, dest)
		err := os.Link(src, dest)

		if err != nil {
			return err
		}
	}
	return nil
}

func Find(bp string) *Buildpack {
	pack := Buildpack{Name: bp}
	pack.Dir = filepath.Join("packs", bp)
	files, err := ioutil.ReadDir(pack.Dir)
	if err != nil {
		log.Fatal("Error reading pack files: ", err)
	}
	pack.Files = files

	b, err := ioutil.ReadFile(filepath.Join(pack.Dir, "cmds.yml"))
	if err != nil {
		log.Fatal(err)
	}

	c := make(map[string]string)

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}

	pack.Bootstrap = Cmd{command: c["bootstrap"]}
	pack.Build = Cmd{command: c["build"]}
	pack.Test = Cmd{command: c["test"]}
	pack.Run = Cmd{command: c["run"]}

	return &pack
}

func Setup(repo string, buildpack string) (string, *Buildpack) {
	ws := Workspace()
	log.Info("Using workspace: ", ws)

	err := Checkout(repo, ws)
	if err != nil {
		log.Error("Error checking out workspace.")
		log.Fatal(err)
	}
	log.Infof("Checked out %s into %s", repo, ws)

	bp := Find(buildpack)

	err = bp.Install(ws)
	if err != nil {
		log.Fatal("Error installing buildpack: ", err)
	}

	return ws, bp
}

func Build(repo string, buildpack string) {
	ws, bp := Setup(repo, buildpack)

	log.Info("Running Bootstrap")
	bp.Bootstrap.Do(ws)
}

func Test(repo string, buildpack string) {
	ws, bp := Setup(repo, buildpack)

	log.Info("Running Bootstrap")
	bp.Bootstrap.Do(ws)

	log.Info("Running Test")
	bp.Test.Do(ws)
}
