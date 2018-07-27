package dirtemplate

// DirTemplate is text template that can be used not only on one file
// but on a directory. It helps to create projects and generate code
// by template directory
type DirTemplate struct {
	Config
}

// New returns new instance of DirTemplate
func New(c Config) DirTemplate {
	return DirTemplate{
		Config: c,
	}
}

// Execute generates code for some language from Go-style templates using config
func (d DirTemplate) Execute() error {
	err := d.copyDir()
	if err != nil {
		return err
	}

	err = d.executeDir()
	if err != nil {
		return err
	}

	return nil
}
