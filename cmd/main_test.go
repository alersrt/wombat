package main

import (
	"context"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"testing"
)

func Test(t *testing.T) {
	mainCtx, mainCancelCauseFunc = context.WithCancelCause(context.Background())

	compose, err := tc.NewDockerCompose("../docker/docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(
			t,
			compose.Down(mainCtx, tc.RemoveOrphans(true), tc.RemoveImagesLocal, tc.RemoveVolumes(true)),
			"compose.Down()",
		)
	})

	_, cancelFunc := context.WithCancel(mainCtx)
	t.Cleanup(cancelFunc)

	require.NoError(t, compose.Up(mainCtx, tc.Wait(true)), "compose.Up()")

	main()
}
