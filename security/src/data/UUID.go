package data

import (
	"log"
	"strings"

	"github.com/gofrs/uuid"
)

// GenerateUUID génère un UUID Version 4 (aléatoire)
func GenerateUUID() []byte {
	// Générer un nouvel UUID (Version 4)
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Erreur lors de la génération de l'UUID : %v", err)
	}
	for strings.Contains(string(id.Bytes()),",") {
		id, _ = uuid.NewV4()
	}
	return id.Bytes() // Retourne l'UUID en tant que chaîne
}
