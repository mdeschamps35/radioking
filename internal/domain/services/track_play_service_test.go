package services

import (
	"errors"
	"testing"
	"time"

	"radioking-app/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTrackPlayRepository is a mock implementation of TrackPlayRepository
type MockTrackPlayRepository struct {
	mock.Mock
}

func (m *MockTrackPlayRepository) Create(trackPlay *models.TrackPlay) error {
	args := m.Called(trackPlay)
	return args.Error(0)
}

func (m *MockTrackPlayRepository) GetByPlaylistID(playlistID int) ([]*models.TrackPlay, error) {
	args := m.Called(playlistID)
	return args.Get(0).([]*models.TrackPlay), args.Error(1)
}

func (m *MockTrackPlayRepository) GetByTrackID(trackID int) ([]*models.TrackPlay, error) {
	args := m.Called(trackID)
	return args.Get(0).([]*models.TrackPlay), args.Error(1)
}

func TestTrackPlayService_RecordTrackPlay(t *testing.T) {
	tests := []struct {
		name    string
		event   models.TrackPlayedEvent
		mockFn  func(*MockTrackPlayRepository)
		wantErr bool
	}{
		{
			name: "successful record",
			event: models.TrackPlayedEvent{
				PlaylistID: 1,
				TrackID:    1,
				Position:   0,
				PlayedAt:   time.Now(),
				EventID:    "test-event-id",
			},
			mockFn: func(repo *MockTrackPlayRepository) {
				repo.On("Create", mock.AnythingOfType("*models.TrackPlay")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			event: models.TrackPlayedEvent{
				PlaylistID: 1,
				TrackID:    1,
				Position:   0,
				PlayedAt:   time.Now(),
				EventID:    "test-event-id",
			},
			mockFn: func(repo *MockTrackPlayRepository) {
				repo.On("Create", mock.AnythingOfType("*models.TrackPlay")).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockTrackPlayRepository)
			tt.mockFn(repo)

			service := NewTrackPlayService(repo)
			err := service.RecordTrackPlay(tt.event)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to record track play")
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestTrackPlayService_GetPlaylistPlays(t *testing.T) {
	tests := []struct {
		name       string
		playlistID int
		mockFn     func(*MockTrackPlayRepository)
		want       []*models.TrackPlay
		wantErr    bool
	}{
		{
			name:       "successful retrieval",
			playlistID: 1,
			mockFn: func(repo *MockTrackPlayRepository) {
				plays := []*models.TrackPlay{
					{ID: 1, PlaylistID: 1, TrackID: 1, Position: 0},
					{ID: 2, PlaylistID: 1, TrackID: 2, Position: 1},
				}
				repo.On("GetByPlaylistID", 1).Return(plays, nil)
			},
			want: []*models.TrackPlay{
				{ID: 1, PlaylistID: 1, TrackID: 1, Position: 0},
				{ID: 2, PlaylistID: 1, TrackID: 2, Position: 1},
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			playlistID: 1,
			mockFn: func(repo *MockTrackPlayRepository) {
				repo.On("GetByPlaylistID", 1).Return([]*models.TrackPlay(nil), errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockTrackPlayRepository)
			tt.mockFn(repo)

			service := NewTrackPlayService(repo)
			result, err := service.GetPlaylistPlays(tt.playlistID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "failed to get playlist plays")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestTrackPlayService_GetTrackPlays(t *testing.T) {
	tests := []struct {
		name    string
		trackID int
		mockFn  func(*MockTrackPlayRepository)
		want    []*models.TrackPlay
		wantErr bool
	}{
		{
			name:    "successful retrieval",
			trackID: 1,
			mockFn: func(repo *MockTrackPlayRepository) {
				plays := []*models.TrackPlay{
					{ID: 1, PlaylistID: 1, TrackID: 1, Position: 0},
					{ID: 3, PlaylistID: 2, TrackID: 1, Position: 2},
				}
				repo.On("GetByTrackID", 1).Return(plays, nil)
			},
			want: []*models.TrackPlay{
				{ID: 1, PlaylistID: 1, TrackID: 1, Position: 0},
				{ID: 3, PlaylistID: 2, TrackID: 1, Position: 2},
			},
			wantErr: false,
		},
		{
			name:    "repository error",
			trackID: 1,
			mockFn: func(repo *MockTrackPlayRepository) {
				repo.On("GetByTrackID", 1).Return([]*models.TrackPlay(nil), errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockTrackPlayRepository)
			tt.mockFn(repo)

			service := NewTrackPlayService(repo)
			result, err := service.GetTrackPlays(tt.trackID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "failed to get track plays")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			repo.AssertExpectations(t)
		})
	}
}
