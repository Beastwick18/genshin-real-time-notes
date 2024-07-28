package embedded

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"resin/pkg/logging"
)

//go:embed assets/*
var AssetFiles embed.FS

//go:embed login/WebView2Loader.dll login/WebViewLogin-v0.0.5.exe
var LoginFiles embed.FS

func ReadAssets[T any](a *T) {
	val := reflect.ValueOf(a)
	elem := val.Elem()
	for i := 0; i < elem.NumField(); i++ {
		file, ok := elem.Type().Field(i).Tag.Lookup("asset")
		if !ok {
			continue // no tag
		}
		bytes, err := AssetFiles.ReadFile(fmt.Sprintf("assets/%s", file))
		// Panic on failure to read any asset
		if err != nil {
			logging.Panic("Failed to read assets:\n%v", err)
			os.Exit(1)
			return
		}

		elem.Field(i).SetBytes(bytes)
	}
	return
}

func ExtractEmbeddedFiles() {
	read, err := LoginFiles.ReadDir("login")
	if err != nil {
		logging.Fail("Failed to read asset dir \"login\":\n%v", err)
		return
	}
	err = os.MkdirAll("login", 0755)
	for i, e := range read {
		path := fmt.Sprintf("login/%s", e.Name())
		winPath := filepath.Join(".", path)

		if _, err := os.Stat(winPath); err == nil {
			continue // File already exists
		}

		file, err := LoginFiles.ReadFile(path)
		if err != nil {
			logging.Fail("failed to read file %d:\n%v", i, err)
			continue
		}

		newFile, err := os.Create(winPath)
		defer newFile.Close()
		if err != nil {
			logging.Fail("failed to create file %d:\n%v", i, err)
			continue
		}

		n, err := newFile.Write(file)
		if err != nil {
			logging.Fail("failed to write file %d:\n%v", i, err)
			continue
		}

		logging.Info("%s: wrote %d bytes", newFile.Name(), i, n)
	}
}
