package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	name         = "benplg"
	ig_name      = "ben_plg"
	baseFolder   = "./" + ig_name + "/:tagged"
	assetsFolder = "../assets/images/" + name + "/"
	photoFile    = "../data/" + name + "/photos.yaml"
)

var photosToDelete = make([]string, 0)

type photo struct {
	Url string `yaml:"url"`
}

func updatePhotosToDelete(name string) {
	max := len(photosToDelete) - 1
	for k := range photosToDelete {
		if photosToDelete[max-k] == name {
			photosToDelete = append(photosToDelete[:max-k], photosToDelete[max-k+1:]...)
		}
	}
}

func main() {
	files, _ := ioutil.ReadDir(assetsFolder)
	for _, f := range files {
		if f.Name() == ".gitignore" {
			continue
		}
		photosToDelete = append(photosToDelete, f.Name())
	}

	// pswd := os.Getenv("IGP")
	// uname := os.Getenv("IGU")

	//	cmd := exec.Command("instaloader", "--tagged", "--login", uname, "-p", pswd, "--no-videos", "--no-captions", "--no-metadata-json", "--no-profile-pic", ig_name)
	cmd := exec.Command("instaloader", "--tagged", "--no-videos", "--no-captions", "--no-metadata-json", "--no-profile-pic", ig_name)
	fmt.Println(`exec.Command("instaloader", "--tagged", "--no-videos", "--no-captions", "--no-metadata-json", "--no-profile-pic", ` + ig_name + ")")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		if e := cmd.Run(); e != nil {
			if e.Error() != "signal: killed" {
				panic(e)
			}
		}
	}()

	time.Sleep(50 * time.Second)
	cmd.Process.Kill()

	files, _ = ioutil.ReadDir(baseFolder)

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	limit := 12
	if limit > len(files) {
		limit = len(files)
	}
	filenames := make([]photo, limit)
	files = files[:limit]

	for k, f := range files {
		filenames[k] = photo{Url: f.Name()}
		updatePhotosToDelete(f.Name())
	}

	for _, f := range photosToDelete {
		if err := os.Remove(assetsFolder + f); err != nil {
			fmt.Println("An error occured L86:", err)
		}
	}

	for _, f := range filenames {
		b, _ := ioutil.ReadFile(baseFolder + "/" + f.Url)
		if err := ioutil.WriteFile(assetsFolder+f.Url, b, 0755); err != nil {
			fmt.Println("An error occured L93:", err)
		}
	}

	b, _ := yaml.Marshal(filenames)
	if err := ioutil.WriteFile(photoFile, b, 0755); err != nil {
		fmt.Println("An error occured L99:", err)
	}
}
