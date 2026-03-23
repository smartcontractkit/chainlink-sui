//go:build integration

package mcms

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMCMSandCCIPSuite(t *testing.T) {
	// t.Run("CCIPMCMSTestSuite", func(t *testing.T) {
	// 	suite.Run(t, new(CCIPMCMSTestSuite))
	// })

	// t.Run("TokenPoolTestSuite", func(t *testing.T) {
	// 	suite.Run(t, new(TokenPoolTestSuite))
	// })

	// t.Run("MCMSUserTestSuite", func(t *testing.T) {
	// 	suite.Run(t, new(UpgradeTestSuite))
	// })

	t.Run("CCIPCurseMCMSSuite", func(t *testing.T) {
		suite.Run(t, new(CCIPCurseMCMSTestSuite))
	})
}
