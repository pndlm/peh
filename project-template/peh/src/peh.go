package main

import (
	"github.com/pndlm/peh/peh3"
	"pndlm.com/example-project/cmd"
)

func main() {
	proj := peh3.ProjectAtCwd("example-project")
	cmd.SetProject(proj)
	cmd.AddNamedCommand("down", proj.CmdDown())
	cmd.AddNamedCommand("exitedrm", proj.CmdExitedRm())
	cmd.AddNamedCommand("restart", proj.CmdRestart())
	cmd.AddNamedCommand("sh", proj.CmdSh())
	cmd.AddNamedCommand("tail", proj.CmdTail())
	cmd.Execute()
}
