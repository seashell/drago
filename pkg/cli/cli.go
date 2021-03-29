package cli

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// Credits: https://github.com/mitchellh/cli

// CLI contains the state necessary to run subcommands and parse the
// command line arguments.
//
// CLI also supports nested subcommands, such as "cli foo bar". To use
// nested subcommands, the key in the Commands mapping below contains the
// full subcommand. In this example, it would be "foo bar".
//
// If you use a CLI with nested subcommands, some semantics change due to
// ambiguities:
//
//   * We use longest prefix matching to find a matching subcommand. This
//     means if you register "foo bar" and the user executes "cli foo qux",
//     the "foo" command will be executed with the arg "qux". It is up to
//     you to handle these args. One option is to just return the special
//     help return code `RunResultHelp` to display help and exit.
//
//   * The help flag "-h" or "-help" will look at all args to determine
//     the help function. For example: "otto apps list -h" will show the
//     help for "apps list" but "otto apps -h" will show it for "apps".
//     In the normal CLI, only the first subcommand is used.
//
//   * The help flag will list any subcommands that a command takes
//     as well as the command's help itself. If there are no subcommands,
//     it will note this. If the CLI itself has no subcommands, this entire
//     section is omitted.
//
//   * Any parent commands that don't exist are automatically created as
//     no-op commands that just show help for other subcommands. For example,
//     if you only register "foo bar", then "foo" is automatically created.
type CLI struct {
	config     *Config
	router     *Router
	helpFunc   HelpFunc
	helpWriter io.Writer
}

// Args contains parsed input to the CLI
type Args struct {
	rootFlags      []string
	subcommand     string
	subcommandArgs []string
	isHelp         bool
	isVersion      bool
}

// New returns a new CLI instance with sensible defaults.
func New(config *Config) *CLI {

	config = DefaultConfig().Merge(config)

	cli := &CLI{
		config:     config,
		router:     NewRouter(),
		helpFunc:   config.HelpFunc,
		helpWriter: config.HelpWriter,
	}

	// Populate router based on config object
	for name, cmd := range config.Commands {
		name = strings.TrimSpace(name)
		cli.router.AddCommand(name, cmd)
	}

	cli.router.AddMissingParents(func() Command {
		return &MockCommand{
			HelpText:      "This command is accessed by using one of the subcommands below.",
			RunReturnCode: CommandReturnCodeHelp,
		}
	})

	return cli
}

// WriteHelp :
func (cli *CLI) WriteHelp(data string) (int, error) {
	return cli.helpWriter.Write([]byte(data))
}

// Run :
func (cli *CLI) Run(ctx context.Context, args []string) (int, error) {

	parsed := cli.parseArgs(args)

	// Print version help
	if parsed.isVersion && cli.config.Version != "" {
		cli.WriteHelp(cli.config.Version + "\n")
		return 0, nil
	}

	// Print global help
	if parsed.isHelp && parsed.subcommand == "" {
		allSubcommands := cli.router.GetSubcommands(parsed.subcommand)
		cli.WriteHelp(cli.helpFunc(allSubcommands) + "\n")
		return 0, nil
	}

	// Print help for command parent
	command, err := cli.router.GetCommand(parsed.subcommand)
	if err != nil {
		parent := cli.router.GetParent(parsed.subcommand)
		parentSubcommands := cli.router.GetSubcommands(parent)
		cli.WriteHelp(cli.helpFunc(parentSubcommands) + "\n")
		return 1, nil
	}

	// Print help for subcommand
	if parsed.isHelp {
		cli.commandHelp(parsed.subcommand)
		return 0, nil
	}

	// If there is an invalid flag, then error
	if len(parsed.rootFlags) > 0 {
		cli.WriteHelp("Invalid flags before the subcommand. If these flags are for\n" +
			"the subcommand, please place them after the subcommand.\n\n")
		cli.commandHelp(parsed.subcommand)
		return 1, nil
	}

	code := command.Run(ctx, parsed.subcommandArgs)
	if code == CommandReturnCodeHelp {
		cli.commandHelp(parsed.subcommand)
		return 1, nil
	}

	return code, nil
}

