package v1

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HotNews Resource Validation", func() {
	var hotNews HotNews

	BeforeEach(func() {
		k8sClient = fake.NewClientBuilder().WithObjects(&v1.ConfigMap{
			Data: map[string]string{
				"group1": "feed1",
				"group2": "feed2",
			},
		}).Build()
	})

	Describe("Spec Validation", func() {
		Context("when validating Spec fields", func() {
			It("should pass with valid Spec", func() {
				hotNews = HotNews{
					Spec: HotNewsSpec{
						Keywords:  []string{"news", "update"},
						DateStart: &metav1.Time{Time: time.Now()},
						DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
					},
				}
				err := hotNews.validateHotNews()
				Expect(err).To(HaveLen(0))
			})

			It("should fail when Keywords are empty", func() {
				hotNews = HotNews{
					Spec: HotNewsSpec{
						Keywords:  []string{},
						DateStart: &metav1.Time{Time: time.Now()},
						DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
					},
				}
				err := hotNews.validateHotNews()
				Expect(err).To(HaveLen(1))
				Expect(err[0].Type.String()).To(Equal("Required value"))
				Expect(err[0].Field).To(Equal("spec.keywords"))
				Expect(err[0].Detail).To(Equal("keywords must be provided"))
			})

			It("should fail when DateEnd is before DateStart", func() {
				hotNews = HotNews{
					Spec: HotNewsSpec{
						Keywords:  []string{"news"},
						DateStart: &metav1.Time{Time: time.Now()},
						DateEnd:   &metav1.Time{Time: time.Now().Add(-24 * time.Hour)},
					},
				}
				err := hotNews.validateHotNews()
				Expect(err).To(HaveLen(1))
				Expect(err[0].Type.String()).To(Equal("Invalid value"))
				Expect(err[0].Field).To(Equal("spec.dateEnd"))
				Expect(err[0].Detail).To(Equal("dateEnd must be after dateStart"))
			})

			It("should fail when DateStart is set without DateEnd", func() {
				hotNews = HotNews{
					Spec: HotNewsSpec{
						Keywords:  []string{"news"},
						DateStart: &metav1.Time{Time: time.Now()},
					},
				}
				err := hotNews.validateHotNews()
				Expect(err).To(HaveLen(1))
				Expect(err[0].Type.String()).To(Equal("Required value"))
				Expect(err[0].Field).To(Equal("spec.dateEnd"))
				Expect(err[0].Detail).To(Equal("dateEnd must be provided"))
			})
		})
	})

	Describe("FeedGroups Validation", func() {
		Context("when validating FeedGroups", func() {
			It("should pass with existing FeedGroups", func() {
				hotNews = HotNews{
					Spec: HotNewsSpec{
						FeedGroups: []string{"group1", "group2"},
					},
				}
				err := hotNews.validateFeedGroups()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should fail with non-existent FeedGroups", func() {
				hotNews = HotNews{
					Spec: HotNewsSpec{
						FeedGroups: []string{"group1", "group3"},
					},
				}
				err := hotNews.validateFeedGroups()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("feedGroup group3 does not exist in ConfigMap /"))
			})
		})
	})

	Describe("CRUD Operations Validation", func() {
		BeforeEach(func() {
			hotNews = HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news"},
					DateStart: &metav1.Time{Time: time.Now()},
					DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
				},
			}
		})

		Context("when creating HotNews", func() {
			It("should pass the creation validation", func() {
				warnings, err := hotNews.ValidateCreate()
				Expect(err).ToNot(HaveOccurred())
				Expect(warnings).To(BeNil())
			})
		})

		Context("when updating HotNews", func() {
			It("should pass the update validation", func() {
				warnings, err := hotNews.ValidateUpdate(nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(warnings).To(BeNil())
			})
		})

		Context("when deleting HotNews", func() {
			It("should pass the deletion validation", func() {
				warnings, err := hotNews.ValidateDelete()
				Expect(err).ToNot(HaveOccurred())
				Expect(warnings).To(BeNil())
			})
		})
	})
})
