package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMain(m *testing.M) {
	k8sClient = fake.NewClientBuilder().WithObjects(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hotnews-feeds-group",
			Namespace: "news-aggregator-namespace",
		},
		Data: map[string]string{
			"group1": "feed1",
			"group2": "feed2",
		},
	}).Build()

	m.Run()
}

func TestValidateHotNews(t *testing.T) {
	tests := []struct {
		name      string
		hotNews   HotNews
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid HotNews",
			hotNews: HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news", "update"},
					DateStart: &metav1.Time{Time: time.Now()},
					DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
				},
			},
			expectErr: false,
		},
		{
			name: "empty keywords",
			hotNews: HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{},
					DateStart: &metav1.Time{Time: time.Now()},
					DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
				},
			},
			expectErr: true,
			errMsg:    "keywords must not be empty",
		},
		{
			name: "dateEnd before dateStart",
			hotNews: HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news"},
					DateStart: &metav1.Time{Time: time.Now()},
					DateEnd:   &metav1.Time{Time: time.Now().Add(-24 * time.Hour)},
				},
			},
			expectErr: true,
			errMsg:    "dateEnd must be after dateStart",
		},
		{
			name: "dateStart without dateEnd",
			hotNews: HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news"},
					DateStart: &metav1.Time{Time: time.Now()},
				},
			},
			expectErr: true,
			errMsg:    "dateEnd must be provided if dateStart is specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hotNews.validateHotNews()
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFeedGroups(t *testing.T) {
	tests := []struct {
		name      string
		hotNews   HotNews
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid feed groups",
			hotNews: HotNews{
				Spec: HotNewsSpec{
					FeedGroups: []string{"group1", "group2"},
				},
			},
			expectErr: false,
		},
		{
			name: "non-existent feed group",
			hotNews: HotNews{
				Spec: HotNewsSpec{
					FeedGroups: []string{"group1", "group3"},
				},
			},
			expectErr: true,
			errMsg:    "feedGroup group3 does not exist in ConfigMap news-aggregator-namespace/hotnews-feeds-group",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hotNews.validateFeedGroups()
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCreate(t *testing.T) {
	hotNews := HotNews{
		Spec: HotNewsSpec{
			Keywords:  []string{"news"},
			DateStart: &metav1.Time{Time: time.Now()},
			DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
		},
	}

	warnings, err := hotNews.ValidateCreate()
	assert.NoError(t, err)
	assert.Nil(t, warnings)
}

func TestValidateUpdate(t *testing.T) {
	hotNews := HotNews{
		Spec: HotNewsSpec{
			Keywords:  []string{"news"},
			DateStart: &metav1.Time{Time: time.Now()},
			DateEnd:   &metav1.Time{Time: time.Now().Add(24 * time.Hour)},
		},
	}

	warnings, err := hotNews.ValidateUpdate(nil)
	assert.NoError(t, err)
	assert.Nil(t, warnings)
}

func TestValidateDelete(t *testing.T) {
	hotNews := HotNews{}

	warnings, err := hotNews.ValidateDelete()
	assert.NoError(t, err)
	assert.Nil(t, warnings)
}
