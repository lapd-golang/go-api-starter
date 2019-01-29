package upload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"go-admin-starter/utils"
	"go-admin-starter/utils/config"
	"go-admin-starter/utils/file"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"
)

var conf = config.New()

func GetImageFullUrl(path string) string {
	return conf.App.ImagePrefixUrl + "/" + path
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	var newFileName string

	newName := com.UrlEncode(name)

	nowTime := time.Now()
	s := nowTime.Unix() * 1000

	dt := nowTime.UnixNano() / int64(time.Millisecond)
	dt = dt - s
	dateRd := utils.GenerateRangeNum(0, int(dt))
	newFileName = utils.EncodeMD5(newName + string(dateRd))

	return newFileName + ext
}

func GetImagePath() string {
	return conf.App.ImageSavePath
}

func CheckImageType(fileBytes []byte) bool {
	fileType := http.DetectContentType(fileBytes)

	for _, t := range conf.App.ImageAllowTypes {
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
		utils.Log.Warn(err)
		return false
	}

	return size <= conf.App.ImageMaxSize
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

func SaveImage(f multipart.File, path string, imageName string) error {
	fileBytes, _ := ioutil.ReadAll(f)
	defer f.Close()

	if ! CheckImageType(fileBytes) || ! CheckImageSize(f) {
		return errors.New("校验图片错误，图片格式或大小有问题")
	}

	err := CheckImage(path)
	if err != nil {
		utils.Log.Warn(err)
		return errors.New("检查图片失败")
	}

	savePath := path + imageName

	out, err := os.Create(savePath)
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(fileBytes))
	if err != nil {
		utils.Log.Fatal(err)
		return errors.New("上传图片失败")
	}

	return nil
}
