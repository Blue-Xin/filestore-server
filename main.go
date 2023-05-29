package main

import (
	"filestore-serve/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	//http.HandleFunc("/file/query", handler.Fil)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
