package compose

type compose struct {
	Version  string         `yaml:version`
	Services map[string]any `yaml:services`
}
