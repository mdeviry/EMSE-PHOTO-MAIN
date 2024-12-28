package models

type Media struct {
    ID       int
    MediaURL []MediaURL
}

type MediaURL struct {
    Purpose string
    URL     string
}

func (url *MediaURL) CachedPath() (string, error) {
    // Simulez le chemin du cache
    return "/cache/" + url.URL, nil
}

type ImageFace struct {
    ID         int
    MediaID    int
    Descriptor []float32 // Aligné avec face.Descriptor
    Rectangle  string    // Coordonnées de détection du visage
}

type FaceGroup struct {
    ID         int
    Label      *string
    ImageFaces []ImageFace
}

type Photo struct {
	PhotoID     uint   `gorm:"primaryKey"`
	PathToPhoto string `gorm:"not null"`
	EventID     uint
}

type RecognizedUser struct {
	RecognizedUserID uint `gorm:"primaryKey"`
	UserID           uint
	PhotoID          uint
}
