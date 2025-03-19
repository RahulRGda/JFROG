// nolint:golint,errcheck,staticcheck
package logging

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	"path/filepath"

	"github.com/RahulRGda/jfrog/pkg/readenv"
	"github.com/astaxie/beego/logs"
	"github.com/francoispqt/onelog"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Alogger struct {
}

var Environment, Product_id, Application_id string
var filepathenv, logsToInclude string
var maxSize, maxFiles int
var maxDays int
var compress bool
var errorFileLogger, infoFileLogger, warnFileLogger, debugFileLogger, fatalFileLogger, oneLoglogger *onelog.Logger
var logLevelsToWrite map[string]bool
var Logger *onelog.Logger

func getCurrentFileDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("could not get caller info")
	}
	return filepath.Dir(filename), nil
}

// Initialize all standard fields that is coming from app conf.
func InitLogger() {
	currentDir, err := getCurrentFileDir()
	if err != nil {
		Panic("Error getting current directory: %v", err)
	}
	// Navigate up one directory.
	projectRoot1st := filepath.Dir(currentDir)
	projectRoot := filepath.Dir(projectRoot1st)

	// Construct the absolute path to the .env file.
	envPath := filepath.Join(projectRoot, ".env")
	err1 := godotenv.Load(envPath)
	if err1 != nil {
		if len(err1.Error()) > 200 {
			Panic("Error reading .env err:", err1.Error()[:200])
		} else {
			Panic("Error reading .env err:", err1.Error())
		}
	}
	Environment = readenv.GetEnvString("runmode")
	Application_id = readenv.GetEnvString("appname")
	Product_id = readenv.GetEnvString("application_id")

	RegisterLogger()
	getLogTypeMapping(logsToInclude)

	f := func(e onelog.Entry) {
		e.String("timestamp", time.Now().Format("2006-01-02T15:04:05.000"))
	}

	var parentLogger = onelog.New(os.Stdout, onelog.ALL).Hook(f)

	oneLoglogger = parentLogger.WithContext("")

	var parentfileLoggerError = onelog.New(&lumberjack.Logger{
		Filename:   filepathenv, // File path the logs has to stored
		MaxSize:    maxSize,     // Megabites in int
		MaxBackups: maxFiles,    // Total number of files
		MaxAge:     maxDays,     // Max number of days the file has to present.
		Compress:   compress,    // Disabled by default, bool compresses files?.
	}, onelog.ERROR).Hook(f)
	var parentfileLoggerDebug = onelog.New(&lumberjack.Logger{
		Filename:   filepathenv, // File path the logs has to stored
		MaxSize:    maxSize,     // Megabites in int
		MaxBackups: maxFiles,    // Total number of files
		MaxAge:     maxDays,     // Max number of days the file has to present.
		Compress:   compress,    // Disabled by default, bool compresses files?.
	}, onelog.DEBUG).Hook(f)
	var parentfileLoggerWarn = onelog.New(&lumberjack.Logger{
		Filename:   filepathenv, // File path the logs has to stored
		MaxSize:    maxSize,     // Megabites in int
		MaxBackups: maxFiles,    // Total number of files
		MaxAge:     maxDays,     // Max number of days the file has to present.
		Compress:   compress,    // Disabled by default, bool compresses files?.
	}, onelog.WARN).Hook(f)
	var parentfileLoggerInfo = onelog.New(&lumberjack.Logger{
		Filename:   filepathenv, // File path the logs has to stored
		MaxSize:    maxSize,     // Megabites in int
		MaxBackups: maxFiles,    // Total number of files
		MaxAge:     maxDays,     // Max number of days the file has to present.
		Compress:   compress,    // Disabled by default, bool compresses files?.
	}, onelog.INFO).Hook(f)
	var parentfileLoggerFatal = onelog.New(&lumberjack.Logger{
		Filename:   filepathenv, // File path the logs has to stored
		MaxSize:    maxSize,     // Megabites in int
		MaxBackups: maxFiles,    // Total number of files
		MaxAge:     maxDays,     // Max number of days the file has to present.
		Compress:   compress,    // Disabled by default, bool compresses files?.
	}, onelog.FATAL).Hook(f)

	errorFileLogger = parentfileLoggerError.WithContext("")
	infoFileLogger = parentfileLoggerInfo.WithContext("")
	warnFileLogger = parentfileLoggerWarn.WithContext("")
	debugFileLogger = parentfileLoggerDebug.WithContext("")
	fatalFileLogger = parentfileLoggerFatal.WithContext("")
}

