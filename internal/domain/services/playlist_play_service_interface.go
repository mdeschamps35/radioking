package services

// IPlaylistPlayService interface pour jouer des playlists
type IPlaylistPlayService interface {
	PlayPlaylist(playlistID int) error
}
