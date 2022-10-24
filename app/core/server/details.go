package server

type Details struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Repository  string `yaml:"repository"`
	Environment string `yaml:"environment"`
	Hash        string `yaml:"commit_hash"`
	BuildDate   string `yaml:"build_date"`
	BuildEpoch  int    `yaml:"build_epoch_sec"`
}

func FetchApplicationDetails() Details {
	//yamlFile, err := Asset("resources/application.yml")
	details := Details{}
	// if err != nil {
	// 	panic(err)
	// }
	// yaml.Unmarshal(yamlFile, &details)
	return details
}
