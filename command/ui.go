package command

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// UICommand :
type UICommand struct {
	UI cli.UI

	Command
}

// Name :
func (c *UICommand) Name() string {
	return "acl"
}

// Synopsis :
func (c *UICommand) Synopsis() string {
	return "Open the Drago web UI"
}

// Run :
func (c *UICommand) Run(ctx context.Context, args []string) int {

	url := "http://127.0.0.1:8080"
	if c.address != "" {
		url = c.address
	}

	c.UI.Output(fmt.Sprintf(`Opening URL "%s"`, url))

	if err := openBrowser(url); err != nil {
		c.UI.Error(err.Error())
	}

	return 0
}

// Help :
func (c *UICommand) Help() string {
	h := `
Usage: drago ui [options] [args]

  Open the Drago web UI in the default browser.

General Options:
` + GlobalOptions() + `
`
	return strings.TrimSpace(h)
}

// Credits: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) error {

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil

}