func getLogTypeMapping(logLevelConfig string) {

	logLevelsToWrite = make(map[string]bool)

	allLogLevels := []string{"INFO", "DEBUG", "ERROR", "WARN", "FATAL"}

	var logLevelstoInclude []string
	for i, logLevel := range allLogLevels {
		if logLevel == strings.ToUpper(logLevelConfig) {
			logLevelstoInclude = allLogLevels[i:]
			break
		}
	}

	for _, loglevel := range logLevelstoInclude {
		logLevelsToWrite[loglevel] = true
	}
}

func RegisterLogger() {

	filepathenv = readenv.GetEnvString("onelog_path", "../onelog/errors.log")
	maxSize = readenv.GetEnvInt("onelog_maxsize", 2000)
	maxFiles = readenv.GetEnvInt("onelog_maxfiles", 3)
	maxDays = readenv.GetEnvInt("onelog_maxdays", 1)
	compress = readenv.GetEnvBool("onelog_compress", false)
	logsToInclude = readenv.GetEnvString("onelog_include_log_types", "Error")
}

func Info(message string, args ...interface{}) {
	message, standardFieldsMap, extraFieldsMap := getArgs(message, args...)
	LogWithJSON(INFO, message, standardFieldsMap, extraFieldsMap)

}

func Warn(message string, args ...interface{}) {
	message, standardFieldsMap, extraFieldsMap := getArgs(message, args...)
	LogWithJSON(WARN, message, standardFieldsMap, extraFieldsMap)
}

func Error(message string, args ...interface{}) {
	message, standardFieldsMap, extraFieldsMap := getArgs(message, args...)
	LogWithJSON(ERROR, message, standardFieldsMap, extraFieldsMap)
}

func Panic(message string, args ...interface{}) {
	message, standardFieldsMap, extraFieldsMap := getArgs(message, args...)
	LogWithJSON(WARN, message, standardFieldsMap, extraFieldsMap)
	panic(message)
}

func Fatal(message string, args ...interface{}) {
	message, standardFieldsMap, extraFieldsMap := getArgs(message, args...)
	LogWithJSON(FATAL, message, standardFieldsMap, extraFieldsMap)
}

func Debug(message string, args ...interface{}) {
	message, standardFieldsMap, extraFieldsMap := getArgs(message, args...)
	LogWithJSON(DEBUG, message, standardFieldsMap, extraFieldsMap)
}

func getArgs(msg string, args ...interface{}) (message string, standardFieldsMap StandardFields, extraFieldsMap ExtraFieldsMap) {
	argTemp := make([]interface{}, 0)
	argTemp = append(argTemp, msg)
	for _, arg := range args {
		switch t := arg.(type) {
		case StandardFields:
			standardFieldsMap = t
		case ExtraFieldsMap:
			extraFieldsMap = t
		default:
			// panic("Unknown argument")
			argTemp = append(argTemp, t)
		}
	}

	// emulate beego log behaviour, i.e. Append all the arguments in a message other than standard & extra fields
	message = strings.TrimSpace(formatLog(generateFmtStr(len(argTemp)), argTemp...))
	return message, standardFieldsMap, extraFieldsMap
}

