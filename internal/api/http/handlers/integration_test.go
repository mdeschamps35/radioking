package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"radioking-app/internal/api/http/beans"
	"radioking-app/internal/config"
	"radioking-app/internal/domain/services"
	"radioking-app/internal/infrastructure/db"
	"radioking-app/internal/infrastructure/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Test constants
const (
	PlaylistsEndpoint = "/playlists"
	ContentTypeJSON   = "application/json"

	// Test data
	TestPlaylistName1 = "Playlist 1"
	TestPlaylistName2 = "Playlist 2"
	TestPlaylistName  = "Test Playlist"
	TestSong1Title    = "Song 1"
	TestSong2Title    = "Song 2"
	TestSong3Title    = "Song 3"
	TestArtist1Name   = "Artist 1"
	TestArtist2Name   = "Artist 2"
	TestArtist3Name   = "Artist 3"

	// Error messages
	ValidationErrorMsg = "Validation"
	InvalidIDErrorMsg  = "Invalid"
	NotFoundErrorMsg   = "not found"

	// Test IDs
	NonExistentID = 999
	InvalidIDStr  = "invalid"
)

type IntegrationTestSuite struct {
	suite.Suite
	router *chi.Mux
	db     *gorm.DB
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Set environment variable to disable auth for tests
	_ = os.Setenv("RADIOKING_AUTH_ENABLED", "false")

	// Load config
	cfg, err := config.Load()
	suite.Require().NoError(err)
	suite.Require().False(cfg.Auth.Enabled, "Auth should be disabled for integration tests")

	// Setup test database
	testDB, err := db.InitDb()
	suite.Require().NoError(err)
	suite.db = testDB

	// Setup router without auth middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	// Setup handlers
	repo := repositories.NewPlaylistRepository(testDB)
	service := services.PlaylistService{Repo: repo}
	playService := services.NewPlaylistPlayService(&service, nil) // pas besoin du publisher pour les tests
	appService := services.NewPlaylistApplicationService(&service, playService)
	handler := NewPlaylistHandler(&service, appService)
	handler.Routes(router)

	suite.router = router
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	// Clean up environment variable
	_ = os.Unsetenv("RADIOKING_AUTH_ENABLED")
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.db.Exec("DELETE FROM tracks").Error
	suite.Require().NoError(err)
	err = suite.db.Exec("DELETE FROM playlists").Error
	suite.Require().NoError(err)
}

// Helper Methods

func (suite *IntegrationTestSuite) createTestPlaylist(name string, tracks []beans.TrackCreateRequest) int {
	requestBody := suite.buildPlaylistRequest(name, tracks...)
	rr := suite.makePostRequest(PlaylistsEndpoint, requestBody)
	suite.Require().Equal(http.StatusCreated, rr.Code)

	response := suite.parsePlaylistResponse(rr)
	return int(response.ID)
}

func (suite *IntegrationTestSuite) makePostRequest(endpoint string, payload interface{}) *httptest.ResponseRecorder {
	body, err := json.Marshal(payload)
	suite.Require().NoError(err)

	req := httptest.NewRequest("POST", endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", ContentTypeJSON)
	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)
	return rr
}

func (suite *IntegrationTestSuite) makeGetRequest(endpoint string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", endpoint, nil)
	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)
	return rr
}

func (suite *IntegrationTestSuite) parsePlaylistResponse(rr *httptest.ResponseRecorder) beans.PlaylistResponseApiBean {
	var response beans.PlaylistResponseApiBean
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.Require().NoError(err)
	return response
}

func (suite *IntegrationTestSuite) parsePlaylistsResponse(rr *httptest.ResponseRecorder) []beans.PlaylistResponseApiBean {
	var response []beans.PlaylistResponseApiBean
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.Require().NoError(err)
	return response
}

func (suite *IntegrationTestSuite) parseErrorResponse(rr *httptest.ResponseRecorder) beans.ErrorResponse {
	var response beans.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.Require().NoError(err)
	return response
}

// Test Data Builders

func (suite *IntegrationTestSuite) buildPlaylistRequest(name string, tracks ...beans.TrackCreateRequest) beans.PlaylistCreateRequest {
	return beans.PlaylistCreateRequest{
		Name:   name,
		Tracks: tracks,
	}
}

func (suite *IntegrationTestSuite) buildTrack(title, artist string) beans.TrackCreateRequest {
	return beans.TrackCreateRequest{
		Title:  title,
		Artist: artist,
	}
}

