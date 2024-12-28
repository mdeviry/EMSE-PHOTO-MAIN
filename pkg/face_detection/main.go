package main

import (
    "log"
    "face_detection"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
	"face_detection/models"
)

func connectMariaDB() (*gorm.DB, error) {
    dsn := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
    db, err := connectMariaDB()
    if err != nil {
        log.Fatalf("Échec de la connexion à MariaDB : %v", err)
    }

    err = face_detection.InitializeFaceDetector(db)
    if err != nil {
        log.Fatalf("Erreur d'initialisation : %v", err)
    }

    media := &models.Media{
        ID: 1,
        MediaURL: []models.MediaURL{
            {Purpose: "thumbnail", URL: "example_thumbnail.jpg"},
        },
    }

    if face_detection.GlobalFaceDetector != nil {
        err = face_detection.GlobalFaceDetector.DetectFaces(db, media)
        if err != nil {
            log.Fatalf("Erreur de détection : %v", err)
        } else {
            log.Println("Détection réussie.")
        }
    } else {
        log.Println("Détecteur non initialisé.")
    }
}
