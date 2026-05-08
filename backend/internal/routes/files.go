package routes

import (
	"backend/internal/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func handlerFileUpload(w http.ResponseWriter, r *http.Request) {
	senderId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	privateIdStr := r.PathValue("private_id")
	privateId, err := strconv.ParseInt(privateIdStr, 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid private_id", nil)
		return
	}

	err = r.ParseMultipartForm(50 << 20) // 50 MB max
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Failed to parse multipart form", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Failed to retrieve file", nil)
		return
	}
	defer file.Close()

	dirPath := filepath.Join("files", "chats", fmt.Sprintf("%d", privateId), fmt.Sprintf("%d", senderId))
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Failed to create directory", nil)
		return
	}

	filePath := filepath.Join(dirPath, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Failed to create file", nil)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Failed to save file", nil)
		return
	}

	fileURL := fmt.Sprintf("/files/chats/%d/%d/%s", privateId, senderId, header.Filename)

	utils.JSON(w, http.StatusOK, true, "File uploaded successfully", map[string]any{
		"file_url": fileURL,
	})

}

func handlerGetFile() http.Handler {
	fs := http.FileServer(http.Dir("files"))
	return http.StripPrefix("/api/files/", fs)
}