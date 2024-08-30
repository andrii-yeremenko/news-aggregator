package v1

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("ConfigMapWebhook", func() {
	var (
		ctx        context.Context
		configMap  *corev1.ConfigMap
		webhook    *ConfigMapWebhook
		fakeClient client.Client
	)

	BeforeEach(func() {
		ctx = context.TODO()
		_ = AddToScheme(scheme.Scheme)

		bbcFeed := &Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bbc-world",
				Namespace: "news-aggregator-namespace",
			},
		}
		abcFeed := &Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "abc-news",
				Namespace: "news-aggregator-namespace",
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

		webhook = &ConfigMapWebhook{
			Client:             fakeClient,
			ConfigMapName:      "hotnews-feeds-group",
			ConfigMapNamespace: "news-aggregator-namespace",
		}

		configMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hotnews-feeds-group",
				Namespace: "news-aggregator-namespace",
			},
			Data: map[string]string{
				"all":   "bbc-world,abc-news",
				"usa":   "abc-news",
				"world": "bbc-world",
			},
		}

		Expect(fakeClient.Create(ctx, bbcFeed)).To(Succeed())
		Expect(fakeClient.Create(ctx, abcFeed)).To(Succeed())
	})

	Context("When creating a ConfigMap", func() {
		It("should fail if any field is empty", func() {
			configMap.Data["empty-field"] = ""

			_, err := webhook.ValidateCreate(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("data field 'empty-field' is empty"))
		})

		It("should fail if feeds do not exist", func() {
			configMap.Data["non-existent"] = "invalid-feed"

			_, err := webhook.ValidateCreate(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("feeds \"invalid-feed\" do not exist"))
		})

		It("should succeed if all fields are non-empty and feeds exist", func() {
			_, err := webhook.ValidateCreate(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
