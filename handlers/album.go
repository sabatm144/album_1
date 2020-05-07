package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"album.com/messagequeue"

	"path/filepath"
)

const (
	gallery = "gallery"
)

func getPath(path string) string {
	return filepath.Join(gallery, path)
}

// InsertAlbum creates an album  directory
func InsertAlbum(w http.ResponseWriter, r *http.Request) {

	albumName := strings.TrimSpace(r.URL.Query().Get("name"))
	// log.Println(albumName)
	//album name mandatory
	if albumName == "" {
		e := &result{
			Message: "ERROR: Album name should not be empty",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	//Check album directory exists
	if err := os.Mkdir(getPath(albumName), 0755); os.IsExist(err) {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: Album[%s] already exists, err:%q", albumName, err),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	res := &result{
		Message: fmt.Sprintf("Album[%s] created successfully", albumName),
		Code:    200,
	}

	// Notify the changes to client
	messagequeue.Notify(res.Message)

	renderJSON(w, http.StatusOK, res)
}

// DeleteAlbum deletes the specific album
func DeleteAlbum(w http.ResponseWriter, r *http.Request) {

	albumName := strings.TrimSpace(r.URL.Query().Get("name"))
	log.Println(albumName)
	//album name mandatory
	if albumName == "" {
		e := &result{
			Message: "ERROR: Album name should not be empty",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	folderPath := getPath(albumName)

	//Check whether album  exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		e := &result{
			Message: fmt.Sprintf("ERROR: Album[%s]  doesn't exists, err:%q", albumName, err),
			Code:    404,
		}
		renderERROR(w, e)
		return
	}

	if err := removeAllFiles(folderPath); err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unable to delete Album[%s], err:%q", albumName, err),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	if err := os.Remove(folderPath); err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unable to delete Album[%s], err:%q", albumName, err),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	res := &result{
		Message: fmt.Sprintf("Album[%s] deleted successfully.", albumName),
		Code:    200,
	}

	// Notify the changes to client
	messagequeue.Notify(res.Message)

	renderJSON(w, http.StatusOK, res)
}

// removeAll files in the album
func removeAllFiles(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
