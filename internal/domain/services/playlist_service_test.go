package services

import (
	"errors"
	"testing"

	"radioking-app/internal/domain/constants"
	domainErrors "radioking-app/internal/domain/errors"
	"radioking-app/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockPlaylistRepository est un mock du repository
type MockPlaylistRepository struct {
	mock.Mock
}

func (m *MockPlaylistRepository) Create(playlist *models.Playlist) error {
	args := m.Called(playlist)
	return args.Error(0)
}

func (m *MockPlaylistRepository) GetAll() ([]*models.Playlist, error) {
	args := m.Called()
	return args.Get(0).([]*models.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) GetByID(id int) (*models.Playlist, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Playlist), args.Error(1)
}

func TestPlaylistService_CreatePlaylist_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	playlist := &models.Playlist{
		Name: "Test Playlist",
		Tracks: []models.Track{
			{Title: "Song 1", Artist: "Artist 1"},
		},
	}

	mockRepo.On("Create", playlist).Return(nil)

	// Act
	err := service.CreatePlaylist(playlist)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_CreatePlaylist_EmptyName(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	playlist := &models.Playlist{
		Name:   "",
		Tracks: []models.Track{},
	}

	// Act
	err := service.CreatePlaylist(playlist)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domainErrors.ErrEmptyPlaylistName, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestPlaylistService_CreatePlaylist_NameTooLong(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	longName := make([]rune, constants.MaxPlaylistNameLength+1)
	for i := range longName {
		longName[i] = 'a'
	}

	playlist := &models.Playlist{
		Name:   string(longName),
		Tracks: []models.Track{},
	}

	// Act
	err := service.CreatePlaylist(playlist)

	// Assert
	assert.Error(t, err)
	var validationErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &validationErr))
	assert.True(t, validationErr.IsValidation())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestPlaylistService_CreatePlaylist_TooManyTracks(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	tracks := make([]models.Track, constants.MaxTracksPerPlaylist+1)
	for i := range tracks {
		tracks[i] = models.Track{Title: "Song", Artist: "Artist"}
	}

	playlist := &models.Playlist{
		Name:   "Test Playlist",
		Tracks: tracks,
	}

	// Act
	err := service.CreatePlaylist(playlist)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domainErrors.ErrTooManyTracks, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestPlaylistService_CreatePlaylist_InvalidTrack(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	playlist := &models.Playlist{
		Name: "Test Playlist",
		Tracks: []models.Track{
			{Title: "", Artist: "Artist 1"}, // Empty title
		},
	}

	// Act
	err := service.CreatePlaylist(playlist)

	// Assert
	assert.Error(t, err)
	var validationErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &validationErr))
	assert.True(t, validationErr.IsValidation())
	assert.Contains(t, err.Error(), "track 1 invalid")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestPlaylistService_CreatePlaylist_RepoError(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	playlist := &models.Playlist{
		Name:   "Test Playlist",
		Tracks: []models.Track{},
	}

	expectedErr := errors.New("database error")
	mockRepo.On("Create", playlist).Return(expectedErr)

	// Act
	err := service.CreatePlaylist(playlist)

	// Assert
	assert.Error(t, err)
	var internalErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &internalErr))
	assert.Equal(t, domainErrors.InternalError, internalErr.Type)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_GetPlaylist_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	expectedPlaylist := &models.Playlist{
		Name: "Test Playlist",
		Tracks: []models.Track{
			{Title: "Song 1", Artist: "Artist 1"},
		},
	}

	mockRepo.On("GetByID", 1).Return(expectedPlaylist, nil)

	// Act
	result, err := service.GetPlaylist(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPlaylist, result)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_GetPlaylist_InvalidID(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	// Act
	result, err := service.GetPlaylist(0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrInvalidPlaylistID, err)
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestPlaylistService_GetPlaylist_NegativeID(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	// Act
	result, err := service.GetPlaylist(-1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrInvalidPlaylistID, err)
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestPlaylistService_GetPlaylist_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	mockRepo.On("GetByID", 999).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := service.GetPlaylist(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrPlaylistNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_GetPlaylist_RepoError(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	expectedErr := errors.New("database connection failed")
	mockRepo.On("GetByID", 1).Return(nil, expectedErr)

	// Act
	result, err := service.GetPlaylist(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var internalErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &internalErr))
	assert.Equal(t, domainErrors.InternalError, internalErr.Type)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_ListPlaylists_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	expectedPlaylists := []*models.Playlist{
		{
			Name: "Playlist 1",
			Tracks: []models.Track{
				{Title: "Song 1", Artist: "Artist 1"},
			},
		},
		{
			Name: "Playlist 2",
			Tracks: []models.Track{
				{Title: "Song 2", Artist: "Artist 2"},
				{Title: "Song 3", Artist: "Artist 3"},
			},
		},
	}

	mockRepo.On("GetAll").Return(expectedPlaylists, nil)

	// Act
	result, err := service.ListPlaylists()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPlaylists, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_ListPlaylists_EmptyResult(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	expectedPlaylists := []*models.Playlist{}
	mockRepo.On("GetAll").Return(expectedPlaylists, nil)

	// Act
	result, err := service.ListPlaylists()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPlaylists, result)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestPlaylistService_ListPlaylists_RepoError(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlaylistRepository)
	service := &PlaylistService{Repo: mockRepo}

	expectedErr := errors.New("database connection failed")
	mockRepo.On("GetAll").Return([]*models.Playlist(nil), expectedErr)

	// Act
	result, err := service.ListPlaylists()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var internalErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &internalErr))
	assert.Equal(t, domainErrors.InternalError, internalErr.Type)
	assert.Contains(t, err.Error(), "failed to list playlists")
	mockRepo.AssertExpectations(t)
}

