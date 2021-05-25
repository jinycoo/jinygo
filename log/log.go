/**------------------------------------------------------------**
 * @filename zap/log.go
 * @author   jinycoo - jinycoo@jinycoo.com
 * @version  1.0.0
 * @date     2021/3/15 下午3:44
 * @desc     jinycoo.com-zap-log: summary
 **------------------------------------------------------------**/
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"jinycoo.com/jinygo/utils/go.uuid"

	"jinycoo.com/jinygo/log/zapcore"
)

const (
	LOG = "log"

	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	WARN  = "warn"

	DEFAULT = "default"
	CONSOLE = "console"
	JSON    = "json"
)

var (
	Coo           *CooLogger
	cooConf       *CooLogConfig
	encoderConfig zapcore.EncoderConfig
)

type CooLogger struct {
	core        zapcore.Core
	development bool
	name        string
	errorOutput zapcore.WriteSyncer
	addCaller   bool
	addStack    zapcore.LevelEnabler
	callerSkip  int
	appName     string
}

type CooLogConfig struct {
	Dev     bool              `json:"dev" yaml:"dev"`
	Level   string            `json:"level" yaml:"level"`
	Encoder []string          `json:"encoder" yaml:"encoder"`
	Encode  map[string]string `json:"encode" yaml:"encode"`
	Key     map[string]string `json:"key" yaml:"key"`
	OutPuts []string          `json:"outputs" yaml:"outputs"`
}

func defaultLogConfig() {
	cooConf = &CooLogConfig{
		Dev:     true,
		Level:   DEBUG,
		Encoder: []string{CONSOLE},
		Encode:  map[string]string{"time": "local", "level": "capital", "duration": "string", "caller": "short"},
		Key: map[string]string{
			"name":       "logger",
			"time":       "time",
			"level":      "level",
			"caller":     "caller",
			"message":    "msg",
			"stacktrace": "stacktrace",
		},
		OutPuts: []string{"stderr", DEFAULT},
	}
}

