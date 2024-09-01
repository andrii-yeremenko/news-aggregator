package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func NewConfigMapPredicate(namespace, name string) predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.GetNamespace() == namespace && e.ObjectNew.GetName() == name
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return e.Object.GetNamespace() == namespace && e.Object.GetName() == name
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return e.Object.GetNamespace() == namespace && e.Object.GetName() == name
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return e.Object.GetNamespace() == namespace && e.Object.GetName() == name
		},
	}
}
