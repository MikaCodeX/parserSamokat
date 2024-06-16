package zip

import (
	"archive/zip"
	"os"
)

func CreateZipFile(pluginFile string) error {
	// Create a new zip file
	f, _ := os.Create(pluginFile)

	defer f.Close()

	// Create a new zip writer
	zw := zip.NewWriter(f)

	// Add manifest.json to the zip file
	var fw, err = zw.Create("manifest.json")
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(manifestJson))
	if err != nil {
		return err
	}

	// Add background.js to the zip file
	fw, err = zw.Create("background.js")
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(backgroundJs))
	if err != nil {
		return err
	}

	// Close the zip writer
	err = zw.Close()
	if err != nil {
		return err
	}

	return nil
}
