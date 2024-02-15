package remix

// Clip describes the JSON format used to represent a single video clip that can be
// queued for playback via the remix service
type Clip struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	TapeId   int    `json:"tapeId"`
}

// ClipListing is the response payload for GET /clips
type ClipListing struct {
	Clips []Clip `json:"clips"`
}
