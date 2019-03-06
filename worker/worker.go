package worker

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/jiajunhuang/huang/pkg/deployment"
	"github.com/jiajunhuang/huang/pkg/random"
	"go.uber.org/zap"
)

const (
	ContainerdAddr = "/run/containerd/containerd.sock"
	HuangNS        = "huang"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
	ctx       = namespaces.WithNamespace(context.Background(), HuangNS) // global namespace... a little bit bad taste
)

func Main() error {
	defer logger.Sync() // flushes buffer, if any

	client, err := containerd.New(ContainerdAddr)
	if err != nil {
		return err
	}
	defer client.Close()

	// create 2 redis container for test
	worker := Worker{}
	containers, err := worker.Deploy(client, &deployment.Deployment{Name: "redis-server"}, "docker.io/library/redis:alpine", 2)

	sugar.Infof("create containers with result: %+v, %s", containers, err)

	for _, c := range containers {
		c.Delete(ctx, containerd.WithSnapshotCleanup)
	}

	return nil
}

type Worker struct{}

func (w *Worker) Deploy(client *containerd.Client, deploy *deployment.Deployment, imagePath string, replica int) ([]containerd.Container, error) {
	image, err := client.Pull(ctx, imagePath, containerd.WithPullUnpack)
	if err != nil {
		return nil, err
	}

	containers := []containerd.Container{}

	for i := 0; i < replica; i++ {
		// NOTE: caller should hanle container.Delete(ctx, containerd.WithSnapshotCleanup)
		container, err := client.NewContainer(
			ctx,
			deploy.Name+random.UUID4String(),
			containerd.WithImage(image),
			containerd.WithNewSnapshot(random.UUID4String(), image),
			containerd.WithNewSpec(oci.WithImageConfig(image)),
		)
		if err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}

	return containers, nil
}
