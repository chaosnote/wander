package utils

import (
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CustomField 任意欄位類型
type CustomField map[string]any

// encodeTimeExcludingTimezone 自定義時間編碼器，排除時區資訊
func encodeTimeExcludingTimezone(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("01-02 15:04:05.000"))
}

// 自定義日誌等級編碼器 (已註解)
// func levelEncoder(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
// 	encoder.AppendString("[" + strings.ToUpper(level.String()[:1]) + "]") // 只取首字母，並加上中括號
// }

func NewEmptyLogger() *zap.Logger {
	return zap.NewNop()
}

func NewConsoleLogger(callerSkip int) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.LevelKey = zapcore.OmitKey
	// config.EncoderConfig.MessageKey = zapcore.OmitKey
	config.EncoderConfig.EncodeTime = encodeTimeExcludingTimezone
	// config.EncoderConfig.TimeKey = zapcore.OmitKey
	// config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.Encoding = "json"

	logger, err := config.Build(zap.AddCallerSkip(callerSkip))
	if err != nil {
		panic(err)
	}
	return logger
}

func NewFileLogger(dir string, callerSkip int) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   path.Join(dir, ".log"), // 文件輸出路徑
		MaxSize:    10,                     // 文件最大大小 (MB)
		LocalTime:  true,                   // 使用本地時間
		Compress:   false,                  // 是否壓縮檔案
		MaxAge:     30,                     // 舊檔案保留天數
		MaxBackups: 50,                     // 最多備份檔案數量
	}
	writeSyncer := zapcore.AddSync(&hook)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = zapcore.OmitKey
	// encoderConfig.MessageKey = zapcore.OmitKey
	encoderConfig.EncodeTime = encodeTimeExcludingTimezone
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 短檔案路徑

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // zapcore.NewConsoleEncoder(encoderConfig),
		writeSyncer,
		zapcore.InfoLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip))
}
