package testutils

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/utils"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

type EigenDA struct {
	compose       tc.ComposeStack
	composeCancel context.CancelFunc
	testName      string
	rootPath      string
}

func NewEigenDA(rootPath, testName string) *EigenDA {
	return &EigenDA{
		rootPath: rootPath,
		testName: testName,
	}
}

func (e *EigenDA) MustStart() {
	var err error
	e.compose, err = tc.NewDockerCompose(
		filepath.Join(e.rootPath, "inabox/testdata", e.testName, "docker-compose.yml"),
	)
	if err != nil {
		panic(err)
	}

	var ctx context.Context
	ctx, e.composeCancel = context.WithCancel(context.Background())
	err = e.compose.Up(ctx, tc.Wait(true), tc.RemoveOrphans(true))
	if err != nil {
		container := strings.Split(err.Error(), " ")[1]
		utils.RunCommand("docker", "logs", container)
		panic(err)
	}
}

func (a *EigenDA) MustStop() {
	a.composeCancel()
	err := a.compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
	if err != nil {
		panic(err)
	}
}