// Tests for validation methods (private methods tested through public methods)

func TestPlaylistService_ValidateTrack_EmptyTitle(t *testing.T) {
	// Arrange
	service := &PlaylistService{}
	track := &models.Track{Title: "", Artist: "Artist"}

	// Act
	err := service.validateTrack(track)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domainErrors.ErrEmptyTrackTitle, err)
}

func TestPlaylistService_ValidateTrack_EmptyArtist(t *testing.T) {
	// Arrange
	service := &PlaylistService{}
	track := &models.Track{Title: "Song", Artist: ""}

	// Act
	err := service.validateTrack(track)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domainErrors.ErrEmptyTrackArtist, err)
}

func TestPlaylistService_ValidateTrack_TitleTooLong(t *testing.T) {
	// Arrange
	service := &PlaylistService{}

	longTitle := make([]rune, constants.MaxTrackNameLength+1)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	track := &models.Track{Title: string(longTitle), Artist: "Artist"}

	// Act
	err := service.validateTrack(track)

	// Assert
	assert.Error(t, err)
	var validationErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &validationErr))
	assert.True(t, validationErr.IsValidation())
	assert.Contains(t, err.Error(), "track title too long")
}

func TestPlaylistService_ValidateTrack_ArtistTooLong(t *testing.T) {
	// Arrange
	service := &PlaylistService{}

	longArtist := make([]rune, constants.MaxArtistNameLength+1)
	for i := range longArtist {
		longArtist[i] = 'a'
	}

	track := &models.Track{Title: "Song", Artist: string(longArtist)}

	// Act
	err := service.validateTrack(track)

	// Assert
	assert.Error(t, err)
	var validationErr *domainErrors.BusinessError
	assert.True(t, errors.As(err, &validationErr))
	assert.True(t, validationErr.IsValidation())
	assert.Contains(t, err.Error(), "artist name too long")
}

func TestPlaylistService_ValidateTrack_Success(t *testing.T) {
	// Arrange
	service := &PlaylistService{}
	track := &models.Track{Title: "Valid Song", Artist: "Valid Artist"}

	// Act
	err := service.validateTrack(track)

	// Assert
	assert.NoError(t, err)
}

func TestPlaylistService_ValidateTrack_WhitespaceOnly(t *testing.T) {
	// Arrange
	service := &PlaylistService{}
	track := &models.Track{Title: "   ", Artist: "   "}

	// Act
	err := service.validateTrack(track)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domainErrors.ErrEmptyTrackTitle, err)
}
