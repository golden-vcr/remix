package remix

// ClipSync is the payload for an operation to register a new clip via POST /admin/clip
type ClipSync struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	TapeId   int    `json:"tapeId"`
}