func (clog *CooLogConfig) lvlEncoder() {
	lvl := clog.Encode["level"]
	switch lvl {
	case "capital":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	case "capitalColor":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "color":
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	default:
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
}

func (clog *CooLogConfig) timeEncoder() {
	time := clog.Encode["time"]
	switch time {
	case "local":
		encoderConfig.EncodeTime = logEncodeTime
	case "iso8601", "ISO8601":
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case "millis":
		encoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	case "nanos":
		encoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
	default:
		encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	}
}

func (clog *CooLogConfig) durEncoder() {
	dur := clog.Encode["duration"]
	switch dur {
	case "string":
		encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	case "nanos":
		encoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
	default:
		encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	}
}

func (clog *CooLogConfig) callerEncoder() {
	caller := clog.Encode["caller"]
	switch caller {
	case "full":
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	default:
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}
}

func NewLogger(app string) *CooLogger {

	encoderConfig.NameKey = cooConf.Key["name"]
	encoderConfig.TimeKey = cooConf.Key["time"]
	encoderConfig.LevelKey = cooConf.Key["level"]
	encoderConfig.CallerKey = cooConf.Key["caller"]
	encoderConfig.MessageKey = cooConf.Key["message"]
	encoderConfig.StacktraceKey = cooConf.Key["stacktrace"]

	encoderConfig.LineEnding = zapcore.DefaultLineEnding

	cooConf.timeEncoder()
	cooConf.lvlEncoder()
	cooConf.durEncoder()
	cooConf.callerEncoder()

	var lvl AtomicLevel
	switch cooConf.Level {
	case DEBUG:
		lvl = NewAtomicLevelAt(zapcore.DebugLevel)
	case WARN:
		lvl = NewAtomicLevelAt(zapcore.WarnLevel)
	case ERROR:
		lvl = NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		lvl = NewAtomicLevelAt(zapcore.InfoLevel)
	}
	var outputs []string
	for _, p := range cooConf.OutPuts {
		if p == DEFAULT {
			logfile := fmt.Sprintf("%s_%s.log", app, time.Now().Format("2006-01-02_15-04-05"))
			outputs = append(outputs, logfile)
		} else {
			if strings.HasPrefix(p, "./") {
				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
				root := strings.Replace(dir, "\\", "/", -1)
				cpath := filepath.Join(root, strings.Replace(p, "./", "", 1))
				outputs = append(outputs, cpath)
			} else {
				outputs = append(outputs, p)
			}
		}
	}
	sink, _, _ := Open(outputs...)

	var cores []zapcore.Core
	for _, e := range cooConf.Encoder {
		switch e {
		case CONSOLE:
			consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
			cores = append(cores, zapcore.NewCore(consoleEncoder, sink, lvl))
		case JSON:
			jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
			cores = append(cores, zapcore.NewCore(jsonEncoder, sink, lvl))
		}
	}

	//errSink, _, err := zap.Open("stderr")
	//if err != nil {
	//	closeOut()
	//}

	log := &CooLogger{
		core:        zapcore.NewTee(cores...),
		development: cooConf.Dev,
		errorOutput: zapcore.Lock(os.Stderr),
		addStack:    zapcore.FatalLevel + 1,
		addCaller:   cooConf.Dev,
		appName:     app,
	}
	return log
}

func Init(appName, mode, filew string) {
	if mode == "dev" {
		defaultLogConfig()
	} else {
		cooConf = &CooLogConfig{
			Level:   INFO,
			Dev: false,
			Encoder: []string{JSON},
			Encode:  map[string]string{"time": "iso8601", "level": "lowercase", "duration": "string", "caller": "short"},
			Key: map[string]string{
				"name":       "logger",
				"time":       "time",
				"level":      "level",
				"caller":     "caller",
				"message":    "msg",
				"stacktrace": "stacktrace",
			},
			OutPuts: []string{"stderr", filew},
		}
	}
	Coo = NewLogger(appName)
}

func (clog *CooLogger) check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	const callerSkipOffset = 2
	ent := zapcore.Entry{
		LoggerName: Coo.name + uuid.Must(uuid.NewV4(), nil).String(),
		Time:       time.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce := Coo.core.Check(ent, nil)
	willWrite := ce != nil

	switch ent.Level {
	case zapcore.PanicLevel:
		ce = ce.Should(ent, zapcore.WriteThenPanic)
	case zapcore.FatalLevel:
		ce = ce.Should(ent, zapcore.WriteThenFatal)
	case zapcore.DPanicLevel:
		if Coo.development {
			ce = ce.Should(ent, zapcore.WriteThenPanic)
		}
	}

	if !willWrite {
		return ce
	}

	ce.ErrorOutput = Coo.errorOutput
	if Coo.addCaller {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(Coo.callerSkip + callerSkipOffset))
		if !ce.Entry.Caller.Defined {
			fmt.Fprintf(Coo.errorOutput, "%v Logger.check error: failed to get caller\n", time.Now().Local())
			Coo.errorOutput.Sync()
		}
	}
	if Coo.addStack.Enabled(ce.Entry.Level) {
		ce.Entry.Stack = Stack("").String
	}

	return ce
}
func Debug(details ...interface{}) {
	if ce := Coo.check(zapcore.DebugLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Info(details ...interface{}) {
	if ce := Coo.check(zapcore.InfoLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Infof(template string, args ...interface{}) {
	if ce := Coo.check(zapcore.InfoLevel, sprint(template, args...)); ce != nil {
		ce.Write()
	}
}

func Warn(details ...interface{}) {
	if ce := Coo.check(zapcore.WarnLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Warnf(template string, args ...interface{}) {
	if ce := Coo.check(zapcore.WarnLevel, sprint(template, args...)); ce != nil {
		ce.Write()
	}
}

func ZError(details ...interface{}) {
	if ce := Coo.check(zapcore.ErrorLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Errorf(template string, args ...interface{}) {
	if ce := Coo.check(zapcore.ErrorLevel, sprint(template, args...)); ce != nil {
		ce.Write()
	}
}

func DPanic(details ...interface{}) {
	if ce := Coo.check(zapcore.DPanicLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Panic(details ...interface{}) {
	if ce := Coo.check(zapcore.PanicLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Fatal(details ...interface{}) {
	if ce := Coo.check(zapcore.FatalLevel, fmt.Sprint(details...)); ce != nil {
		ce.Write()
	}
}

func Fatalf(template string, args ...interface{}) {
	if ce := Coo.check(zapcore.FatalLevel, sprint(template, args...)); ce != nil {
		ce.Write()
	}
}

func Sync() error {
	return Coo.core.Sync()
}

func logEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05] "))
}

func sprint(template string, args ...interface{}) (message string) {
	message = template
	argsLen := len(args)
	if message == "" && argsLen > 0 {
		message = fmt.Sprint(args...)
	} else if message != "" && argsLen > 0 {
		message = fmt.Sprintf(template, args...)
	}
	return
}
