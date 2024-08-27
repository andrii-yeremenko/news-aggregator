package controller_test

import (
	"bytes"
	v1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HotNews Controller", func() {

	var (
		reconcile  controller.HotNewsReconciler
		httpClient *MockHTTPClient
		fakeClient client.Client
		hotNews    *v1.HotNews
	)

	BeforeEach(func() {
		httpClient = new(MockHTTPClient)
		_ = v1.AddToScheme(scheme.Scheme)

		hotNews = &v1.HotNews{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-hotnews",
				Namespace: "default",
			},
			Spec: v1.HotNewsSpec{
				Keywords:      []string{"test", "news"},
				DateStart:     &metav1.Time{Time: metav1.Now().AddDate(0, 0, -1)},
				DateEnd:       &metav1.Time{Time: metav1.Now().AddDate(0, 0, 1)},
				SummaryConfig: v1.SummaryConfig{TitlesCount: 5},
			},
			Status: v1.HotNewsStatus{},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(hotNews).Build()

		reconcile = controller.HotNewsReconciler{
			Client:             fakeClient,
			HTTPClient:         httpClient,
			Scheme:             scheme.Scheme,
			NewsAggregatorURL:  "http://localhost:8080",
			ConfigMapName:      "test-configmap",
			ConfigMapNamespace: "default",
		}
	})

	AfterEach(func() {
	})

	Context("Test Successful Reconcile", func() {

		BeforeEach(func() {
			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"test-feed-group": "test-feed",
				},
			}

			err := fakeClient.Create(context.TODO(), &configMap)
			Expect(err).To(BeNil())
		})

		It("Should reconcile valid HotNews with only FeedGroups defined", func() {
			hotNews.Spec.FeedGroups = []string{"test-feed-group"}

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"test title\"}]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(1))
		})

		It("Should reconcile valid HotNews with only Feeds defined", func() {
			hotNews.Spec.Feeds = []string{"test-feed"}

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"test title\"}]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(1))
		})

		It("Should reconcile valid HotNews with both FeedGroups and Feeds defined", func() {
			hotNews.Spec.FeedGroups = []string{"test-feed-group"}
			hotNews.Spec.Feeds = []string{"test-feed"}

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"test title\"}, {\"title\": \"test title 2\"}]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(2))
			Expect(hotNews.Status.ArticlesTitles).To(ConsistOf("test title", "test title 2"))
		})

		It("Should reconcile valid HotNews with only FeedGroups defined and more articles than TitlesCount", func() {
			hotNews.Spec.FeedGroups = []string{"test-feed-group"}
			hotNews.Spec.SummaryConfig.TitlesCount = 1

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"!test title\"}, {\"title\": \"test title 2\"}]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(1))
			Expect(hotNews.Status.ArticlesTitles).To(ConsistOf("!test title"))
		})

		It("Should reconcile valid HotNews with only FeedGroups defined and less articles than TitlesCount", func() {
			hotNews.Spec.FeedGroups = []string{"test-feed-group"}
			hotNews.Spec.SummaryConfig.TitlesCount = 2

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"test title\"}]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(1))
			Expect(hotNews.Status.ArticlesTitles).To(ConsistOf("test title"))
		})

		It("Should reconcile valid HotNews with only FeedGroups defined and no articles", func() {
			hotNews.Spec.FeedGroups = []string{"test-feed-group"}

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(0))
		})

		It("Should reconcile valid HotNews and return articles in the correct alphabetical order", func() {
			hotNews.Spec.FeedGroups = []string{"test-feed-group"}
			hotNews.Spec.SummaryConfig.TitlesCount = 3

			err := fakeClient.Create(context.TODO(), hotNews)
			Expect(err).To(BeNil())

			httpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"a\"}, {\"title\": \"b\"}, {\"title\": \"c\"}]")),
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
			_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
			Expect(err).To(BeNil())

			Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			Expect(hotNews.Status.Conditions).To(HaveLen(1))
			Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionUpdated))
			Expect(hotNews.Status.ArticlesCount).To(Equal(3))
			Expect(hotNews.Status.ArticlesTitles).To(Equal([]string{"a", "b", "c"}))
		})
	})

	Context("Test Failed Reconcile", func() {

		Context("Test ConfigMap absent", func() {
			It("Should fail to reconcile when ConfigMap is absent and FeedGroups defined", func() {
				hotNews.Spec.FeedGroups = []string{"test-feed-group"}

				err := fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
				Expect(hotNews.Status.Conditions).To(HaveLen(1))
				Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionFailed))
			})

			It("Should fail to reconcile when ConfigMap is absent and only Feeds defined", func() {
				hotNews.Spec.Feeds = []string{"test-feed"}

				err := fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
				Expect(hotNews.Status.Conditions).To(HaveLen(1))
				Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionFailed))
			})
		})

		Context("Test Failed HTTP Request", func() {

			BeforeEach(func() {
				configMap := corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-configmap",
						Namespace: "default",
					},
					Data: map[string]string{
						"test-feed-group": "test-feed",
					},
				}

				err := fakeClient.Create(context.TODO(), &configMap)
				Expect(err).To(BeNil())
			})

			It("Should fail to reconcile when HTTP request returns non-200 status code", func() {
				hotNews.Spec.FeedGroups = []string{"test-feed-group"}

				err := fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				httpClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
				Expect(hotNews.Status.Conditions).To(HaveLen(1))
				Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionFailed))
			})

			It("Should fail to reconcile when HTTP request returns invalid JSON", func() {
				hotNews.Spec.FeedGroups = []string{"test-feed-group"}

				err := fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				httpClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
				}, nil)

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
				Expect(hotNews.Status.Conditions).To(HaveLen(1))
				Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionFailed))
			})
		})

		Context("Test Failed due to error in url", func() {
			It("Should fail to reconcile when URL is invalid", func() {
				hotNews.Spec.FeedGroups = []string{"test-feed-group"}

				err := fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				httpClient.On("Do", mock.Anything).Return(nil, &url.Error{})

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
				Expect(hotNews.Status.Conditions).To(HaveLen(1))
				Expect(hotNews.Status.Conditions[0]).To(Equal(v1.ConditionFailed))
			})
		})
	})

})

// MockHTTPClient is a mock for HTTPClient interface.
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	panic("unexpected call to Post")
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
