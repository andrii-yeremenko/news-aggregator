package controller_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feed Controller")
}

func TestHotNewsController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HotNews Controller")
}
