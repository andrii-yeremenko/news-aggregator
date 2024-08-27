package v1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestFeedWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feed Validation")
}

func TestHotNewsWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HotNews Validation")
}