func (suite *IntegrationTestSuite) buildTestPlaylist() beans.PlaylistCreateRequest {
	return suite.buildPlaylistRequest(TestPlaylistName,
		suite.buildTrack(TestSong1Title, TestArtist1Name),
		suite.buildTrack(TestSong2Title, TestArtist2Name),
	)
}

func (suite *IntegrationTestSuite) buildEmptyPlaylist() beans.PlaylistCreateRequest {
	return suite.buildPlaylistRequest("" /* no tracks */)
}

func (suite *IntegrationTestSuite) TestCreatePlaylist_Success() {
	// Arrange
	requestBody := suite.buildTestPlaylist()

	// Act
	rr := suite.makePostRequest(PlaylistsEndpoint, requestBody)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, rr.Code)

	response := suite.parsePlaylistResponse(rr)
	assert.Equal(suite.T(), TestPlaylistName, response.Name)
	assert.Len(suite.T(), response.Tracks, 2)
	assert.Equal(suite.T(), TestSong1Title, response.Tracks[0].Title)
	assert.Equal(suite.T(), TestArtist1Name, response.Tracks[0].Artist)
	assert.NotZero(suite.T(), response.ID)
}

func (suite *IntegrationTestSuite) TestCreatePlaylist_EmptyName() {
	// Arrange
	requestBody := suite.buildEmptyPlaylist()

	// Act
	rr := suite.makePostRequest(PlaylistsEndpoint, requestBody)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)

	errorResponse := suite.parseErrorResponse(rr)
	assert.Contains(suite.T(), errorResponse.Error, ValidationErrorMsg)
}

func (suite *IntegrationTestSuite) TestGetPlaylists_Success() {
	// Arrange - Create test data
	suite.createTestPlaylist(TestPlaylistName1, []beans.TrackCreateRequest{
		suite.buildTrack(TestSong1Title, TestArtist1Name),
	})
	suite.createTestPlaylist(TestPlaylistName2, []beans.TrackCreateRequest{
		suite.buildTrack(TestSong2Title, TestArtist2Name),
		suite.buildTrack(TestSong3Title, TestArtist3Name),
	})

	// Act
	rr := suite.makeGetRequest(PlaylistsEndpoint)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	response := suite.parsePlaylistsResponse(rr)
	assert.Len(suite.T(), response, 2)
	assert.Equal(suite.T(), TestPlaylistName1, response[0].Name)
	assert.Len(suite.T(), response[0].Tracks, 1)
	assert.Equal(suite.T(), TestPlaylistName2, response[1].Name)
	assert.Len(suite.T(), response[1].Tracks, 2)
}

func (suite *IntegrationTestSuite) TestGetPlaylists_Empty() {
	// Arrange - No test data

	// Act
	rr := suite.makeGetRequest(PlaylistsEndpoint)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	response := suite.parsePlaylistsResponse(rr)
	assert.Len(suite.T(), response, 0)
}

func (suite *IntegrationTestSuite) TestGetPlaylistByID_Success() {
	// Arrange - Create test playlist
	testPlaylist := suite.buildTestPlaylist()
	playlistID := suite.createTestPlaylist(testPlaylist.Name, testPlaylist.Tracks)

	// Act
	rr := suite.makeGetRequest(fmt.Sprintf("%s/%d", PlaylistsEndpoint, playlistID))

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	response := suite.parsePlaylistResponse(rr)
	assert.Equal(suite.T(), int64(playlistID), response.ID)
	assert.Equal(suite.T(), TestPlaylistName, response.Name)
	assert.Len(suite.T(), response.Tracks, 2)
}

func (suite *IntegrationTestSuite) TestGetPlaylistByID_NotFound() {
	// Act
	rr := suite.makeGetRequest(fmt.Sprintf("%s/%d", PlaylistsEndpoint, NonExistentID))

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, rr.Code)

	errorResponse := suite.parseErrorResponse(rr)
	assert.Contains(suite.T(), errorResponse.Error, NotFoundErrorMsg)
}

func (suite *IntegrationTestSuite) TestGetPlaylistByID_InvalidID() {
	// Act
	rr := suite.makeGetRequest(fmt.Sprintf("%s/%s", PlaylistsEndpoint, InvalidIDStr))

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)

	errorResponse := suite.parseErrorResponse(rr)
	assert.Contains(suite.T(), errorResponse.Error, InvalidIDErrorMsg)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
