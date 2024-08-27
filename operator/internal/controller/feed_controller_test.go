package controller_test

import (
	"bytes"
	"com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	"com.teamdev/news-aggregator/internal/controller/mocks"
	"context"
	"errors"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	serviceURL      = "http://localhost:8080"
	serviceEndpoint = "/sources"
	finalizer       = "finalizer.news-aggregator.com.teamdev"
)

var _ = Describe("Feed Controller", func() {

	var (
		reconciler *controller.FeedReconcile
		httpClient *mocks.MockHTTPClient
		fakeClient client.Client
		mockCtrl   *gomock.Controller
		ctx        context.Context
		feed       *v1.Feed
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		httpClient = mocks.NewMockHTTPClient(mockCtrl)
		scheme := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(scheme)
		_ = v1.AddToScheme(scheme)

		feed = &v1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
			Spec: v1.FeedSpec{
				Name: "test-feed",
				Link: "http://example.com/rss",
			},
			Status: v1.FeedStatus{
				Conditions: []v1.Condition{},
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(feed).Build()
		reconciler = &controller.FeedReconcile{
			Client:     fakeClient,
			Scheme:     scheme,
			HTTPClient: httpClient,
			ServiceURL: serviceURL,
			Finalizer:  finalizer,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Test Successful Reconciliation", func() {

		It("Should reconcile successfully when Feed is newly created", func() {
			Expect(fakeClient.Create(context.TODO(), feed)).To(Succeed())

			httpClient.EXPECT().
				Post(serviceURL+serviceEndpoint, "application/json", gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(BeNil(), "Expected no error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected successful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(1), "Expected one condition to be added")
			Expect(feed.Status.Conditions[0].Type).To(Equal(v1.ConditionAdded), "Expected ConditionAdded status")
		})

		It("Should reconcile successfully when Feed is updated", func() {
			feed.Status.Conditions = append(feed.Status.Conditions, v1.Condition{Type: v1.ConditionAdded})
			feed.CreationTimestamp = metav1.Now()
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())

			httpClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(BeNil(), "Expected no error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected successful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(2), "Expected second condition to be added")
			Expect(feed.Status.Conditions[1].Type).To(Equal(v1.ConditionUpdated), "Expected ConditionUpdated status")
		})

		It("Should reconcile successfully when Feed is deleted", func() {

			feed.Status.Conditions = append(feed.Status.Conditions, v1.Condition{Type: v1.ConditionAdded})
			feed.CreationTimestamp = metav1.Now()
			feed.Finalizers = append(feed.Finalizers, finalizer, "test-finalizer")
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())

			httpClient.EXPECT().
				Do(gomock.Any()).Times(1).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			Expect(fakeClient.Delete(ctx, feed)).To(Succeed())

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(BeNil(), "Expected no error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected successful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(2), "Expected second condition to be added")
			Expect(feed.Status.Conditions[1].Type).To(Equal(v1.ConditionDeleted), "Expected ConditionDeleted status")
		})
	})

	Context("Test Error Reconciliation", func() {

		It("Should fail to reconcile when Feed creation fails due to HTTP error", func() {
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())

			httpClient.EXPECT().
				Post(serviceURL+serviceEndpoint, "application/json", gomock.Any()).
				Return(nil, errors.New("network error"))

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(HaveOccurred(), "Expected error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected unsuccessful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(1), "Expected one condition to be added")
			Expect(feed.Status.Conditions[0].Type).To(Equal(v1.ConditionFailed), "Expected ConditionFailed status")
		})

		It("Should fail to reconcile when Feed is in unknown state", func() {
			feed.Status.Conditions = append(feed.Status.Conditions, v1.Condition{Type: v1.ConditionFailed})
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(HaveOccurred(), "Expected error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected unsuccessful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(1), "Expected no new conditions to be added")
		})

		It("Should fail to reconcile when json response is invalid", func() {
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())

			httpClient.EXPECT().
				Post(serviceURL+serviceEndpoint, "application/json", gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
				}, nil)

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(HaveOccurred(), "Expected error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected unsuccessful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(1), "Expected one condition to be added")
			Expect(feed.Status.Conditions[0].Type).To(Equal(v1.ConditionFailed), "Expected ConditionFailed status")
		})

		It("Should fail to reconcile when invalid response status code is returned during creation", func() {
			Expect(fakeClient.Create(ctx, feed)).To(Succeed())

			httpClient.EXPECT().
				Post(serviceURL+serviceEndpoint, "application/json", gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(HaveOccurred(), "Expected error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected unsuccessful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())
			Expect(feed.Status.Conditions).To(HaveLen(1), "Expected one condition to be added")
			Expect(feed.Status.Conditions[0].Type).To(Equal(v1.ConditionFailed), "Expected ConditionFailed status")
		})
	})

	Context("Test adding Finalizer", func() {

		It("Should add finalizer to Feed when it is created", func() {
			Expect(fakeClient.Create(context.TODO(), feed)).To(Succeed())

			httpClient.EXPECT().
				Post(serviceURL+serviceEndpoint, "application/json", gomock.Any()).
				Return(&http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)

			result, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(feed)})

			Expect(err).To(BeNil(), "Expected no error during reconciliation")
			Expect(result).To(Equal(ctrl.Result{}), "Expected successful reconciliation result")

			Expect(fakeClient.Get(ctx, client.ObjectKeyFromObject(feed), feed)).To(Succeed())

			Expect(feed.ObjectMeta.Finalizers).To(ContainElement(finalizer), "Expected finalizer to be added")
		})
	})
})
