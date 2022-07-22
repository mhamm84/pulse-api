//go:build integration
// +build integration

package economic

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// SETUP

	os.Exit(m.Run())

	// CLEANUP
}
