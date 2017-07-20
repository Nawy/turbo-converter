package model

// ImageJSON response for upload request in convertation process
type ImageJSON struct {
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
}

// StatusResponseJSON status response, contains time and available space on the disk
type StatusResponseJSON struct {
	Time  string `json:"time"`
	Space string `json:"space"`
}

// ErrorResponseJSON error response contains only error description in JSON format
type ErrorResponseJSON struct {
	Description string `json:"description"`
}

// SimpleResponseJSON is simple one line response
type SimpleResponseJSON struct {
	Response string `json:"response"`
}
