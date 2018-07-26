package dirtemplate

type DirTemplate struct {
	Config
}

func New(c Config) *DirTemplate {
	return &DirTemplate{
		Config: c,
	}
}
