package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"radioking-app/internal/api/http/beans"
	"radioking-app/internal/domain/models"
	"radioking-app/internal/domain/services"
	"strconv"

	domainErrors "radioking-app/internal/domain/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

const IdParameter = "id"
const mappingError = "Response mapping error"

var validate = validator.New()

type PlaylistHandler struct {
	service            services.IPlaylistService
	applicationService services.IPlaylistApplicationService
}

func NewPlaylistHandler(service services.IPlaylistService, applicationService services.IPlaylistApplicationService) *PlaylistHandler {
	return &PlaylistHandler{
		service:            service,
		applicationService: applicationService,
	}
}

func (handler *PlaylistHandler) Routes(router *chi.Mux) chi.Router {
	router.Post("/playlists", handler.CreatePlaylist)
	router.Get("/playlists", handler.ListPlaylists)
	router.Get("/playlists/{id}", handler.GetPlaylist)
	router.Post("/playlists/{id}/play", handler.PlayPlaylist)
	return router
}

func (handler *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req beans.PlaylistCreateRequest

	playlist, done := mapRequest(w, r, req, handler)
	if done {
		return
	}

	if err := handler.service.CreatePlaylist(&playlist); err != nil {
		handler.handleBusinessError(w, err)
		return
	}

	var resp beans.PlaylistResponseApiBean
	if err := copier.Copy(&resp, &playlist); err != nil {

		handler.handleError(w, mappingError, http.StatusInternalServerError, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func mapRequest(w http.ResponseWriter, r *http.Request, req beans.PlaylistCreateRequest, handler *PlaylistHandler) (models.Playlist, bool) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.handleError(w, "Invalid JSON payload", http.StatusBadRequest, err)
		return models.Playlist{}, true
	}

	if err := validate.Struct(req); err != nil {
		handler.handleError(w, "Validation failed", http.StatusBadRequest, err)
		return models.Playlist{}, true
	}

	var playlist models.Playlist
	if err := copier.Copy(&playlist, &req); err != nil {
		handler.handleError(w, "Internal mapping error", http.StatusInternalServerError, err)
		return models.Playlist{}, true
	}
	return playlist, false
}

func (handler *PlaylistHandler) ListPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := handler.service.ListPlaylists()
	if err != nil {
		handler.handleBusinessError(w, err)
		return
	}

	render.JSON(w, r, playlists)
}

func (handler *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.extractPlaylistID(w, r)
	if !ok {
		return
	}

	playlist, err := handler.service.GetPlaylist(id)
	if err != nil {
		handler.handleBusinessError(w, err)
		return
	}

	var resp beans.PlaylistResponseApiBean
	if err := copier.Copy(&resp, playlist); err != nil {
		handler.handleError(w, mappingError, http.StatusInternalServerError, err)
		return
	}

	render.JSON(w, r, resp)
}

func (handler *PlaylistHandler) PlayPlaylist(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.extractPlaylistID(w, r)
	if !ok {
		return
	}

	result, err := handler.applicationService.PlayPlaylist(id)
	if err != nil {
		handler.handleBusinessError(w, err)
		return
	}

	var resp beans.PlaylistPlayResponse
	if err := copier.Copy(&resp, result); err != nil {
		handler.handleError(w, mappingError, http.StatusInternalServerError, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// extractPlaylistID helper pour extraire et valider l'ID depuis l'URL
func (handler *PlaylistHandler) extractPlaylistID(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := chi.URLParam(r, IdParameter)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handler.handleError(w, "Invalid playlist ID format", http.StatusBadRequest, err)
		return 0, false
	}
	return id, true
}

func (handler *PlaylistHandler) writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(beans.ErrorResponse{Error: message})
}

func (handler *PlaylistHandler) handleError(w http.ResponseWriter, message string, statusCode int, err error) {
	log.Printf("Handler error: %s - %v", message, err)
	handler.writeJSONError(w, message, statusCode)
}

func (handler *PlaylistHandler) handleBusinessError(w http.ResponseWriter, err error) {
	var businessErr *domainErrors.BusinessError
	if errors.As(err, &businessErr) {
		log.Printf("Business error: %v", err)
		switch {
		case businessErr.IsValidation():
			handler.writeJSONError(w, businessErr.Error(), http.StatusBadRequest)
		case businessErr.IsNotFound():
			handler.writeJSONError(w, businessErr.Error(), http.StatusNotFound)
		default:
			handler.writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	handler.writeJSONError(w, "Internal server error", http.StatusInternalServerError)
}
