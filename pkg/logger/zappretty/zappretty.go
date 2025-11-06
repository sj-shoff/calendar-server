package zappretty

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Err returns a zap field for error
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Field returns a zap field with key-value pair
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

type PrettyEncoderOptions struct {
	TimeZone *time.Location
}

type PrettyEncoder struct {
	zapcore.Encoder
	pool     buffer.Pool
	timeZone *time.Location
}

func (opts PrettyEncoderOptions) NewPrettyEncoder() zapcore.Encoder {
	timezone := opts.TimeZone
	if timezone == nil {
		timezone = time.Local
	}

	return &PrettyEncoder{
		pool:     buffer.NewPool(),
		timeZone: timezone,
	}
}

func (e *PrettyEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := e.pool.Get()

	// Format time
	localTime := entry.Time.In(e.timeZone)
	timeStr := localTime.Format("2006-01-02T15:04:05.000-0700")
	buf.AppendString(timeStr)

	// Format level with colors
	var levelStr string
	switch entry.Level {
	case zapcore.InfoLevel:
		levelStr = colorBrightGreen("INFO")
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		levelStr = colorBrightRed("ERROR")
	case zapcore.WarnLevel:
		levelStr = colorYellow("WARN")
	case zapcore.DebugLevel:
		levelStr = colorMagenta("DEBUG")
	default:
		levelStr = entry.Level.CapitalString()
	}

	// Add level with padding
	buf.AppendString("\t" + levelStr + "\t")

	// Caller (gray color)
	if entry.Caller.Defined {
		caller := entry.Caller.TrimmedPath()
		buf.AppendString(colorGray(caller) + "\t")
	}

	// Message
	msg := entry.Message
	if entry.Level == zapcore.ErrorLevel || entry.Level == zapcore.FatalLevel || entry.Level == zapcore.PanicLevel {
		msg = colorBrightRed(msg) // Red for errors
	} else {
		msg = colorCyan(msg) // Cyan for others
	}
	buf.AppendString(msg)

	// Fields (JSON)
	if len(fields) > 0 {
		fieldMap := make(map[string]interface{})

		tempEncoder := zapcore.NewMapObjectEncoder()
		for _, field := range fields {
			field.AddTo(tempEncoder)
		}

		for key, value := range tempEncoder.Fields {
			fieldMap[key] = value
		}

		if len(fieldMap) > 0 {
			jsonBytes, err := json.Marshal(fieldMap)
			if err != nil {
				return nil, err
			}

			fieldStr := string(jsonBytes)
			buf.AppendString("\t" + fieldStr)
		}
	}

	buf.AppendString("\n")

	return buf, nil
}

func (e *PrettyEncoder) Clone() zapcore.Encoder {
	return &PrettyEncoder{
		pool:     e.pool,
		timeZone: e.timeZone,
	}
}

type PrettyCore struct {
	zapcore.Core
	enc *PrettyEncoder
	out io.Writer
}

func NewPrettyCore(writer io.Writer, level zapcore.LevelEnabler, opts PrettyEncoderOptions) zapcore.Core {
	encoder := opts.NewPrettyEncoder()

	return &PrettyCore{
		Core: zapcore.NewCore(
			encoder,
			zapcore.AddSync(writer),
			level,
		),
		enc: encoder.(*PrettyEncoder),
		out: writer,
	}
}

func (c *PrettyCore) With(fields []zap.Field) zapcore.Core {
	encoderClone := c.enc.Clone().(*PrettyEncoder)
	return &PrettyCore{
		Core: c.Core.With(fields),
		enc:  encoderClone,
		out:  c.out,
	}
}

func (c *PrettyCore) Write(entry zapcore.Entry, fields []zap.Field) error {
	buf, err := c.enc.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}
	defer buf.Free()

	_, err = c.out.Write(buf.Bytes())
	return err
}

func (c *PrettyCore) Sync() error {
	if syncer, ok := c.out.(zapcore.WriteSyncer); ok {
		return syncer.Sync()
	}
	return nil
}

// NewPrettyLogger creates a new pretty logger with Debug level
func NewPrettyLogger(writer io.Writer, opts PrettyEncoderOptions) *zap.Logger {
	core := NewPrettyCore(writer, zapcore.DebugLevel, opts)
	return zap.New(core, zap.AddCaller())
}

// SetupLogger creates a default pretty logger
func SetupLogger() *zap.Logger {
	return NewPrettyLogger(os.Stdout, PrettyEncoderOptions{
		TimeZone: time.Local,
	})
}

// SetupLoggerWithLevel creates a pretty logger with specified level
func SetupLoggerWithLevel(level zapcore.Level) *zap.Logger {
	return NewPrettyLoggerWithLevel(os.Stdout, level, PrettyEncoderOptions{
		TimeZone: time.Local,
	})
}

// NewPrettyLoggerWithLevel creates a pretty logger with specified level
func NewPrettyLoggerWithLevel(writer io.Writer, level zapcore.LevelEnabler, opts PrettyEncoderOptions) *zap.Logger {
	core := NewPrettyCore(writer, level, opts)
	return zap.New(core, zap.AddCaller())
}

// Color functions
func colorMagenta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func colorYellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}

func colorBrightRed(s string) string {
	return "\033[91m" + s + "\033[0m"
}

func colorBrightGreen(s string) string {
	return "\033[92m" + s + "\033[0m"
}

func colorCyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}

func colorGray(s string) string {
	return "\033[90m" + s + "\033[0m"
}
