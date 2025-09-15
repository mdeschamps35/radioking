package beans

type PlaylistPlayResponse struct {
	Message     string `json:"message"`
	PlaylistID  int    `json:"playlist_id"`
	TracksCount int    `json:"tracks_count"`
}
