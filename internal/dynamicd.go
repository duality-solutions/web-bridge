package webbridge

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

/*
docker volume create --name=dynamicd-data
docker run -e DISABLEWALLET=0 -v dynamicd-data:/dynamic --name=dynamicd -d -p 33300:33300 -p 127.0.0.1:33350:33350 dualitysolutions/docker-dynamicd
*/

// LoadDockerDynamicd is used to create and run a Docker dynamicd full node and cli.
func LoadDockerDynamicd() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "dynamicd",
		ExposedPorts: nat.PortSet{"33350": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("33350"): {{HostIP: "127.0.0.1", HostPort: "8080"}}},
	}, nil, "dynamicd")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return cli, nil
}

// LoadRPCDynamicd is used to create and run a managed dynamicd full node and cli.
func LoadRPCDynamicd() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "dynamicd",
		ExposedPorts: nat.PortSet{"33350": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("33350"): {{HostIP: "127.0.0.1", HostPort: "8080"}}},
	}, nil, "dynamicd")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return cli, nil
}
