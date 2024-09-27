package handlers

import (
	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("ReconcileAllHotNews", func() {
	var (
		ctx                context.Context
		mockClient         client.Client
		configMapNamespace string
		configMapName      string
		namespace          string
		configMap          *v1.ConfigMap
		hotNews1           *newsaggregatorv1.HotNews
		hotNews2           *newsaggregatorv1.HotNews
	)

	BeforeEach(func() {
		ctx = context.TODO()
		scheme := runtime.NewScheme()
		Expect(newsaggregatorv1.AddToScheme(scheme)).To(Succeed())
		Expect(v1.AddToScheme(scheme)).To(Succeed())
		mockClient = fake.NewClientBuilder().WithScheme(scheme).Build()

		configMapNamespace = "default"
		configMapName = "test-configmap"
		namespace = "test-namespace"

		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: configMapNamespace,
			},
		}

		hotNews1 = &newsaggregatorv1.HotNews{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hotnews-1",
				Namespace: namespace,
			},
			Spec: newsaggregatorv1.HotNewsSpec{
				FeedGroups: []string{"group1"},
			},
		}

		hotNews2 = &newsaggregatorv1.HotNews{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hotnews-2",
				Namespace: namespace,
			},
			Spec: newsaggregatorv1.HotNewsSpec{
				FeedGroups: nil,
			},
		}
	})

	Context("when the object is not a ConfigMap", func() {
		It("should return nil", func() {
			invalidObject := &v1.Pod{}
			requests := ReconcileAllHotNews(ctx, invalidObject, mockClient, namespace, configMapName, configMapNamespace)

			Expect(requests).To(BeNil())
		})
	})

	Context("when the ConfigMap does not match", func() {
		It("should return nil", func() {
			nonMatchingConfigMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "non-matching",
					Namespace: "wrong-namespace",
				},
			}

			requests := ReconcileAllHotNews(ctx, nonMatchingConfigMap, mockClient, namespace, configMapName, configMapNamespace)
			Expect(requests).To(BeNil())
		})
	})

	Context("when the ConfigMap matches", func() {
		BeforeEach(func() {
			Expect(mockClient.Create(ctx, hotNews1)).To(Succeed())
			Expect(mockClient.Create(ctx, hotNews2)).To(Succeed())
		})

		It("should enqueue HotNews resources that use the ConfigMap", func() {
			requests := ReconcileAllHotNews(ctx, configMap, mockClient, namespace, configMapName, configMapNamespace)

			Expect(requests).To(HaveLen(1))
			Expect(requests[0].NamespacedName.Name).To(Equal("hotnews-1"))
		})

		It("should not enqueue HotNews resources that don't use FeedGroups", func() {
			requests := ReconcileAllHotNews(ctx, configMap, mockClient, namespace, configMapName, configMapNamespace)

			for _, req := range requests {
				Expect(req.NamespacedName.Name).NotTo(Equal("hotnews-2"))
			}
		})
	})

	Context("when listing HotNews fails", func() {
		It("should log an error and return nil", func() {
			timeoutCtx, cancel := context.WithTimeout(ctx, 0)
			defer cancel()

			requests := ReconcileAllHotNews(timeoutCtx, configMap, mockClient, namespace, configMapName, configMapNamespace)

			Expect(requests).To(BeNil())
		})
	})
})
