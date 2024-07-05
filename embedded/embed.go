package embedded

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"resin/pkg/logging"
)

//go:embed login/WebView2Loader.dll login/WebViewLogin-v0.0.5.exe
var LoginFiles embed.FS

func ExtractEmbeddedFiles() {
	read, err := LoginFiles.ReadDir("login")
	if err != nil {
		logging.Fail("Failed to read dir \".\" :(")
		return
	}
	err = os.MkdirAll("login", 0750)
	for i, e := range read {
		path := fmt.Sprintf("login/%s", e.Name())
		winPath := filepath.Join(".", path)

		if _, err := os.Stat(winPath); err == nil {
			continue // File already exists
		}

		file, err := LoginFiles.ReadFile(path)
		if err != nil {
			logging.Fail("failed to read file %d", i)
			continue
		}

		newFile, err := os.Create(winPath)
		defer newFile.Close()
		if err != nil {
			logging.Fail("failed to create file %d", i)
			continue
		}

		n, err := newFile.Write(file)
		if err != nil {
			logging.Fail("failed to write file %d", i)
			continue
		}

		logging.Info("%s: wrote %d bytes", newFile.Name(), i, n)
	}
}
