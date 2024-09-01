package predicates

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

// TestConfigMapPredicates is the test suite for ConfigMap predicates
func TestFeedPredicates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feed Predicates Suite")
}
