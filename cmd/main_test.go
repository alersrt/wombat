package main

import (
	"context"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test(t *testing.T) {
	compose, err := tc.NewDockerCompose("../docker/docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(
			t,
			compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal, tc.RemoveVolumes(true)),
			"compose.Down()",
		)
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
}
