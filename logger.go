// Package logger allows the abstraction of logging within a go program.
package logger

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"
)

// Type Level is the representation of the different logging levels.
type Level int

const (
	DEBUG Level = 1 // Messages aimed only to the developers.
	INFO  Level = 2 // Indicates the state of the program (without errors).
	WARN  Level = 3 // Detected problem which doesn't disturb the service itself.
	ERROR Level = 4 // Significant problem which disturb the service.
	FATAL Level = 5 // Problem which makes the service stop working totaly.
)

// formatedJSON is a struct containing the base representation of the JSON
// output. It is filled by the logger itself accordingly to the various
// function calls. It contains:
//    - Process: Name of the current program using the logger.
//    - Timestamp: Time of the logged event.
//    - Level: Level of the logged event.
//    - Client: ID of the Sewan client running this program.
//    - Data: Generic placeholder for what needs to be logged.
type formatedJSON struct {
	Process   string                 `json:"process"`
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Client    string                 `json:"client"`
	Data      map[string]interface{} `json:"data"`
}

// Logger contains all the variables needed by the logger:
//    - level: Minimum level which will be logged, lower levels are simply
//             discarded.
//    - logger: JSON encoder which will be used to log.
//    - json: base struct which will be logged (more infos above).
type Logger struct {
	level  Level
	logger *json.Encoder
	json   formatedJSON
}

// levelToString returns a map allowing the conversion of a Level into a string
// Usage exemple:
//    `levelToString()[INFO]`
// will return the string "INFO".
func levelToString() map[Level]string {
	var lvlstrings = make(map[Level]string)
	lvlstrings[DEBUG] = "DEBUG"
	lvlstrings[INFO] = "INFO"
	lvlstrings[WARN] = "WARN"
	lvlstrings[ERROR] = "ERROR"
	lvlstrings[FATAL] = "FATAL"
	return lvlstrings
}

// stringToLevel returns a map allowing the conversion of a string to a level
// Usage exemple:
//    `stringToLevel()["FATAL"]`
// will return a Level typed FATAL variable.
func stringToLevel() map[string]Level {
	var stringlvls = make(map[string]Level)
	stringlvls["DEBUG"] = DEBUG
	stringlvls["INFO"] = INFO
	stringlvls["WARN"] = WARN
	stringlvls["ERROR"] = ERROR
	stringlvls["FATAL"] = FATAL
	stringlvls[""] = INFO // Case where env var is empty, set log level to INFO.
	return stringlvls
}

// stringToWriter returns an io.Writer based on a string.
// Add a case here if you want to add a new type of writer to the logger.
func stringToWriter(writer string) io.Writer {
	switch writer {
	case "stdout":
		return os.Stdout
	default:
		return os.Stdout
		// var f, err = os.OpenFile("/var/log/"+os.Args[0]+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	panic(err)
		// }
		// return f
	}
}

// New creates, initializes and returns a new Logger.
// Parameters of the function are:
//    - level: string representing the minimum logging level.
//    - writer: string representing the selected writer.
//    - prettyprint: string representing a bool, activating the pretty print or
//                   not (`true` or `false`).
// Note that the function protoptype has been made so you can use it with
// environment variables, like that:
//    `var lg = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_WRITER"), os.Getenv("LOG_PRETTY"))`
// LOG_LEVEL, LOG_WRITER and LOG_PRETTY can be set just before the program, or
// exported in the environment:
//    `LOG_LEVEL=DEBUG LOG_WRITER=stdout, LOG_PRETTY=true ./my_program`
// OR
//    `export LOG_LEVEL=WARN`
//    `./my_program` (which will be run with a WARN log level)
func New(level string, writer string, prettyprint string) Logger {
	var lg Logger
	lg.level = stringToLevel()[level]
	lg.logger = json.NewEncoder(stringToWriter(writer))
	if prettyprint == "true" {
		lg.logger.SetIndent("", "  ")
	}

	var process = strings.Split(os.Args[0], "/")
	lg.json.Process = process[len(process)-1]

	return lg
}

// Log writes and outputs a given data *v* to a given Level *level*.
func (lg Logger) Log(level Level, v ...interface{}) {
	if level >= lg.level {
		lg.json.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		lg.json.Level = levelToString()[level]
		lg.json.Data = make(map[string]interface{})
		var n = len(v)
		for i := 0; i < n; i = i + 2 {
			if str, ok := v[i].(string); ok {
				lg.json.Data[str] = v[i+1]
			} else {
				panic("json key must be a string")
			}
		}

		lg.logger.Encode(lg.json)
	}
}
