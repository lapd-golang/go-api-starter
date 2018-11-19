package upload

import (
	"admin-server/pkg/config"
	"admin-server/pkg/file"
	"admin-server/pkg/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"
)

func GetImageFullUrl(path string) string {
	return config.Conf.App.ImagePrefixUrl + "/" + path
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	var newFileName string

	newName := com.UrlEncode(name)

	nowTime := time.Now()
	s := nowTime.Unix() * 1000

	dt := nowTime.UnixNano() / int64(time.Millisecond)
	dt = dt - s
	dateRd := util.GenerateRangeNum(0, int(dt))
	newFileName = util.EncodeMD5(newName + string(dateRd))

	return newFileName + ext
}

func GetImagePath() string {
	return config.Conf.App.ImageSavePath
}

func CheckImageType(fileBytes []byte) bool {
	fileType := http.DetectContentType(fileBytes)

	for _, t := range config.Conf.App.ImageAllowTypes {
		if t == fileType {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		util.Log.Warn(err)
		return false
	}

	return size <= config.Conf.App.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func SaveImage(file multipart.File, savePath string) error {
	savePath = config.Conf.App.RuntimeRootPath + savePath

	fileBytes, _ := ioutil.ReadAll(file)
	defer file.Close()

	if ! CheckImageType(fileBytes) || ! CheckImageSize(file) {
		return errors.New("校验图片错误，图片格式或大小有问题")
	}

	err := CheckImage(savePath)
	if err != nil {
		util.Log.Warn(err)
		return errors.New("检查图片失败")
	}

	out, err := os.Create(savePath)
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(fileBytes))
	if err != nil {
		util.Log.Fatal(err)
		return errors.New("上传图片失败")
	}

	return nil
}
