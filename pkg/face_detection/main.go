package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"face_detection"
	"face_detection/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connecte à MariaDB
func connectMariaDB() (*gorm.DB, error) {
	dsn := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// Analyse les photos dans un répertoire
func analyzePhotos(db *gorm.DB, folderPath string) ([]map[string]interface{}, error) {
	// Liste pour stocker les résultats
	results := []map[string]interface{}{}

	// Initialisez le détecteur de visages
	err := face_detection.InitializeFaceDetector(db)
	if err != nil {
		return nil, fmt.Errorf("erreur d'initialisation du détecteur : %w", err)
	}

	// Parcourez toutes les photos dans le dossier
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du dossier : %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Analysez chaque photo
		filePath := filepath.Join(folderPath, file.Name())
		media := &models.Media{
			ID:       0, // Optionnel : Remplissez l'ID si nécessaire
			MediaURL: []models.MediaURL{{Purpose: "photo", URL: filePath}},
		}

		err := face_detection.GlobalFaceDetector.DetectFaces(db, media)
		if err != nil {
			log.Printf("Erreur lors de la détection sur %s : %v", filePath, err)
			continue
		}

		// Enregistrez les résultats
		results = append(results, map[string]interface{}{
			"photo_path": filePath,
			"faces":      media.MediaURL,
		})
	}

	return results, nil
}

// Sauvegarde les résultats en JSON
func saveResultsToFile(results []map[string]interface{}, outputPath string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de la conversion en JSON : %w", err)
	}

	err = ioutil.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier : %w", err)
	}

	return nil
}

func main() {
	// Connectez-vous à la base de données
	db, err := connectMariaDB()
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données : %v", err)
	}

	// Assurez-vous que le dossier /photos/new existe
	folderPath := "./photos/new"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Fatalf("Le dossier %s n'existe pas", folderPath)
	}

	// Analysez les photos dans le dossier
	results, err := analyzePhotos(db, folderPath)
	if err != nil {
		log.Fatalf("Erreur d'analyse des photos : %v", err)
	}

	// Sauvegardez les résultats dans results.json
	outputPath := "results.json"
	err = saveResultsToFile(results, outputPath)
	if err != nil {
		log.Fatalf("Erreur lors de l'écriture des résultats : %v", err)
	}

	log.Println("Analyse terminée. Résultats sauvegardés dans", outputPath)
}
