package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level string
	Encoding string
	OutputPaths []string
	ErrorOutputPaths []string
	FileRotation bool
	Filename     string
	MaxSize      int  
	MaxBackups   int  
	MaxAge       int  
	Compress     bool 
	Development bool
	TimeEncoder zapcore.TimeEncoder
}

var (
	zapLogger   *zap.Logger
	sugar       *zap.SugaredLogger
	initialized = false
)

func Init(cfg Config) error {
	if initialized {
		_ = Sync()
		zapLogger = nil
		sugar = nil
		initialized = false
	}

	if cfg.Encoding == "" {
		if cfg.Development {
			cfg.Encoding = "console"
		} else {
			cfg.Encoding = "json"
		}
	}

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}

	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     cfg.TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if encoderCfg.EncodeTime == nil {
		encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// ISO8601-ish format
			enc.AppendString(t.UTC().Format(time.RFC3339Nano))
		}
	}

	var encoder zapcore.Encoder
	if strings.EqualFold(cfg.Encoding, "console") {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	// Build write syncers
	var syncers []zapcore.WriteSyncer

	// Always include stdout as a default sink (so logs appear in containers)
	syncers = append(syncers, zapcore.AddSync(os.Stdout))

	for _, p := range cfg.OutputPaths {
		lower := strings.ToLower(p)
		switch lower {
		case "stdout":
			// already added
		case "stderr":
			syncers = append(syncers, zapcore.AddSync(os.Stderr))
		default:
			// treat as file path
			f, ferr := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
			if ferr != nil {
				return fmt.Errorf("failed to open output path %s: %w", p, ferr)
			}
			syncers = append(syncers, zapcore.AddSync(f))
		}
	}

	if cfg.FileRotation && cfg.Filename != "" {
		if cfg.MaxSize == 0 {
			cfg.MaxSize = 100 // sensible default
		}
		if cfg.MaxBackups == 0 {
			cfg.MaxBackups = 7
		}
		if cfg.MaxAge == 0 {
			cfg.MaxAge = 30
		}
		l := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		syncers = append(syncers, zapcore.AddSync(l))
	} else if cfg.Filename != "" && !cfg.FileRotation {
		// if FileRotation is false but a filename is provided, open file without rotation
		f, ferr := os.OpenFile(cfg.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if ferr != nil {
			return fmt.Errorf("failed to open file %s: %w", cfg.Filename, ferr)
		}
		syncers = append(syncers, zapcore.AddSync(f))
	}

	// Combine syncers into one core sink
	var core zapcore.Core
	if len(syncers) == 1 {
		core = zapcore.NewCore(encoder, syncers[0], level)
	} else {
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(syncers...), level)
	}

	// Options
	opts := []zap.Option{
		zap.AddCaller(),      // include caller info
		zap.AddCallerSkip(1), // adjust for wrapper functions
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	zapLogger = zap.New(core, opts...)
	sugar = zapLogger.Sugar()
	initialized = true

	return nil
}

// Sync flushes any buffered logs. It is safe to call multiple times.
func Sync() error {
	if sugar != nil {
		_ = sugar.Sync() // sugar.Sync delegates to underlying logger
	}
	if zapLogger != nil {
		return zapLogger.Sync()
	}
	return nil
}

func Logger() (*zap.Logger, error) {
	if !initialized {
		if err := Init(Config{}); err != nil {
			return nil, err
		}
	}
	return zapLogger, nil
}

func Sugar() (*zap.SugaredLogger, error) {
	if !initialized {
		if err := Init(Config{}); err != nil {
			return nil, err
		}
	}
	return sugar, nil
}

func parseLevel(l string) (zapcore.LevelEnabler, error) {
	if l == "" {
		return zapcore.InfoLevel, nil
	}
	switch strings.ToLower(strings.TrimSpace(l)) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "dpanic":
		return zapcore.DPanicLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		var zl zapcore.Level
		if err := zl.UnmarshalText([]byte(l)); err == nil {
			return zl, nil
		}
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", l)
	}
}