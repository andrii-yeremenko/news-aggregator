package predicates

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"testing"
)

// TestConfigMapPredicates is the test suite for ConfigMap predicates
func TestFeedPredicates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feed Predicates Suite")
}

var _ = Describe("Feed Predicates", func() {
	var (
		namespace = "test-namespace"
		predicate predicate.Predicate
		obj       *unstructured.Unstructured
	)

	BeforeEach(func() {
		predicate = NewFeedPredicate(namespace)

		obj = &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
	})

	Context("when validating Feed predicates", func() {
		It("should allow events from the correct namespace", func() {
			obj.SetNamespace(namespace)
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}

			Expect(predicate.Create(createEvent)).To(BeTrue())
			Expect(predicate.Update(updateEvent)).To(BeTrue())
		})

		It("should deny events from other namespaces", func() {
			obj.SetNamespace("other-namespace")
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}

			Expect(predicate.Create(createEvent)).To(BeFalse())
			Expect(predicate.Update(updateEvent)).To(BeFalse())
		})
	})
})
