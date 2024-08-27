package predicates

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"testing"
)

func TestFeedPredicates(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Feed Predicates Suite")
}

var _ = ginkgo.Describe("Feed Predicates", func() {
	var (
		namespace = "test-namespace"
		predicate predicate.Predicate
		obj       *unstructured.Unstructured
	)

	ginkgo.BeforeEach(func() {
		predicate = NewFeedPredicate(namespace)

		obj = &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
	})

	ginkgo.It("should allow events from the correct namespace", func() {
		obj.SetNamespace(namespace)
		createEvent := event.CreateEvent{Object: obj}
		updateEvent := event.UpdateEvent{ObjectNew: obj}

		gomega.Expect(predicate.Create(createEvent)).To(gomega.BeTrue())
		gomega.Expect(predicate.Update(updateEvent)).To(gomega.BeTrue())
	})

	ginkgo.It("should deny events from other namespaces", func() {
		obj.SetNamespace("other-namespace")
		createEvent := event.CreateEvent{Object: obj}
		updateEvent := event.UpdateEvent{ObjectNew: obj}

		gomega.Expect(predicate.Create(createEvent)).To(gomega.BeFalse())
		gomega.Expect(predicate.Update(updateEvent)).To(gomega.BeFalse())
	})
})
