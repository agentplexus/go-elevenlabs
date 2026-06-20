package elevenlabs

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// AnalyticsService handles conversation analytics operations.
type AnalyticsService struct {
	client *Client
}

// GetLiveCount returns the number of active ongoing conversations.
func (s *AnalyticsService) GetLiveCount(ctx context.Context) (int, error) {
	resp, err := s.client.apiClient.GetLiveCount(ctx, api.GetLiveCountParams{})
	if err != nil {
		return 0, err
	}

	result, ok := resp.(*api.GetLiveCountResponse)
	if !ok {
		return 0, errors.New("unexpected response type")
	}

	return result.Count, nil
}
