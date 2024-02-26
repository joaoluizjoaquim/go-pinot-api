package container

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Pinot struct {
	Container testcontainers.Container
	TearDown  func()
	URI       string
}

func RunPinotContainer(ctx context.Context) (*Pinot, error) {

	absPath, err := filepath.Abs(filepath.Join(".", "testdata", "pinot-controller.conf"))
	if err != nil {
		return nil, fmt.Errorf("failed to add data: %s", err)
	}

	// newNetwork, err := network.New(ctx, network.WithCheckDuplicate())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create network: %s", err)
	// }

	// networkName := newNetwork.Name

	pinotZkContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			// Networks: []string{networkName},
			// NetworkAliases: map[string][]string{
			// 	networkName: {"pinot-zk"},
			// },
			Name:         "pinot-zk",
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2181/tcp"},
			Cmd:          []string{"StartZookeeper"},
			WaitingFor:   wait.ForLog("Start zookeeper at localhost:2181 in thread main").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	pinotZKMappedPort, err := pinotZkContainer.MappedPort(ctx, "2181")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %s", err)
	}

	pinotZKHost, err := pinotZkContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %s", err)
	}

	pinotZKURI := fmt.Sprintf("%s:%s", pinotZKHost, pinotZKMappedPort.Port())

	fmt.Println("ZK URI: ", pinotZKURI)

	defer pinotZkContainer.Terminate(ctx)

	pinotContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			// Networks: []string{networkName},
			// NetworkAliases: map[string][]string{
			// 	networkName: {"pinot-controller"}, // this does work btrw
			// },
			Name:         "pinot-controller",
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2123/tcp", "9000/tcp", "8000/tcp", "7050/tcp", "6000/tcp"},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/config/pinot-controller.conf",
					FileMode:          0o700,
				},
			},
			Cmd:        []string{"StartController", "-zkAddress", pinotZKURI}, //"StartController"  "-configFileName", "/config/pinot-controller.conf"
			WaitingFor: wait.ForLog("INFO: [HttpServer] Started.").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	defer pinotContainer.Terminate(ctx)

	tearDown := func() {
		if err := pinotContainer.Terminate(ctx); err != nil {
			log.Panicf("failed to terminate container: %s", err)
		}
	}

	pinotControllerMappedPort, err := pinotContainer.MappedPort(ctx, "9000")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %s", err)
	}

	pinotControllerHost, err := pinotContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %s", err)
	}

	pinotControllerIP := fmt.Sprintf("%s:%s", pinotControllerHost, pinotControllerMappedPort.Port())

	return &Pinot{
		Container: pinotContainer,
		TearDown:  tearDown,
		URI:       pinotControllerIP,
	}, nil
}

func (p *Pinot) CreateUser(ctx context.Context, userBytes []byte) (*model.UserActionResponse, error) {
	client := goPinotAPI.NewPinotAPIClient("http://" + p.URI)

	userCreationResponse, err := client.CreateUser(userBytes)
	if err != nil {
		log.Fatal(err)
	}

	return userCreationResponse, nil

}

func (p *Pinot) GetUsers(ctx context.Context) (*model.GetUsersResponse, error) {
	client := goPinotAPI.NewPinotAPIClient("http://" + p.URI)

	userResp, err := client.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	return userResp, nil
}
