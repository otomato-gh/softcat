package models

// componet represents data about a software component.
type Component struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Team     int    `json:"team"`
	Language string `json:"language"`
}

// team represents data about a team.
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// image represents data about an image.
type Image struct {
	ID    int    `json:"id"`
	Image []byte `json:"image"`
}
