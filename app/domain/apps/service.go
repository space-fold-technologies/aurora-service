package apps

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AppService struct {
	repository ApplicationRepository
	hasher     security.HashHandler
	provider   providers.Provider
}

func NewService(provider providers.Provider, hasher security.HashHandler, repository ApplicationRepository) *AppService {
	instance := new(AppService)
	instance.provider = provider
	instance.hasher = hasher
	instance.repository = repository
	return instance
}

func (as *AppService) SetupDeployment(order *DeployAppOrder) (*DeploymentPass, error) {
	// Register  the deployment details
	// encrypt the credentials and sent them out
	pass := &DeploymentPass{}
	pass.Identifier = uuid.NewString()
	addedAt := time.Now().UTC()
	if err := as.repository.RegisterDeploymentEntry(&DeploymentOrder{
		Identifier:      pass.Identifier,
		ApplicationName: order.GetName(),
		Status:          "INITIATED",
		ImageURI:        order.GetImageUri(),
		CreatedAt:       &addedAt,
	}); err != nil {
		return nil, err
	}
	if order.GetCredentials() != nil {
		if content, err := proto.Marshal(order.GetCredentials()); err != nil {
			return nil, err
		} else if hash, err := as.hasher.Encrypt(content); err != nil {
			return nil, err
		} else {
			pass.Token = hash
		}
	}
	go func(name string) {
		if err := as.reset(name); err != nil {
			logging.GetInstance().Error(err)
		}
	}(order.GetName())
	return pass, nil
}

func (as *AppService) Deploy(ws *websocket.Conn, properties *DeploymentProperties) error {
	if uri, err := as.repository.FetchImageURI(properties.Identifier); err != nil {
		return err
	} else if vars, err := as.repository.FetchEnvVars(properties.Name); err != nil {
		return err
	} else if app, err := as.repository.FetchDetails(properties.Name); err != nil {
		return err
	} else {
		order := &providers.DeploymentOrder{
			Name:                 app.Name,
			ID:                   properties.Identifier,
			ImageURI:             uri,
			EnvironmentVariables: as.variables(vars),
			Volumes:              []providers.Mount{},
			Scale:                uint(app.Scale),
		}
		if len(properties.Token) > 0 {
			if data, err := base64.URLEncoding.DecodeString(properties.Token); err != nil {
				return err
			} else {
				credentials := &RegistryCredentials{}
				if content, err := as.hasher.Decrypt(data); err != nil {
					return err
				} else if err := proto.Unmarshal(content, credentials); err != nil {
					return err
				} else {
					order.Username = credentials.GetUsername()
					order.Password = credentials.GetPassword()
					order.Temporary = credentials.GetTemporary()
				}
			}
		}
		reporter := providers.Reporter{
			Done: func(report providers.Report) error {
				return as.processReport(properties.Identifier, properties.Name, &report)
			},
			Progress: func(line []byte) { ws.WriteMessage(websocket.TextMessage, line) },
		}
		return as.provider.DeployService(ws, order, &reporter)
	}
}

func (as *AppService) Log(ws *websocket.Conn, properties *LogProperties) error {
	if deployment, err := as.repository.Deployed(properties.Name); err != nil {
		return err
	} else {
		return as.provider.Log(ws, &providers.LogProperties{
			ServiceID:      deployment.ServiceID,
			Details:        properties.ShowDetails,
			ShowTimestamps: properties.ShowTimeStamps})
	}
}

func (as *AppService) Shell(ws *websocket.Conn, properties *ShellProperties) error {
	// Will need to read a lot more on this one
	if details, err := as.repository.FetchContainer(properties.Identifier); err != nil {
		return err
	} else {
		logging.GetInstance().Infof("Shell Container: %s For App : %s", details.ID, properties.Name)
		return as.provider.Shell(ws, &providers.ShellProperties{
			ContainerID: details.ID,
			Host:        details.NodeIP,
			Width:       properties.Width,
			Heigth:      properties.Height,
			Terminal:    properties.Term,
		})
	}
}

func (as *AppService) Create(order *CreateAppOrder) error {

	return as.repository.Create(&ApplicationEntry{
		Identifier:  uuid.NewString(),
		Name:        order.GetName(),
		Description: order.GetDescription(),
		Team:        order.GetTeam(),
		Cluster:     order.GetCluster(),
		Scale:       uint32(order.GetScale()),
	})
}

func (as *AppService) Update(order *UpdateAppOrder) error {
	return as.repository.Update(
		order.GetName(),
		order.GetDescription(),
		order.GetScale(),
	)
}

func (as *AppService) List(cluster string) (*AppSummary, error) {
	if entries, err := as.repository.List(cluster); err != nil {
		return nil, err
	} else {
		summary := &AppSummary{Entries: make([]*AppSummary_Entry, 0)}
		for _, entry := range entries {
			summary.Entries = append(summary.Entries, &AppSummary_Entry{
				Name:  entry.Name,
				Scale: int32(entry.Scale),
			})
		}
		return summary, nil
	}
}

