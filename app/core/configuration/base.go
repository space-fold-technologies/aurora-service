package configuration

import (
	"fmt"
	"os"

	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/core/server"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Host                           string                         `yaml:"host"`
	Port                           int                            `yaml:"port"`
	Provider                       string                         `yaml:"provider"`
	ProfileDIR                     string                         `yaml:"profile-directory"`
	DefaultCluster                 string                         `yaml:"default-cluster"`
	DefaultClusterListenAddress    string                         `yaml:"default-cluster-listen-address"`
	DefaultClusterAdvertiseAddress string                         `yaml:"default-cluster-advertise-address"`
	DefaultTeam                    string                         `yaml:"default-team"`
	Domain                         string                         `yaml:"domain"`
	Https                          bool                           `yaml:"https"`
	CertResolverName               string                         `yaml:"cert-resolver-name"`
	CertResolverEmail              string                         `yaml:"cert-resolver-email"`
	SessionDuration                int                            `yaml:"admin-session-duration"`
	EncryptionParameters           *security.EncryptionParameters `yaml:"encryption"`
	AgentParameters                *providers.ClientConfiguration `yaml:"agent-client"`
	NetworkPrefix                  string                         `yaml:"network-prefix"`
	NetworkName                    string                         `yaml:"network-name"`
	Plugins                        []string                       `yaml:"plugins"`
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

func ParseFromPath(filePath string) Configuration {
	config := Configuration{}
	if data, err := os.ReadFile(filePath); err != nil {
		fmt.Println("Failed to open the file")
		panic(err)
	} else if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return config
}

// func ParseFromResource() Configuration {
// 	yamlFile, err := os.ReadFile(filepath.Join("./resources", "settings.yml"))
// 	config := Configuration{}
// 	if err != nil {
// 		fmt.Println("Failed to open the file")
// 		panic(err)
// 	}
// 	yaml.Unmarshal(yamlFile, &config)
// 	return config
// }
