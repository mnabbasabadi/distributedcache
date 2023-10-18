//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/mnabbasbaadi/distributedcache/node/service/internal/node"
	"github.com/mnabbasbaadi/distributedcache/node/service/tests/support/client"
	"github.com/mnabbasbaadi/distributedcache/node/service/tests/support/receiver"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	node   node.Node
	client *client.CacheAPITestClient
	rec    *receiver.Receiver
}

func (s *E2ETestSuite) TestE2E() {
	s.T().Run("success", func(t *testing.T) {
		ctx := context.Background()
		key := "key"
		value := "value"
		err := s.client.SetValue(ctx, key, value)
		s.NoError(err)
		resp, _ := s.client.GetValue(ctx, key)
		s.Equal(key, resp.Key)
		s.Equal(value, resp.Value)

		val, b := s.node.Get([]byte(key))
		s.True(b)
		s.EqualValues(value, val)
	})
}
