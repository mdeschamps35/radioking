package beans

type PlaylistResponseApiBean struct {
	ID     int64                  `json:"id"`
	Name   string                 `json:"name"`
	Tracks []TrackResponseApiBean `json:"tracks"`
}

type PlaylistCreateRequest struct {
	Name   string               `json:"name" validate:"required,min=1,max=255"`
	Tracks []TrackCreateRequest `json:"tracks" validate:"max=100,dive"`
}
