package dirtemplate

import "text/template"

// Config stores template configs
type Config struct {
	SrcPath  string
	DestPath string

	LeftDelim  string
	RightDelim string

	Data interface{}

	FuncMap template.FuncMap

	SkipPaths    []string
	ReplacePaths map[string]string
}