func (cli *CLI) parseArgs(args []string) (out *Args) {

	out = &Args{}

	for i, arg := range args {

		if arg == "--" {
			break
		}

		if arg == "-h" || arg == "-help" || arg == "--help" {
			out.isHelp = true
			continue
		}

		// If no subcommand has been found yet, continue setting top-level flags
		if out.subcommand == "" {
			if arg == "-v" || arg == "-version" || arg == "--version" {
				out.isVersion = true
				continue
			}
			if arg != "" && arg[0] == '-' {
				out.rootFlags = append(out.rootFlags, arg)
			}
		}

		// If no subcommand has been found yet and we stumble upon a non-flag arg,
		// then it must be the subcommand.
		if out.subcommand == "" && arg != "" && arg[0] != '-' {
			out.subcommand = arg
			// If the command has a space in it, then it is invalid.
			// Set a blank command so that it fails.
			if strings.ContainsRune(arg, ' ') {
				out.subcommand = ""
				return
			}

			// Determine the argument we look to to end subcommands.
			// We look at all arguments until one has a space. This
			// disallows commands like: ./cli foo "bar baz". An argument
			// with a space is always an argument.
			j := 0
			for k, v := range args[i:] {
				if strings.ContainsRune(v, ' ') {
					break
				}
				j = i + k + 1
			}

			// Nested CLI, the subcommand is actually the entire
			// arg list up to a flag that is still a valid subcommand.
			searchKey := strings.Join(args[i:j], " ")
			k, _, ok := cli.router.GetLongestPrefix(searchKey)
			if ok {
				// k could be a prefix that doesn't contain the full
				// command such as "foo" instead of "foobar", so we
				// need to verify that we have an entire key. To do that,
				// we look for an ending in a space or an end of string.
				reVerify := regexp.MustCompile(regexp.QuoteMeta(k) + `( |$)`)
				if reVerify.MatchString(searchKey) {
					out.subcommand = k
					i += strings.Count(k, " ")
				}
			}
			// The remaining args the subcommand arguments
			out.subcommandArgs = args[i+1:]
		}

		// If a subcommand was not found, then use a default command if available
		if out.subcommand == "" {
			if _, err := cli.router.GetCommand(""); err == nil {
				out.rootFlags = nil
				out.subcommandArgs = append(out.rootFlags, out.subcommandArgs...)
			}
		}
	}

	return
}

func (cli *CLI) commandHelp(n string) {
	tpl := `
{{.Help}}
{{if gt (len .Subcommands) 0}}
Subcommands:
{{- range $value := .Subcommands }}
	{{ $value.NameAligned }}    {{ $value.Synopsis }}
{{- end }}
{{- end }}
`
	tpl = strings.TrimSpace(tpl)
	if !strings.HasSuffix(tpl, "\n") {
		tpl += "\n"
	}

	t, err := template.New("root").Parse(tpl)
	if err != nil {
		t = template.Must(template.New("root").Parse(fmt.Sprintf(
			"Failed to parse command help template: %s\n", err)))
	}

	cmd, err := cli.router.GetCommand(n)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"Name": cli.config.Name,
		"Help": cmd.Help(),
	}

	// Build subcommand list if we have it
	var subcommandsTpl []map[string]interface{}

	subcommands := cli.router.GetSubcommands(n)

	keys := make([]string, 0, len(subcommands))
	for k := range subcommands {
		keys = append(keys, k)
	}

	// Sort the keys
	sort.Strings(keys)

	// Figure out the padding length
	var longest int
	for _, k := range keys {
		if v := len(k); v > longest {
			longest = v
		}
	}

	// Go through and create their structures
	subcommandsTpl = make([]map[string]interface{}, 0, len(subcommands))
	for _, k := range keys {

		subcommand, ok := subcommands[k]
		if !ok {
			cli.WriteHelp(fmt.Sprintf(
				"Error getting subcommand %q", k))
		}

		name := k
		if idx := strings.LastIndex(k, " "); idx > -1 {
			name = name[idx+1:]
		}

		subcommandsTpl = append(subcommandsTpl, map[string]interface{}{
			"Name":        name,
			"NameAligned": name + strings.Repeat(" ", longest-len(k)),
			"Help":        subcommand.Help(),
			"Synopsis":    subcommand.Synopsis(),
		})
	}

	data["Subcommands"] = subcommandsTpl

	err = t.Execute(cli.helpWriter, data)
	if err == nil {
		return
	}

	cli.WriteHelp(fmt.Sprintf("Error rendering help: %s", err))
}
