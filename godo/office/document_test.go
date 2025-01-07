package office

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestGetDocument(t *testing.T) {
	// Get the absolute path to the testdata directory
	testdataDir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatalf("Failed to get absolute path to testdata directory: %v", err)
	}

	// Read all files in the testdata directory
	files, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	// Iterate over each file and test GetDocument
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(testdataDir, file.Name())
			t.Run(file.Name(), func(t *testing.T) {
				doc, err := GetDocument(filePath)
				if err != nil {
					t.Errorf("Failed to get document for %s: %v", file.Name(), err)
				} else {
					log.Printf("Document file.Name: %s\ncontent: %d\n", file.Name(), len(doc.Content))
					//t.Logf("Document file.Name: %s\n", file.Name())
				}
			})
		}
	}
}
