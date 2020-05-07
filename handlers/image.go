package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"album.com/messagequeue"
	"album.com/models"
)

// InsertImage insert img under the album
func InsertImage(w http.ResponseWriter, r *http.Request) {

	albumName := strings.TrimSpace(r.URL.Query().Get("album"))
	//album name mandatory
	if albumName == "" {
		e := &result{
			Message: "ERROR: Missing mandatory query params[album]",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	//check whether album exists
	if err := os.Mkdir(getPath(albumName), 0755); os.IsNotExist(err) {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: album[%s] not exist", albumName),
			Code:    404,
		}
		renderERROR(w, e)
		return

	}

	imgFile, head, e := r.FormFile("image")
	if e != nil {
		defer imgFile.Close()
		e := &result{
			Message: "ERROR: Invalid query params, image file missing",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	if _, err := os.Stat(getPath(albumName + "/" + head.Filename)); err == nil {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: %s image already exists in %s", head.Filename, albumName),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	contentType := head.Header.Get("Content-Type")
	log.Println(contentType)
	if ok := strings.HasPrefix(contentType, "image/"); !ok {
		e := &result{
			Message: "ERROR: Invalid image file format",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	data, err := ioutil.ReadAll(imgFile)
	if err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: while reading %s image in %s", head.Filename, getPath(albumName)),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	err = ioutil.WriteFile(getPath(albumName+"/"+head.Filename), data, 0666)
	if e != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unable to insert image[%s] to the  album[%s]", head.Filename, albumName),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	res := &result{
		Message: fmt.Sprintf("Image[%s] uploaded successfully to the album[%s]", head.Filename, albumName),
		Code:    200,
	}

	// Notify the changes to client
	messagequeue.Notify(res.Message)

	renderJSON(w, http.StatusOK, res)
}

// DeleteImage deletes the specific image within an album
func DeleteImage(w http.ResponseWriter, r *http.Request) {

	albumName := strings.TrimSpace(r.URL.Query().Get("album"))
	imageName := strings.TrimSpace(r.URL.Query().Get("name"))

	//album name mandatory
	if albumName == "" {
		e := &result{
			Message: "ERROR: Invalid query params",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	//Album exists or not
	if _, err := os.Stat(getPath(albumName)); os.IsNotExist(err) {
		e := &result{
			Message: fmt.Sprintf("ERROR: album[%s] not exist", albumName),
			Code:    404,
		}
		renderERROR(w, e)
		return
	}

	//image exists or not
	if _, err := os.Stat(getPath(albumName + "/" + imageName)); err != nil {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: %s image doesn't exists in %s", imageName, albumName),
			Code:    404,
		}
		renderERROR(w, e)
		return
	}

	if err := os.Remove(getPath(albumName + "/" + imageName)); err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unable to remove image[%s] from the  album[%s]", imageName, albumName),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	res := &result{
		Message: fmt.Sprintf("Image[%s] deleted successfully from the album[%s]", imageName, albumName),
		Code:    200,
	}

	// Notify the changes to client
	messagequeue.Notify(res.Message)

	renderJSON(w, http.StatusOK, res)
}

// GetImage fetch the image from the requested album
func GetImage(w http.ResponseWriter, r *http.Request) {

	albumName := strings.TrimSpace(r.URL.Query().Get("album"))
	imgName := strings.TrimSpace(r.URL.Query().Get("name"))

	//album, img name mandatory
	if albumName == "" || imgName == "" {
		e := &result{
			Message: "ERROR: Mandatory query params",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	//Existence of album
	if _, err := os.Stat(getPath(albumName)); os.IsNotExist(err) {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: album[%s] not exist", albumName),
			Code:    404,
		}
		renderERROR(w, e)
		return
	}

	files, _ := ioutil.ReadDir(getPath(albumName))
	fmt.Println(files)
	found := false
	for _, entry := range files {
		fmt.Println(" ", entry.Name(), entry.IsDir())
		if strings.EqualFold(entry.Name(), imgName) {
			found = true
		}
	}

	if !found {
		e := &result{
			Message: fmt.Sprintf("ERROR: Image[%s] not found in the album[%s]", imgName, albumName),
			Code:    404,
		}
		renderERROR(w, e)
		return
	}

	img, err := os.Open(getPath(albumName + "/" + imgName))
	if err != nil {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: Unbale to open imgae[%s] in the album[%s], err:%q", imgName, albumName, err),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}
	defer img.Close()

	buffer := make([]byte, 512)
	if _, err := img.Read(buffer); err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unbale to open imgae[%s] in the album[%s], err:%q", imgName, albumName, err),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	byts, err := ioutil.ReadFile(getPath(albumName + "/" + imgName))
	if err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unable to read image[%s] from the album[%s]", imgName, albumName),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}
	w.Write(byts)

}

// GetImages fetch all the image from the album
func GetImages(w http.ResponseWriter, r *http.Request) {

	albumName := strings.TrimSpace(r.URL.Query().Get("album"))
	//album, img name mandatory
	if albumName == "" {
		e := &result{
			Message: "ERROR: Mandatory query params",
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	//Existence of album
	if _, err := os.Stat(getPath(albumName)); os.IsNotExist(err) {
		log.Println(err)
		e := &result{
			Message: fmt.Sprintf("ERROR: album[%s] not exist", albumName),
			Code:    404,
		}
		renderERROR(w, e)
		return
	}

	images, err := ioutil.ReadDir(getPath(albumName))
	if err != nil {
		e := &result{
			Message: fmt.Sprintf("ERROR: Unable to get Images from the album[%s]", albumName),
			Code:    400,
		}
		renderERROR(w, e)
		return
	}

	albumIns := models.Album{}
	albumIns.AlbumName = albumName
	for _, img := range images {
		albumIns.Images = append(albumIns.Images, img.Name())
	}

	renderJSON(w, http.StatusOK, albumIns)
}
