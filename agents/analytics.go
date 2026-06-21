package agents

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// GetLiveCount returns the number of active ongoing conversations.
func (s *Service) GetLiveCount(ctx context.Context) (int, error) {
	resp, err := s.client.API().GetLiveCount(ctx, api.GetLiveCountParams{})
	if err != nil {
		return 0, err
	}

	result, ok := resp.(*api.GetLiveCountResponse)
	if !ok {
		return 0, errors.New("unexpected response type")
	}

	return result.Count, nil
}
