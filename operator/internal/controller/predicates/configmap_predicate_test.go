package predicates

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ = Describe("ConfigMap Predicates", func() {
	var (
		namespace = "configmap-namespace"
		name      = "configmap-name"
		predicate predicate.Predicate
		obj       *unstructured.Unstructured
	)

	BeforeEach(func() {
		predicate = NewConfigMapPredicate(namespace, name)

		obj = &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
		obj.SetName(name)
	})

	Context("when validating ConfigMap predicates", func() {
		It("should allow events for the correct ConfigMap", func() {
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}
			genericEvent := event.GenericEvent{Object: obj}
			deleteEvent := event.DeleteEvent{Object: obj}

			Expect(predicate.Create(createEvent)).To(BeTrue())
			Expect(predicate.Update(updateEvent)).To(BeTrue())
			Expect(predicate.Generic(genericEvent)).To(BeTrue())
			Expect(predicate.Delete(deleteEvent)).To(BeTrue())
		})

		It("should deny events for other ConfigMaps", func() {
			obj.SetName("other-name")
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}
			genericEvent := event.GenericEvent{Object: obj}
			deleteEvent := event.DeleteEvent{Object: obj}

			Expect(predicate.Create(createEvent)).To(BeFalse())
			Expect(predicate.Update(updateEvent)).To(BeFalse())
			Expect(predicate.Generic(genericEvent)).To(BeFalse())
			Expect(predicate.Delete(deleteEvent)).To(BeFalse())
		})

		It("should deny events from other namespaces", func() {
			obj.SetNamespace("other-namespace")
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}
			genericEvent := event.GenericEvent{Object: obj}
			deleteEvent := event.DeleteEvent{Object: obj}

			Expect(predicate.Create(createEvent)).To(BeFalse())
			Expect(predicate.Update(updateEvent)).To(BeFalse())
			Expect(predicate.Generic(genericEvent)).To(BeFalse())
			Expect(predicate.Delete(deleteEvent)).To(BeFalse())
		})

		It("should deny events from other namespaces and ConfigMaps", func() {
			obj.SetNamespace("other-namespace")
			obj.SetName("other-name")
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}
			genericEvent := event.GenericEvent{Object: obj}
			deleteEvent := event.DeleteEvent{Object: obj}

			Expect(predicate.Create(createEvent)).To(BeFalse())
			Expect(predicate.Update(updateEvent)).To(BeFalse())
			Expect(predicate.Generic(genericEvent)).To(BeFalse())
			Expect(predicate.Delete(deleteEvent)).To(BeFalse())
		})

		It("should deny events from other ConfigMaps in current namespace", func() {
			obj.SetName("other-name")
			createEvent := event.CreateEvent{Object: obj}
			updateEvent := event.UpdateEvent{ObjectNew: obj}
			genericEvent := event.GenericEvent{Object: obj}
			deleteEvent := event.DeleteEvent{Object: obj}

			Expect(predicate.Create(createEvent)).To(BeFalse())
			Expect(predicate.Update(updateEvent)).To(BeFalse())
			Expect(predicate.Generic(genericEvent)).To(BeFalse())
			Expect(predicate.Delete(deleteEvent)).To(BeFalse())
		})
	})
})
