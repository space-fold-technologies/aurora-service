package providers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
	"github.com/sirupsen/logrus"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"google.golang.org/protobuf/proto"
)

type ClientConfiguration struct {
	Retries         uint `yaml:"retries"`
	Jitter          uint `yaml:"maximum-jitter"`
	BackoffInterval uint `yaml:"backoff-interval"`
	Timeout         uint `yaml:"timeout"`
}

type retryClient struct {
	client *httpclient.Client
	logger *logrus.Logger
}

func NewClient(configuration *ClientConfiguration) AgentClient {
	instance := &retryClient{logger: logging.GetInstance()}
	return instance.initialize(configuration)
}

func (rc *retryClient) initialize(configuration *ClientConfiguration) *retryClient {
	backoffInterval := time.Duration(configuration.BackoffInterval) * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := time.Duration(configuration.Jitter) * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	timeout := time.Duration(configuration.Timeout) * time.Millisecond
	// Create a new client, sets the retry mechanism, and the number of times you would like to retry
	rc.client = httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(int(configuration.Retries)),
	)
	rc.client.AddPlugin(plugins.NewRequestLogger(nil, nil))
	return rc
}

func (rc *retryClient) Join(ctx context.Context, order *RegisterAgent, address, token string) error {
	uri := fmt.Sprintf("http://%s:2700/api/v1/aurora-agent/agents/join", address)
	return rc.put(ctx, uri, token, order)
}

func (rc *retryClient) Leave(ctx context.Context, order *RemoveAgent, address, token string) error {
	uri := fmt.Sprintf("http://%s:2700/api/v1/aurora-agent/agents/leave", address)
	return rc.put(ctx, uri, token, order)
}

func (rc *retryClient) Containers(ctx context.Context, serviceId, address string) (*ContainerReport, error) {
	uri := fmt.Sprintf("http://%s:2700/api/v1/aurora-agent/agents/%s/containers", address, serviceId)
	report := &ContainerReport{}
	if err := rc.getNoAuth(ctx, uri, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (rc *retryClient) put(ctx context.Context, path, token string, in proto.Message) error {
	if data, err := proto.Marshal(in); err != nil {

	} else if request, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewBuffer(data)); err != nil {
		return err
	} else if response, err := rc.authorizedClient(request, token); err != nil {
		return err
	} else if response.StatusCode != 204 {
		if _, err := io.ReadAll(response.Body); err != nil {
			return err
		}
	}
	return nil
}

// func (rc *retryClient) get(ctx context.Context, path, token string, out proto.Message) error {
// 	if request, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil); err != nil {
// 		return err
// 	} else if response, err := rc.authorizedClient(request, token); err != nil {
// 		return err
// 	} else if response.StatusCode != 204 {
// 		if _, err := io.ReadAll(response.Body); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (rc *retryClient) getNoAuth(ctx context.Context, path string, out proto.Message) error {
	if request, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil); err != nil {
		return err
	} else if response, err := rc.client.Do(request); err != nil {
		return err
	} else if response.StatusCode == 200 {
		if data, err := io.ReadAll(response.Body); err != nil {
			return err
		} else {
			return proto.Unmarshal(data, out)
		}
	}
	return nil
}

func (rc *retryClient) authorizedClient(request *http.Request, token string) (*http.Response, error) {
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return rc.client.Do(request)
}
