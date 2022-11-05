package app

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/space-fold-technologies/aurora-service/app/core/configuration"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/core/server"
)

type Application struct {
	server           *server.ServerCore
	configs          configuration.Configuration
	details          server.Details
	serviceResources *ServiceResources
	hashHandler      security.HashHandler
	tokenHandler     security.TokenHandler
}

func (a *Application) onStartUp() bool {
	// services to start the system with
	logging.GetInstance().Info("START UP INITIALIZATION")
	a.serviceResources.Initialize()
	return true
}

func (a *Application) onShutdown() bool {
	// services to shutdown and resources to clean up
	logging.GetInstance().Info("SHUTDOWN CALLED")
	return true
}

func (a *Application) Start() {
	a.details = server.FetchApplicationDetails()
	a.configs = configuration.ParseFromResource()
	configurationPath := filepath.Join(a.configs.ProfileDIR, "configurations")
	if _, err := os.Stat(configurationPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configurationPath, fs.ModeAppend); err != nil {
			logging.GetInstance().Error(err)
			os.Exit(-1)
		}
	}
	if privateKeyBytes, err := os.ReadFile(filepath.Join(a.configs.ProfileDIR, "keys", "private-rsa-key.pem")); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else if publicKeyBytes, err := os.ReadFile(filepath.Join(a.configs.ProfileDIR, "keys", "public-rsa-key.pem")); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else if privateKey, err := parseRSAPrivateKeyFromPEM(privateKeyBytes); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else if publicKey, err := parseRSAPublicKeyFromPEM(publicKeyBytes); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else {
		a.hashHandler = security.NewHashHandler(publicKey, privateKey)
		a.tokenHandler = security.NewTokenHandler(publicKey, privateKey)
		a.server = server.New(
			a.details,
			a.configs.Host,
			a.configs.Port,
			a.tokenHandler,
		)
	}

	a.serviceResources = ProduceServiceResources(
		a.server,
		a.configs,
		a.tokenHandler,
		a.hashHandler,
	)
	a.server.OnStartUp(a.onStartUp)
	a.server.OnShutDown(a.onShutdown)
	a.server.Start()
}

func (a *Application) Stop() {
	a.server.Stop()
}

// Parse PEM encoded PKCS1 or PKCS8 private key
func parseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, security.ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, security.ErrNotRSAPrivateKey
	}
	return pkey, nil
}

// Parse PEM encoded PKCS1 or PKCS8 public key
func parseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, security.ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, security.ErrNotRSAPublicKey
	}

	return pkey, nil
}
