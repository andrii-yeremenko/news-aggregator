package v1

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HotNews Validation", func() {
	BeforeEach(func() {
		k8sClient = fake.NewClientBuilder().WithObjects(&v1.ConfigMap{
			Data: map[string]string{
				"group1": "feed1",
				"group2": "feed2",
			},
		}).Build()
	})

	Describe("ValidateHotNews", func() {
		Context("when validating HotNews", func() {

			It("should validate a valid HotNews", func() {
				hotNews := HotNews{
					Spec: HotNewsSpec{
						Keywords:  []string{"news", "update"},
						DateStart: &metav1.Time{Time: time.Now()},
						DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
					},
				}
				err := hotNews.validateHotNews()
				Expect(err).To(HaveLen(0))
			})

			It("should not validate with empty keywords", func() {
				hotNews := HotNews{
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

			It("should not validate when dateEnd is before dateStart", func() {
				hotNews := HotNews{
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

			It("should not validate when dateStart is specified without dateEnd", func() {
				hotNews := HotNews{
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

	Describe("ValidateFeedGroups", func() {
		Context("when validating feed groups", func() {

			It("should validate valid feed groups", func() {
				hotNews := HotNews{
					Spec: HotNewsSpec{
						FeedGroups: []string{"group1", "group2"},
					},
				}
				err := hotNews.validateFeedGroups()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not validate with non-existent feed groups", func() {
				hotNews := HotNews{
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

	Describe("Validation Operations", func() {
		var hotNews HotNews

		BeforeEach(func() {
			hotNews = HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news"},
					DateStart: &metav1.Time{Time: time.Now()},
					DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
				},
			}
		})

		It("should validate creation of HotNews", func() {
			warnings, err := hotNews.ValidateCreate()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should validate update of HotNews", func() {
			warnings, err := hotNews.ValidateUpdate(nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should validate deletion of HotNews", func() {
			warnings, err := hotNews.ValidateDelete()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})
	})
})
