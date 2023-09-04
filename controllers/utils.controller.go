package controllers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/adhupraba/breadit-server/utils"
)

type UtilsController struct{}

type imageData struct {
	Url string `json:"url"`
}

type metaObj struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       imageData `json:"image"`
}

type urlMetadataRes struct {
	Success int     `json:"success"`
	Meta    metaObj `json:"meta"`
}

func (uc *UtilsController) GetUrlMetadata(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	fmt.Println("get url metadata for =>", url)

	if url == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid href")
		return
	}

	res, err := http.Get(url)

	if err != nil {
		utils.RespondWithError(w, res.StatusCode, err.Error())
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	bodyStr := string(body)

	title := ""
	titleRe := regexp.MustCompile(`<title>(.*?)<\/title>`)

	if titleSlice := titleRe.FindStringSubmatch(bodyStr); len(titleSlice) >= 2 {
		title = titleSlice[1]
	}

	description := ""
	descriptionRe := regexp.MustCompile(`<meta name="description" content="(.*?)"`)

	if descriptionSlice := descriptionRe.FindStringSubmatch(bodyStr); len(descriptionSlice) >= 2 {
		description = descriptionSlice[1]
	}

	imageUrl := ""
	imageRe := regexp.MustCompile(`<meta property="og:image" content="(.*?)"`)

	if imageSlice := imageRe.FindStringSubmatch(bodyStr); len(imageSlice) >= 2 {
		imageUrl = imageSlice[1]
	}

	utils.RespondWithJsonDirect(w, 200, urlMetadataRes{
		Success: 1,
		Meta: metaObj{
			Title:       title,
			Description: description,
			Image:       imageData{Url: imageUrl},
		},
	})
}
