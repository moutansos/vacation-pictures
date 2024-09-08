package data

type VacationsRoot struct {
	Vacations []Vacation `json:"vacations"`
}

type Vacation struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	MainThumbnail string        `json:"mainThumbnail"`
	Pictures      []VacaPicture `json:"pictures"`
}

type VacaPicture struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	ThumbnailPath string   `json:"thumbnailPath"`
	ImagePath     string   `json:"imagePath"`
	Tags          []string `json:"tags"`
	Rotate        *string  `json:"rotate"`
}