func (as *AppService) Information(name string) (*AppDetails, error) {
	if application, err := as.repository.FetchDetails(name); err != nil {
		return nil, err
	} else {
		details := &AppDetails{
			Name:        application.Name,
			Description: application.Description,
			Team:        application.Team.Name,
			Cluster:     application.Cluster.Name,
			Containers:  make([]*AppDetails_Container, 0),
		}
		for _, instance := range application.Instances {
			details.Containers = append(details.Containers, &AppDetails_Container{
				Identifier: instance.Identifier,
				Ip:         instance.IP,
				Family:     instance.Family,
			})
		}
		return details, nil
	}
}

func (as *AppService) Deployments(name string) (*Deployments, error) {
	if deployments, err := as.repository.Deployments(name); err != nil {
		return nil, err
	} else {
		pack := &Deployments{Entries: make([]*Deployments_Entry, 0)}
		for _, deployment := range deployments {
			pack.Entries = append(pack.Entries, &Deployments_Entry{
				Identifier: deployment.Identifier,
				Image:      deployment.ImageURI,
				Status:     deployment.Status,
				Report:     deployment.Report,
				Stamp:      timestamppb.New(*deployment.CompletedAt),
			})
		}
		return pack, nil
	}
}

func (as *AppService) Rollback(ws *websocket.Conn, properties *RollbackProperties) error {
	if summary, err := as.repository.FetchDeployment(properties.Identifier); err != nil {
		return err
	} else if vars, err := as.repository.FetchEnvVars(summary.Name); err != nil {
		return err
	} else if app, err := as.repository.FetchDetails(summary.Name); err != nil {
		return err
	} else {
		reporter := providers.Reporter{
			Done: func(report providers.Report) error {
				return as.processReport(properties.Identifier, properties.Name, &report)
			},
			Progress: func(line []byte) { ws.WriteMessage(websocket.TextMessage, line) },
		}
		return as.provider.DeployService(
			ws,
			&providers.DeploymentOrder{
				ID:                   properties.Identifier,
				ImageURI:             summary.ImageURI,
				EnvironmentVariables: as.variables(vars),
				Volumes:              []providers.Mount{},
				Scale:                uint(app.Scale),
			}, &reporter)
	}
}

func (as *AppService) UpdateContainers() {
	logger := logging.GetInstance()
	if checks, err := as.repository.FetchActiveDeployments(); err != nil {
		logger.Errorf("checks failure : %s", err.Error())
	} else {
		identifiers := make([]string, 0)
		for _, check := range checks {
			identifiers = append(identifiers, check.ID)
		}
		if len(identifiers) == 0 {
			return
		}
		if err := as.provider.FetchContainers(identifiers, func(state map[string][]*providers.Instance) {
			for id, instances := range state {
				logger.Infof("updating containers for service id[:%s]", id)
				containers := make([]*ContainerOrder, 0)
				for _, instance := range instances {
					containers = append(containers, &ContainerOrder{
						Identifier: instance.ID,
						IP:         instance.IP,
						Family:     instance.Family,
						ServiceID:  instance.ServiceID,
						Node:       instance.Node,
					})
				}
				if err := as.repository.AddContainers(containers); err != nil {
					logger.Errorf("failed to update containers: Err::[%s]", err.Error())
				}
			}
		}); err != nil {
			logger.Errorf("checks failure : %s", err.Error())
		}
		targetTime := time.Now().Add(-10 * time.Minute)
		if err := as.repository.RemoveContainersOlderThan(&targetTime); err != nil {
			logger.Errorf("failed to remove dead container details : %s", err.Error())
		}
	}
}

func (as *AppService) processReport(identifier, name string, report *providers.Report) error {
	completedAt := time.Now().UTC()
	if err := as.repository.UpdateDeploymentEntry(&DeploymentUpdate{
		Identifier: identifier,
		Status:     report.Status,
		Report:     report.Message,
		ServiceID:  report.ServiceID,
		ImageURI:   report.ImageDigest,
		UpdatedAt:  &completedAt,
	}); err != nil {
		return err
	} else if report.Status == "DEPLOYED" {
		containers := make([]*ContainerOrder, 0)
		for _, instance := range report.Instances {
			containers = append(containers, &ContainerOrder{
				Identifier: instance.ID,
				IP:         instance.IP,
				Family:     instance.Family,
				ServiceID:  instance.ServiceID,
				Node:       instance.Node,
			})
			logging.GetInstance().Infof(
				"ADDED CONTAINER WITH SERVICE ID: %s TASK ID: %s CONTAINER ID: %s ON NODE ID: %s",
				instance.ServiceID, instance.TaskID, instance.ID, instance.Node,
			)
		}
		return as.repository.AddContainers(containers)
	}
	return nil
}

func (as *AppService) reset(name string) error {
	if deployment, err := as.repository.Deployed(name); err != nil {
		return err
	} else if err = as.repository.RemoveContainers(deployment.ApplicationID); err != nil {
		return err
	} else if as.provider.Stop(deployment.ServiceID); err != nil {
		return err
	}
	return nil
}

func (as *AppService) variables(vars []*EnvVarEntry) map[string]string {
	variables := make(map[string]string)
	for _, entry := range vars {
		variables[entry.Key] = entry.Val
	}
	return variables
}
