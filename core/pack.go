package buildpack

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/vcs"
	log "github.com/Sirupsen/logrus"
)

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}
	return
}

func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return
}

func DownloadVcs(dir string, url string, version string) {
	repo, err := vcs.NewRepo(url, dir)
	if err != nil {
		log.Fatal("Unable to check out repo: ", url)
	}

	err = repo.Get()
	if err != nil {
		log.Fatalf("Unable to checkout repo: ", url)
	}

	if version == "" {
		version = "master"
	}

	err = repo.UpdateVersion(version)
	if err != nil {
		log.Fatalf("Unable to checkout ref: ", version)
	}
}

func DownloadLocal(dir string, path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	packName := filepath.Base(path)

	dir, err = filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	target := filepath.Join(dir, packName)

	CopyDir(path, target)
}
