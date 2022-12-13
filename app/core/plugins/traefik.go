package plugins

import (
	"errors"
	"fmt"

	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
)

type ProxyOrder string

const (
	REVERSE_PROXY_REGISTRATION = "registration"
)

var (
	LOGS_PATH         = "/etc/traefik/logs"
	LETS_ENCRYPT_PATH = "/etc/traefik/acme"
)

type TraefikPlugin struct {
	https                  bool
	certResolverName       string
	certResolverEmail      string
	certResolverParameters map[string]string
	domain                 string
	network                string
	provider               providers.Provider
	identifier             string
}

func NewTraefikPlugin(network, domain string, https bool, certResolverName, certResolverEmail string, certResolverParameters map[string]string, provider providers.Provider) Plugin {
	instance := new(TraefikPlugin)
	instance.network = network
	instance.https = https
	instance.certResolverName = certResolverName
	instance.certResolverEmail = certResolverEmail
	instance.certResolverParameters = certResolverParameters
	instance.domain = domain
	instance.provider = provider
	return instance
}

type ProxyRequest struct {
	Hostname string
	Port     uint
}

type ProxyResponse struct {
	Labels map[string]string
}

func (tp *TraefikPlugin) Name() string {
	return "traefik-plugin"
}
func (tp *TraefikPlugin) Category() PluginCategory {
	return REVERSE_PROXY
}
func (tp *TraefikPlugin) OnStartUp() error {
	ports := make([]int, 0)
	ports = append(ports, 80, 8080)
	if tp.https {
		ports = append(ports, 443)
	}
	if identifier, err := tp.provider.DeployDependency(&providers.DependencyOrder{
		ID:                   "traefik-plugin",
		Name:                 "internal-traefik-proxy",
		URI:                  "traefik:2.9.6",
		Digest:               "sha256:0f22b9ca5a3bacd50c43ca193e1970d825fe2320ffc6222fec1e2e81b9a23393",
		Ports:                ports,
		Volumes:              tp.mounts(),
		Command:              tp.commands(),
		EnvironmentVariables: tp.certResolverParameters,
	}); err != nil {
		return err
	} else {
		tp.identifier = identifier
		logging.GetInstance().Infof("TRAEFIK SERVICE ID : [%s]", tp.identifier)
	}
	return nil
}

func (tp *TraefikPlugin) OnShutDown() error {
	return tp.provider.Stop(tp.identifier)
}

func (tp *TraefikPlugin) Call(operation string, request interface{}, response interface{}) error {
	if operation != REVERSE_PROXY_REGISTRATION {
		return fmt.Errorf("wrong operation called in plugin")
	}
	if order, ok := request.(*ProxyRequest); !ok {
		return fmt.Errorf("request type is not compatible with call")
	} else if result, ok := response.(*ProxyResponse); !ok {
		return fmt.Errorf("respose type is not compatible with call")
	} else {
		return tp.register(order, result)
	}
}

func (tp *TraefikPlugin) register(order *ProxyRequest, result *ProxyResponse) error {
	if result.Labels == nil {
		return errors.New("target labels are not initialized")
	} else if order.Port == 0 {
		return errors.New("no ports specified")
	}
	host := fmt.Sprintf("%s.%s", order.Hostname, tp.domain)
	target := result.Labels
	// pack in the required traefik labels
	target["traefik.enable"] = "true"
	target["traefik.docker.network"] = tp.network
	target["traefik.http.routers."+order.Hostname+".entrypoints"] = "web"
	target["traefik.http.services."+order.Hostname+".loadbalancer.server.port"] = fmt.Sprint(order.Port)
	rule := "traefik.http.routers." + order.Hostname + ".rule"
	target[rule] = fmt.Sprintf("Host(`%s`)", host)
	if tp.https {
		tlsLabel := fmt.Sprintf("traefik.http.routers.%s.tls", order.Hostname)
		target[tlsLabel] = "true"
		tlsResolverLabel := fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", order.Hostname)
		target[tlsResolverLabel] = tp.certResolverName
	}
	return nil
}

func (tp *TraefikPlugin) commands() []string {
	command := []string{"traefik"}
	command = append(command, "--entrypoints.web.address=:80")
	if tp.https {
		command = append(
			command,
			"--entrypoints.websecure.address=:443",
			fmt.Sprintf("--certificatesresolvers.%s.acme.email=%s", tp.certResolverName, tp.certResolverEmail),
			fmt.Sprintf("--certificatesresolvers.%s.acme.storage=/letsencrypt/acme.json", tp.certResolverName),
			fmt.Sprintf("--certificatesresolvers.%s.acme.httpchallenge.entrypoint=web", tp.certResolverName),
			// Redirect HTTP -> HTTPS: https://doc.traefik.io/traefik/routing/entrypoints/#redirection
			"--entrypoints.web.http.redirections.entrypoint.to=websecure",
			"--entrypoints.web.http.redirections.entrypoint.scheme=https")
	}
	command = append(command,
		"--api.insecure=true",
		"--providers.docker.endpoint=unix:///var/run/docker.sock",
		"--providers.docker.swarmMode=true",
		"--providers.docker.swarmModeRefreshSeconds=10",
		"--providers.docker.httpClientTimeout=300",
		"--providers.docker.exposedbydefault=false",
		fmt.Sprintf("--providers.docker.network=%s", tp.network),
		"--accesslog=true")
	command = append(command,
		"--metrics.prometheus=true",
		"--metrics.prometheus.addEntryPointsLabels=true",
		"--metrics.prometheus.addrouterslabels=true",
		"--metrics.prometheus.addServicesLabels=true",
		"--entryPoints.metrics.address=:8082",
		"--metrics.prometheus.entryPoint=metrics")
	return command
}

func (tp *TraefikPlugin) mounts() map[string]string {
	volumes := make(map[string]string)
	volumes["/var/run/docker.sock"] = "/var/run/docker.sock"
	volumes[LOGS_PATH] = "/var/logs"
	if tp.https {
		volumes[LETS_ENCRYPT_PATH] = "/letsencrypt"
	}
	return volumes
}
