package agent

import (
	"fmt"

	"github.com/seashell/drago/version"
)

var bannerTmpl = fmt.Sprintf(`
====|===================>
___  ____ ____ ____ ____ 
|  \ |__/ |__| | __ |  | 
|__/ |  \ |  | |__] |__| 
		   
               {{ .AnsiColor.Cyan }}%s{{ .AnsiColor.Default }}
<===================|====

`, version.GetVersion().VersionNumber())
