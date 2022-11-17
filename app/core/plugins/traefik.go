package plugins

import (
	"errors"
	"fmt"

	"github.com/space-fold-technologies/aurora-service/app/core/logging"
)

type ProxyOrder string

const (
	REVERSE_PROXY_REGISTRATION = "registration"
)

type TraefikPlugin struct {
	https            bool
	certResolverName string
	domain           string
	network          string
}

func NewTraefikPlugin(network, domain string, https bool, certResolverName string) Plugin {
	instance := new(TraefikPlugin)
	instance.network = network
	instance.https = https
	instance.certResolverName = certResolverName
	instance.domain = domain
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
	// Make this a docker service and keep the id, inject a custom label to use to query and destroy
	return nil
}
func (tp *TraefikPlugin) OnShutDown() error {
	// Shutdown the docker service and delete it [We are careful to have a dooms day switch]
	return nil
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
	target["traefik.http.services."+order.Hostname+".loadbalancer.server.port"] = fmt.Sprint(order.Port)
	target["traefik.port"] = fmt.Sprint(order.Port)
	rule := "traefik.http.routers." + order.Hostname + ".rule"
	target[rule] = fmt.Sprintf("Host(`%s`)", host)
	if tp.https {
		tlsLabel := fmt.Sprintf("traefik.http.routers.%s.tls", order.Hostname)
		target[tlsLabel] = "true"
		tlsResolverLabel := fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", order.Hostname)
		target[tlsResolverLabel] = tp.certResolverName
	}
	log := logging.GetInstance()
	for label, value := range result.Labels {
		log.Infof("LABEL: %s VALUE: %s", label, value)
	}
	return nil
}
