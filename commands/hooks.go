package commands

import (
	"github.com/mislav/anyenv/cli"
	"github.com/mislav/anyenv/config"
	"os"
	"path/filepath"
	"strings"
)

var hooksHelp = `
Usage: $ProgramName hooks <command>

List hook scripts for a given $ProgramName command
`

func hooksCmd(args cli.Args) {
	commandName := args.Required(0)

	for _, script := range findHookScripts(commandName) {
		cli.Println(script)
	}
}

func findHookScripts(commandName string) []string {
	results := []string{}

	paths := append(strings.Split(os.Getenv("RBENV_HOOK_PATH"), ":"),
		filepath.Join(config.Root, "rbenv.d"),
		"/usr/local/etc/rbenv.d",
		"/etc/rbenv.d",
		"/usr/lib/rbenv/hooks")

	pluginPaths, err := filepath.Glob(config.PluginsDir().Join("*/etc/rbenv.d").String())
	if err != nil {
		panic(err)
	}
	paths = append(paths, pluginPaths...)

	for _, path := range paths {
		if path == "" {
			continue
		}
		hookScripts, err := filepath.Glob(filepath.Join(path, commandName, "*.bash"))
		if err != nil {
			panic(err)
		}
		for _, script := range hookScripts {
			results = append(results, script)
		}
	}

	return results
}

func init() {
	cli.Register("hooks", hooksCmd, hooksHelp)
}
