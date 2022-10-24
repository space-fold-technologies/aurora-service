package configuration

import (
	"fmt"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/core/server"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Host                  string                         `yaml:"host"`
	Port                  int                            `yaml:"port"`
	Provider              string                         `yaml:"provider"`
	ProfileDIR            string                         `yaml:"profile-directory"`
	DefaultCluster        string                         `yaml:"default-cluster"`
	DefaultClusterAddress string                         `yaml:"default-cluster-address"`
	DefaultTeam           string                         `yaml:"default-team"`
	Domain                string                         `yaml:"domain"`
	Https                 bool                           `yaml:"https"`
	CertResolver          string                         `yaml:"cert-resolver"`
	SessionDuration       time.Duration                  `yaml:"admin-session-duration"`
	EncryptionParameters  *security.EncryptionParameters `yaml:"encryption"`
}

func ParseFromResource() Configuration {
	yamlFile, err := server.Asset("resources/settings.yml")
	config := Configuration{}
	if err != nil {
		fmt.Println("Failed to open the file")
		panic(err)
	}
	yaml.Unmarshal(yamlFile, &config)
	return config
}
