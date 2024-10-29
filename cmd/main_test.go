package main

import (
	"context"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"testing"
)

func Test(t *testing.T) {
	testCtx, testCancelFunc := context.WithCancel(context.Background())

	compose, err := tc.NewDockerCompose("../docker/docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(
			t,
			compose.Down(testCtx, tc.RemoveOrphans(true), tc.RemoveImagesLocal, tc.RemoveVolumes(true)),
			"compose.Down()",
		)
	})

	t.Cleanup(testCancelFunc)

	require.NoError(t, compose.Up(testCtx, tc.Wait(true)), "compose.Up()")

	main()
}
