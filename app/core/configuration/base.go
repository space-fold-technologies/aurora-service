package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/security"
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
	CertResolver                   string                         `yaml:"cert-resolver"`
	SessionDuration                time.Duration                  `yaml:"admin-session-duration"`
	EncryptionParameters           *security.EncryptionParameters `yaml:"encryption"`
}

// func ParseFromResource() Configuration {
// 	yamlFile, err := server.Asset("resources/settings.yml")
// 	config := Configuration{}
// 	if err != nil {
// 		fmt.Println("Failed to open the file")
// 		panic(err)
// 	}
// 	yaml.Unmarshal(yamlFile, &config)
// 	return config
// }

func ParseFromResource() Configuration {
	yamlFile, err := os.ReadFile(filepath.Join("./resources", "settings.yml"))
	config := Configuration{}
	if err != nil {
		fmt.Println("Failed to open the file")
		panic(err)
	}
	yaml.Unmarshal(yamlFile, &config)
	return config
}
