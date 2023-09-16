//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/mnabbasbaadi/distributedcache/master/service/tests/support/client"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	client *client.RegisterAPITestClient
}

func (s *E2ETestSuite) TestE2E() {
	s.T().Run("success", func(t *testing.T) {
		addr := "localhost:8080"
		ctx := context.Background()
		err := s.client.RegisterNode(ctx)
		s.NoError(err)
		err = s.client.UnRegisterNode(ctx, addr)
		s.NoError(err)
	})
}
