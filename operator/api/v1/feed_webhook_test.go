package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func init() {
	feedlog = zap.New(zap.UseDevMode(true))
}

func TestValidateFeedName(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()

	tests := []struct {
		name      string
		feed      *Feed
		expectErr bool
	}{
		{
			name: "valid name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "valid-name",
					Link: "https://example.com",
				},
			},
			expectErr: false,
		},
		{
			name: "empty name",
			feed: &Feed{
				Spec: FeedSpec{Name: ""},
			},
			expectErr: true,
		},
		{
			name: "name too long",
			feed: &Feed{
				Spec: FeedSpec{Name: "this-name-is-way-too-long-to-be-valid"},
			},
			expectErr: true,
		},
		{
			name: "name with invalid characters",
			feed: &Feed{
				Spec: FeedSpec{Name: "invalid!name"},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFeed(tt.feed)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateFeedLink тестирует валидацию поля Link в Feed.
func TestValidateFeedLink(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()

	tests := []struct {
		name      string
		feed      *Feed
		expectErr bool
	}{
		{
			name: "valid link",
			feed: &Feed{
				Spec: FeedSpec{
					Link: "http://example.com",
					Name: "valid-link",
				},
			},
			expectErr: false,
		},
		{
			name: "invalid link",
			feed: &Feed{
				Spec: FeedSpec{
					Link: "invalid-link",
					Name: "invalid-link",
				},
			},
			expectErr: true,
		},
		{
			name: "empty link",
			feed: &Feed{
				Spec: FeedSpec{
					Link: "",
					Name: "empty-link",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFeed(tt.feed)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNameUniqueness(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	existingFeed := &Feed{
		Spec: FeedSpec{
			Name: "existing-feed",
			Link: "https://example.com",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			UID:       "existing-uid",
		},
	}
	existingFeedList := &FeedList{
		Items: []Feed{*existingFeed},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(existingFeedList).Build()

	tests := []struct {
		name      string
		feed      *Feed
		expectErr bool
	}{
		{
			name: "unique name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "new-feed",
					Link: "https://example.com",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: existingFeed.Namespace,
					UID:       "new-uid",
				},
			},
			expectErr: false,
		},
		{
			name: "duplicate name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "existing-feed",
					Link: "https://example.com",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: existingFeed.Namespace,
					UID:       "new-uid",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkNameUniqueness(tt.feed)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
