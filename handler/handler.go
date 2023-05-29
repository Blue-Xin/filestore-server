package handler

import (
	"encoding/json"
	"filestore-serve/meta"
	"filestore-serve/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internet server err")
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流存储到目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to get data:", err.Error())
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "C:\\Users\\Ferry\\Desktop\\filestore-serve\\tmp\\",
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location + header.Filename)
		if err != nil {
			fmt.Println("Failed to create file :", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("Failed to save data into file", err.Error())
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Println(fileMeta.FileSha1)
		//存入文件信息
		meta.UpdateFileMeta(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSuccessHandler:上传已完成
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}

// 获取文件信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fileMeta := meta.GetFileMeta(filehash)
	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(filehash)
	open, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer open.Close()
	data, err := ioutil.ReadAll(open)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.FormValue("filehash")
	fileName := r.FormValue("filename")
	op := r.FormValue("opType")
	if op != "0" {
		w.WriteHeader(http.StatusForbidden)
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fileMeta := meta.GetFileMeta(fileHash)
	fmt.Println(fileMeta)
	fileMeta.FileName = fileName
	meta.UpdateFileMeta(fileMeta)

	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.FormValue("filehash")
	fmt.Println(fileHash)
	//查询有没有
	fileMeta := meta.GetFileMeta(fileHash)
	file := fileMeta.Location + fileMeta.FileName
	os.Remove(file)
	meta.RemoveFileMeta(fileHash)
	w.WriteHeader(http.StatusOK)
}
