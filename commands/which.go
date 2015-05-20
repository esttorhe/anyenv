package commands

import (
	"github.com/mislav/everyenv/cli"
	"github.com/mislav/everyenv/config"
	"github.com/mislav/everyenv/utils"
	"os"
	"strings"
)

var whichHelp = `
Usage: $ProgramName which <command>

Displays the full path to the executable that $ProgramName will invoke when
running the given command.
`

func whichCmd(args cli.Args) {
	currentVersion := detectVersion()
	exeName := args.Required(0)
	exePath := findExecutable(exeName, currentVersion)

	if exePath.IsBlank() {
		cli.Errorf("%s: command not found\n", exeName)
		versions := whence(exeName)
		if len(versions) > 0 {
			cli.Errorf("\nThe `%s' command exists in these versions:\n  %s\n", exeName,
				strings.Join(versions, "\n  "))
		}
		cli.Exit(127)
	} else {
		cli.Println(exePath)
	}
}

func findExecutable(exeName string, currentVersion SelectedVersion) utils.Pathname {
	if currentVersion.IsSystem() {
		filename := findInPath(exeName)
		if !filename.IsBlank() {
			return filename
		}
	} else {
		versionDir := config.VersionDir(currentVersion.Name)
		if !versionDir.Exists() {
			cli.Errorf("version `%s' is not installed\n", currentVersion.Name)
			cli.Exit(1)
		}
		filename := versionDir.Join("bin", exeName)
		if filename.IsExecutable() {
			return filename
		}
	}

	return utils.NewPathname("")
}

func findInPath(exeName string) utils.Pathname {
	shimsDir := config.ShimsDir()
	dirs := strings.Split(os.Getenv("PATH"), ":")

	var dir utils.Pathname
	var filename utils.Pathname

	for _, p := range dirs {
		dir = utils.NewPathname(p)
		if dir.IsBlank() || dir.Equal(shimsDir) {
			continue
		}
		filename = dir.Join(exeName)
		if filename.IsExecutable() {
			return filename.Abs()
		}
	}

	return utils.NewPathname("")
}

func init() {
	cli.Register("which", whichCmd, whichHelp)
}
