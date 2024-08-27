package predicates

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"testing"
)

func TestConfigMapPredicates(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "ConfigMap Predicates Suite")
}

var _ = ginkgo.Describe("ConfigMap Predicates", func() {
	var (
		namespace = "configmap-namespace"
		name      = "configmap-name"
		predicate predicate.Predicate
		obj       *unstructured.Unstructured
	)

	ginkgo.BeforeEach(func() {
		predicate = NewConfigMapPredicate(namespace, name)

		obj = &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
		obj.SetName(name)
	})

	ginkgo.It("should allow events for the correct ConfigMap", func() {
		createEvent := event.CreateEvent{Object: obj}
		updateEvent := event.UpdateEvent{ObjectNew: obj}

		gomega.Expect(predicate.Create(createEvent)).To(gomega.BeTrue())
		gomega.Expect(predicate.Update(updateEvent)).To(gomega.BeTrue())
	})

	ginkgo.It("should deny events for other ConfigMaps", func() {
		obj.SetName("other-name")
		createEvent := event.CreateEvent{Object: obj}
		updateEvent := event.UpdateEvent{ObjectNew: obj}

		gomega.Expect(predicate.Create(createEvent)).To(gomega.BeFalse())
		gomega.Expect(predicate.Update(updateEvent)).To(gomega.BeFalse())
	})

	ginkgo.It("should deny events from other namespaces", func() {
		obj.SetNamespace("other-namespace")
		createEvent := event.CreateEvent{Object: obj}
		updateEvent := event.UpdateEvent{ObjectNew: obj}

		gomega.Expect(predicate.Create(createEvent)).To(gomega.BeFalse())
		gomega.Expect(predicate.Update(updateEvent)).To(gomega.BeFalse())
	})
})
