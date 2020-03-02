package korectl

// resourceConfig stores custom resource CLI configurations
type resourceConfig struct {
	Name            string
	APIResourceName string
	Columns         []string
}

type resourceConfigMap map[string]resourceConfig

func (r resourceConfigMap) Get(name string) resourceConfig {
	if config, ok := r[name]; ok {
		return config
	}

	return resourceConfig{
		Name:            name,
		APIResourceName: name,
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	}
}

var teamResourceConfig = resourceConfig{
	Name:            "team",
	APIResourceName: "teams",
	Columns: []string{
		Column("Name", ".metadata.name"),
		Column("Description", ".spec.description"),
	},
}

var resourceConfigs = resourceConfigMap{
	"team":  teamResourceConfig,
	"teams": teamResourceConfig,
}
