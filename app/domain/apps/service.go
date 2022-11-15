package apps

import (
	"context"
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

func (as *AppService) Deploy(ws *websocket.Conn, properties *providers.TerminalProperties) error {
	if uri, err := as.repository.FetchImageURI(properties.Identifier); err != nil {
		return err
	} else if vars, err := as.repository.FetchEnvVars(properties.Name); err != nil {
		return err
	} else if app, err := as.repository.FetchDetails(properties.Name); err != nil {
		return err
	} else {
		order := &providers.Order{
			Name:       app.Name,
			Identifier: properties.Identifier,
			URI:        uri,
			Variables:  as.variables(vars),
			Volumes:    []providers.Mount{},
			Scale:      uint(app.Scale),
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
		return as.provider.Deploy(
			ws,
			properties,
			order,
			func(ctx context.Context, report *providers.Report) error {
				return as.processReport(ctx, properties.Identifier, properties.Name, report)
			})
	}
}

func (as *AppService) LogContainer(ws *websocket.Conn, properties *providers.TerminalProperties) error {
	logging.GetInstance().Infof("Logging Container: %s For App : %s", properties.Identifier, properties.Name)
	if err := as.provider.LogContainer(ws, properties, properties.Identifier); err != nil {
		return err

	}
	return nil
}

func (as *AppService) LogService(ws *websocket.Conn, properties *providers.TerminalProperties) error {
	if deployment, err := as.repository.Deployed(properties.Name); err != nil {
		return err
	} else {
		logging.GetInstance().Infof("Service Id: %s For App : %s", deployment.ServiceID, properties.Name)
		if err := as.provider.LogService(ws, properties, deployment.ServiceID); err != nil {
			return err
		}
	}
	return nil
}

func (as *AppService) Shell(ws *websocket.Conn, properties *providers.TerminalProperties) error {
	// Will need to read a lot more on this one
	if details, err := as.repository.FetchDetails(properties.Name); err != nil {
		return err
	} else if len(details.Instances) > 0 {
		container := details.Instances[0]
		logging.GetInstance().Infof("Shell Container: %s For App : %s", container.Identifier, properties.Name)
		as.provider.Shell(ws, properties, container.Identifier)
	}
	return nil
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

func (as *AppService) Rollback(ws *websocket.Conn, properties *providers.TerminalProperties) error {
	if summary, err := as.repository.FetchDeployment(properties.Identifier); err != nil {
		return err
	} else if vars, err := as.repository.FetchEnvVars(summary.Name); err != nil {
		return err
	} else if app, err := as.repository.FetchDetails(summary.Name); err != nil {
		return err
	} else {
		return as.provider.Deploy(
			ws,
			properties,
			&providers.Order{
				Identifier: properties.Identifier,
				URI:        summary.ImageURI,
				Variables:  as.variables(vars),
				Volumes:    []providers.Mount{},
				Scale:      uint(app.Scale),
			},
			func(ctx context.Context, report *providers.Report) error {
				return as.processReport(ctx, properties.Identifier, properties.Name, report)
			})
	}
}

func (as *AppService) processReport(ctx context.Context, identifier, name string, report *providers.Report) error {
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
	} else if as.provider.Nuke(deployment.ServiceID); err != nil {
		return err
	}
	return nil
}

func (as *AppService) variables(vars []*EnvVarEntry) []*providers.Variable {
	variables := make([]*providers.Variable, 0)
	for _, entry := range vars {
		logging.GetInstance().Infof("PACKER: KEY : %s VALUE : %s", entry.Key, entry.Val)
		variables = append(variables, &providers.Variable{Key: entry.Key, Value: entry.Val})
	}
	return variables
}
