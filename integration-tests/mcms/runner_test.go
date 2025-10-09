//go:build integration

package mcms

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMCMSandCCIPSuite(t *testing.T) {
	suite.Run(t, new(CCIPMCMSTestSuite))
}
