package plugins

import (
	"errors"
	"fmt"
)

type TraefikPlugin struct {
	Https            bool
	CertResolverName string
}

func (tp *TraefikPlugin) Labels(target map[string]string, network, host, hostname string, ports []uint) error {
	if target == nil {
		return errors.New("target labels are not initialized")
	} else if len(ports) == 0 {
		return errors.New("no ports specified")
	}
	// pack in the required traefik labels
	target["traefik.enable"] = "true"
	target["traefik.docker.network"] = network
	target["traefik.http.services."+hostname+".loadbalancer.server.port"] = fmt.Sprint(ports[0])
	target["traefik.port"] = fmt.Sprint(ports[0])
	rule := "traefik.http.routers." + hostname + ".rule"
	target[rule] = fmt.Sprintf("Host(`%s`)", host)
	if tp.Https {
		tlsLabel := fmt.Sprintf("traefik.http.routers.%s.tls", hostname)
		target[tlsLabel] = "true"
		tlsResolverLabel := fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", hostname)
		target[tlsResolverLabel] = tp.CertResolverName
	}
	return nil
}
