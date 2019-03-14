package worker

import (
	"context"
	"sync"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/jiajunhuang/huang/pkg/deployment"
	"github.com/jiajunhuang/huang/pkg/random"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	ContainerdAddr  = "/run/containerd/containerd.sock"
	HuangNS         = "huang"
	DeploymentLabel = "deploy"
	SyncDuration    = 3
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func Main(c *cli.Context) error {
	defer logger.Sync() // flushes buffer, if any

	client, err := containerd.New(ContainerdAddr)
	if err != nil {
		return errors.Wrap(err, "cannot init containerd")
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), HuangNS)

	deploy := deployment.Deployment{Name: "redis-server"}

	// create 2 redis container for test
	worker := Worker{client: client, deploymentMap: DeploymentMap{}}
	go worker.SyncLocalContainerInfos(ctx)
	containers, err := worker.Deploy(ctx, &deploy, "docker.io/library/redis:alpine", 2)

	sugar.Infof("create containers with result: %+v, %s", containers, err)

	time.Sleep(time.Second * time.Duration(10))

	worker.Delete(ctx, &deploy)

	return nil
}

type DeploymentMap map[string][]containerd.Container

type Worker struct {
	sync.RWMutex

	deploymentMap DeploymentMap
	client        *containerd.Client
}

// SyncLocalContainerInfos sync local container infomation
func (w *Worker) SyncLocalContainerInfos(ctx context.Context) {
	for {
		func() {
			w.Lock()
			defer w.Unlock()

			containers, err := w.client.Containers(ctx)
			if err != nil {
				sugar.Errorf("failed to get containers: %+v", err)
				return
			}

			deploymentMap := DeploymentMap{}
			for _, c := range containers {
				labels, err := c.Labels(ctx)
				if err != nil {
					sugar.Errorf("failed to get labels of container %+v: %s", c, err)
					continue
				}

				deploymentName := labels[DeploymentLabel]
				if deploymentName == "" {
					sugar.Errorf("container %+v do not set deployment label", c)
					continue
				}

				cl := deploymentMap[deploymentName]
				cl = append(cl, c)
				deploymentMap[deploymentName] = cl
			}

			w.deploymentMap = deploymentMap
			sugar.Infof("we have those deployments with containers for now: %+v", w.deploymentMap)
		}()

		time.Sleep(time.Second * time.Duration(SyncDuration))
	}
}

// Deploy deploy a deployment
func (w *Worker) Deploy(ctx context.Context, deploy *deployment.Deployment, imagePath string, replica int) ([]containerd.Container, error) {
	image, err := w.client.Pull(ctx, imagePath, containerd.WithPullUnpack)
	if err != nil {
		return nil, errors.Wrap(err, "connot pull image")
	}

	containers := []containerd.Container{}

	for i := 0; i < replica; i++ {
		// NOTE: caller should hanle container.Delete(ctx, containerd.WithSnapshotCleanup)
		container, err := w.client.NewContainer(
			ctx,
			deploy.Name+random.UUID4String(),
			containerd.WithImage(image),
			containerd.WithNewSnapshot(random.UUID4String(), image),
			containerd.WithNewSpec(oci.WithImageConfig(image)),
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create container")
		}
		container.SetLabels(ctx, map[string]string{DeploymentLabel: deploy.Name})
		sugar.Infof("create container %s in deployment %s succeed", container.ID(), deploy.Name)
		containers = append(containers, container)
	}

	return containers, nil
}

func (w *Worker) Delete(ctx context.Context, deploy *deployment.Deployment) error {
	w.Lock()
	defer w.Unlock()

	cl, ok := w.deploymentMap[deploy.Name]
	if !ok {
		return errors.Wrap(nil, "cannot find deploy "+deploy.Name)
	}

	for _, c := range cl {
		if err := c.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
			return errors.Wrap(err, "cannot delete container "+c.ID())
		}
		sugar.Warnf("delete container %s in deployment %s", c.ID(), deploy.Name)
	}

	w.deploymentMap[deploy.Name] = cl[:0]
	return nil
}
