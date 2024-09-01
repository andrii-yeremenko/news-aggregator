package v1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Feed Validation", func() {
	var (
		scheme *runtime.Scheme
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		_ = AddToScheme(scheme)
		k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()
	})

	Describe("ValidateFeedName", func() {
		Context("when validating feed names", func() {

			It("should validate a valid name", func() {
				feed := &Feed{
					Spec: FeedSpec{
						Name: "valid-name",
						Link: "https://example.com",
					},
				}
				err := validateFeed(feed)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not validate an empty name", func() {
				feed := &Feed{
					Spec: FeedSpec{Name: ""},
				}
				err := validateFeed(feed)
				Expect(err).To(HaveOccurred())
			})

			It("should not validate a name that is too long", func() {
				feed := &Feed{
					Spec: FeedSpec{Name: "this-name-is-way-too-long-to-be-valid"},
				}
				err := validateFeed(feed)
				Expect(err).To(HaveOccurred())
			})

			It("should not validate a name with invalid characters", func() {
				feed := &Feed{
					Spec: FeedSpec{Name: "invalid!name"},
				}
				err := validateFeed(feed)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ValidateFeedLink", func() {
		Context("when validating feed links", func() {

			It("should validate a valid link", func() {
				feed := &Feed{
					Spec: FeedSpec{
						Link: "http://example.com",
						Name: "valid-link",
					},
				}
				err := validateFeed(feed)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not validate an invalid link", func() {
				feed := &Feed{
					Spec: FeedSpec{
						Link: "invalid-link",
						Name: "invalid-link",
					},
				}
				err := validateFeed(feed)
				Expect(err).To(HaveOccurred())
			})

			It("should not validate an empty link", func() {
				feed := &Feed{
					Spec: FeedSpec{
						Link: "",
						Name: "empty-link",
					},
				}
				err := validateFeed(feed)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("CheckNameUniqueness", func() {
		Context("when checking name uniqueness", func() {
			var existingFeed *Feed
			var existingFeedList *FeedList

			BeforeEach(func() {
				existingFeed = &Feed{
					Spec: FeedSpec{
						Name: "existing-feed",
						Link: "https://example.com",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						UID:       "existing-uid",
					},
				}
				existingFeedList = &FeedList{
					Items: []Feed{*existingFeed},
				}
				k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(existingFeedList).Build()
			})

			It("should validate a unique name", func() {
				feed := &Feed{
					Spec: FeedSpec{
						Name: "new-feed",
						Link: "https://example.com",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: existingFeed.Namespace,
						UID:       "new-uid",
					},
				}
				err := checkNameUniqueness(feed)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not validate a duplicate name", func() {
				feed := &Feed{
					Spec: FeedSpec{
						Name: "existing-feed",
						Link: "https://example.com",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: existingFeed.Namespace,
						UID:       "new-uid",
					},
				}
				err := checkNameUniqueness(feed)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
