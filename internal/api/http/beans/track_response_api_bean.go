package beans

type TrackResponseApiBean struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type TrackCreateRequest struct {
	Title  string `json:"title" validate:"required,min=1,max=255"`
	Artist string `json:"artist" validate:"required,min=1,max=255"`
}
