//go:generate go run -tags generate gen.go

package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	//"fmt"
	"path/filepath"

	. "github.com/eyedeekay/go-ccw"
)

//go:embed i2pchrome.js/*
//go:embed localcdn/*
//go:embed onionbrowse/*
//go:embed scriptsafe/*
//go:embed ublockorigin/*
var extensionContent embed.FS

var EXTENSIONS = []string{
	"localcdn",
	"onionbrowse",
	"scriptsafe",
	"ublockorigin",
	"i2pchrome.js",
}

var EXTENSIONHASHES = []string{
	"a0ea5acc47bdfdf360f3e3ab2467b26ac84b2886f80aa92e0105bd79b17c1241",
	"8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85",
	"a674a9c883773a181bf2fb5340d014713044070e53d7b9c3789c0295b0b53fc6",
	"8d7ace5c493193c15286690189f81eb72dc50223f5ed185f43975f8c5cc10cf1",
	"986330c6ccc668102d7b187c808440557940316ee23a5b8b0ba2da1a6eb730a4",
}

func extensionPaths(outpath string) []string {
	var paths []string
	for _, extension := range EXTENSIONS {
		paths = append(paths, outpath+"/"+extension)
	}
	return paths
}

func WriteOutExtensions(outdir string) error {
	// Walk the contents of extensionContent and write the files out to disk
	os.MkdirAll(outdir, 0755)
	return fs.WalkDir(extensionContent, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		outpath := filepath.Join(outdir, path)
		if d.IsDir() {
			if err := os.MkdirAll(outpath, 0755); err != nil {
				//log.Println(err)
			}
		}
		bytes, err := extensionContent.ReadFile(path)
		if err != nil {
			log.Println(err)
		}
		if err := ioutil.WriteFile(outpath, bytes, 0644); err != nil {
			log.Println(err)
		}
		return nil
	})
}

var ARGS = []string{
	"--safebrowsing-disable-download-protection",
	"--disable-client-side-phishing-detection",
	"--disable-3d-apis",
	"--disable-accelerated-2d-canvas",
	"--disable-remote-fonts",
	"--disable-sync-preferences",
	"--disable-sync",
	"--disable-speech",
	"--disable-webgl",
	"--disable-reading-from-canvas",
	"--disable-gpu",
	"--disable-32-apis",
	"--disable-auto-reload",
	"--disable-background-networking",
	"--disable-d3d11",
	"--disable-file-system",
}

func main() {
	WriteOutExtensions("i2pchromium-browser")
	CHROMIUM, ERROR = SecureExtendedChromium("i2pchromium-browser", false, extensionPaths("i2pchromium-browser"), EXTENSIONHASHES, ARGS...)
	//CHROMIUM, ERROR = ExtendedChromium("i2pchromium-browser", false, extensionPaths("extensions"), ARGS...)
	if ERROR != nil {
		log.Fatal(ERROR)
	}
	defer CHROMIUM.Close()
	<-CHROMIUM.Done()
}
