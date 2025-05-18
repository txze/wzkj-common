package surprise

import (
	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

var defaultTempl = `{{ .Title "WZKJ" "" 4 }}
{{ .AnsiColor.BrightCyan }}The title will be ascii and indented 4 spaces{{ .AnsiColor.Default }}
GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
GOARCH: {{ .GOARCH }}
NumCPU: {{ .NumCPU }}
GOPATH: {{ .GOPATH }}
GOROOT: {{ .GOROOT }}
Compiler: {{ .Compiler }}
ENV: {{ .Env "GOPATH" }}
Now: {{ .Now "Monday, 2 Jan 2006" }}`

type Banner struct {
	Templ string
}

var defaultBanner = &Banner{
	Templ: defaultTempl,
}

func NewBanner() *Banner {
	return &Banner{}
}

func RunBanner() {
	defaultBanner.Run()
}
func (b *Banner) Run() {
	banner.InitString(colorable.NewColorableStdout(), true, true, defaultTempl)
}

func (b *Banner) SetTempl(templ string) {
	b.Templ = templ
}
