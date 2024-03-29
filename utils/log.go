package utils

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"go-admin-starter/utils/config"
	"go-admin-starter/utils/file"
	"path"
	"time"
)

var Log *logrus.Logger

var c = config.New()

func LogSetup() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)

	logPath := conf.App.LogSavePath
	file.IsNotExistMkDir(logPath)

	ConfigLocalFilesystemLogger(logPath, conf.App.LogSaveName, time.Hour*24*30, time.Second*60*60*24-1)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+"%Y%m%d"+conf.App.LogFileExt,
		//rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		Log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})

	Log.AddHook(lfHook)
}