// Logger Utils
func LogWithJSON(level LoggingLevel, message string, standardFieldsMap StandardFields, extraFieldsMap ExtraFieldsMap) {

	var ok bool
	var line int
	var source string
	var pc uintptr

	// This determines the calling method name
	// so we know what method called and also the line
	if strings.EqualFold(string(level), "warn") {
		pc, source, line, ok = runtime.Caller(6)
	} else {
		pc, source, line, ok = runtime.Caller(2)
	}

	if !ok {
		source = "???"
		line = 0
	}
	_, filename := path.Split(source)

	details := runtime.FuncForPC(pc)
	functionName := details.Name()
	if extraFieldsMap["function_name"] != nil && extraFieldsMap["line"] != nil {
		functionName = extraFieldsMap["function_name"].(string)
		line = extraFieldsMap["line"].(int)
	}

	standardFieldsMap.Environment = Environment
	standardFieldsMap.Product_id = Product_id
	standardFieldsMap.Application_id = Application_id

	if extraFieldsMap == nil {
		extraFieldsMap = ExtraFieldsMap{}
	}
	extraFieldsMap["source"] = filename
	extraFieldsMap["line"] = line
	// extraFieldsMap["function_name"] = source

	switch {
	case strings.EqualFold(string(level), "debug"):
		oneLoglogger.DebugWithFields(
			prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
		)
		if logLevelsToWrite["DEBUG"] {
			debugFileLogger.DebugWithFields(
				prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
			)
		}
	case strings.EqualFold(string(level), "warn"):
		oneLoglogger.WarnWithFields(
			prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
		)
		if logLevelsToWrite["WARN"] {
			warnFileLogger.WarnWithFields(
				prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
			)
		}
	case strings.EqualFold(string(level), "error"):
		oneLoglogger.ErrorWithFields(
			prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
		)
		if logLevelsToWrite["ERROR"] {
			errorFileLogger.ErrorWithFields(
				prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
			)
		}
	case strings.EqualFold(string(level), "fatal"):
		oneLoglogger.FatalWithFields(
			prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
		)
		if logLevelsToWrite["FATAL"] {
			fatalFileLogger.FatalWithFields(
				prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
			)
		}
	default:
		oneLoglogger.InfoWithFields(
			prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
		)
		if logLevelsToWrite["INFO"] {
			infoFileLogger.InfoWithFields(
				prepareMessage(message, functionName), func(e onelog.Entry) { logFields(standardFieldsMap, extraFieldsMap, e) },
			)
		}
	}
}

func SetLogger(adaptername string, config string) error {
	return logs.SetLogger(adaptername, config)
}

func logFields(standardFieldsMap StandardFields, extraFieldsMap ExtraFieldsMap, e onelog.Entry) {

	v := reflect.ValueOf(standardFieldsMap)
	typeOfS := v.Type()

	// We are extending extrafields to encompass standard fields
	// we will rely only on extra fields map for all our needs
	for i := 0; i < v.NumField(); i++ {
		if extraFieldsMap == nil {
			extraFieldsMap = ExtraFieldsMap{}
		}
		// we are assuming integer ids arent of a value 0 but this is debatable
		if v.Field(i).Interface() != "" && v.Field(i).Interface() != 0 {
			// this could overwrite a key with same name in extraFieldsMap
			// but we are expecting keys to be of distinct names across standard fields and extra fields
			attributeName := strings.ToLower(typeOfS.Field(i).Name)
			attributeValue := v.Field(i).Interface()
			extraFieldsMap[attributeName] = attributeValue
		}
	}
	for key, element := range extraFieldsMap {
		switch v := element.(type) {
		default:
			// this must be captured and logged,
			// the unhandled type must be identified and implemented
			e.Err(key, fmt.Errorf("unexpected type %T for %s", v, key))
		case int:
			e.Int(key, element.(int))
		case int64:
			e.Int64(key, element.(int64))
		case float64:
			e.Float(key, element.(float64))
		case float32:
			e.Float(key, element.(float64))
		case bool:
			e.Bool(key, element.(bool))
		case error:
			e.Err(key, element.(error))
		case string:
			e.String(key, fmt.Sprint(v))
		case time.Time:
			e.String(key, fmt.Sprint(v))
		}
	}
}

// The prepare message method will help us bring standardization
// in how we want to prepare the message key in our json log
func prepareMessage(message string, functionName string) string {
	var messageString string

	// messageString += functionName
	// messageString += " "
	messageString += message

	return messageString
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch t := f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			// format string
		} else {
			// do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		fmt.Print("type is", t)
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
