package controller_test

import (
	"bytes"
	v1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	"com.teamdev/news-aggregator/internal/controller/mocks"
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		mockCtrl   *gomock.Controller
		httpClient *mocks.MockHTTPClient
		fakeClient client.Client
		hotNews    *v1.HotNews
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		httpClient = mocks.NewMockHTTPClient(mockCtrl)
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
			Finalizer:          "finalizer.news-aggregator.com.teamdev",
			NewsAggregatorURL:  "http://localhost:8080",
			ConfigMapName:      "test-configmap",
			ConfigMapNamespace: "default",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

			httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

		Context("Test Owner References", func() {
			It("Should set Owner References for the Feed", func() {
				feed := &v1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-feed",
						Namespace: "default",
					},
					Spec: v1.FeedSpec{
						Name: "test-feed",
						Link: "http://localhost:8080",
					},
				}

				err := fakeClient.Create(context.TODO(), feed)
				Expect(err).To(BeNil())

				hotNews.Spec.Feeds = []string{"test-feed"}

				err = fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"test title\"}]")),
				}, nil)

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(BeNil())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())

				Expect(fakeClient.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "test-feed"}, feed)).To(Succeed())
				Expect(feed.OwnerReferences).To(HaveLen(1))
				Expect(feed.OwnerReferences[0].Name).To(Equal(hotNews.Name))
			})

			It("Should remove Owner References for the Feed when HotNews is deleting", func() {
				feed := &v1.Feed{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-feed",
						Namespace: "default",
					},
					Spec: v1.FeedSpec{
						Name: "test-feed",
						Link: "http://localhost:8080",
					},
				}

				err := fakeClient.Create(context.TODO(), feed)
				Expect(err).To(BeNil())

				hotNews.Spec.Feeds = []string{"test-feed"}

				err = fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("[{\"title\": \"test title\"}]")),
				}, nil)

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(BeNil())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())

				Expect(fakeClient.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "test-feed"}, feed)).To(Succeed())
				Expect(feed.OwnerReferences).To(HaveLen(1))
				Expect(feed.OwnerReferences[0].Name).To(Equal(hotNews.Name))

				err = fakeClient.Delete(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				namespacedName = types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(BeNil())

				Expect(fakeClient.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "test-feed"}, feed)).To(Succeed())
				Expect(feed.OwnerReferences).To(HaveLen(0))
			})
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
			})

			It("Should fail to reconcile when ConfigMap is absent and only Feeds defined", func() {
				hotNews.Spec.Feeds = []string{"test-feed"}

				err := fakeClient.Create(context.TODO(), hotNews)
				Expect(err).To(BeNil())

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
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

				httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

				httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
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

				httpClient.EXPECT().Do(gomock.Any()).Times(0).Return(nil, &url.Error{})

				namespacedName := types.NamespacedName{Namespace: "default", Name: "test-hotnews"}
				_, err = reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: namespacedName})
				Expect(err).To(HaveOccurred())

				Expect(fakeClient.Get(context.TODO(), namespacedName, hotNews)).To(Succeed())
			})
		})
	})

})
