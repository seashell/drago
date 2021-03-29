package agent

import (
	"fmt"

	"github.com/seashell/drago/version"
)

// Banner is a banner to be displayed when the Drago
// agent is started
var Banner = fmt.Sprintf(`
====|===================>
___  ____ ____ ____ ____ 
|  \ |__/ |__| | __ |  | 
|__/ |  \ |  | |__] |__| 
		   
                  {{ .AnsiColor.Cyan }}%s{{ .AnsiColor.Default }}
<===================|====

`, version.GetVersion().VersionNumber())
