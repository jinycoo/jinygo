// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package log_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"jinycoo.com/jinygo/log"
	"jinycoo.com/jinygo/log/zapcore"
)

func Example_presets() {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	logger := log.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()

	const url = "http://example.com"


	// In the unusual situations where every microsecond matters, use the
	// Logger. It's even faster than the SugaredLogger, but only supports
	// structured logging.
	logger.Info("Failed to fetch URL.",
		// Structured context as strongly typed fields.
		log.String("url", url),
		log.Int("attempt", 3),
		log.Duration("backoff", time.Second),
	)
	// Output:
	// {"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
	// {"level":"info","msg":"Failed to fetch URL: http://example.com"}
	// {"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
}

func Example_basicConfiguration() {
	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg log.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	// Output:
	// {"level":"info","message":"logger construction succeeded","foo":"bar"}
}

func Example_advancedConfiguration() {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := log.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := log.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(log.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(log.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := log.New(core)
	defer logger.Sync()
	logger.Info("constructed a logger")
}

func ExampleNamespace() {
	logger := log.NewExample()
	defer logger.Sync()

	logger.With(
		log.Namespace("metrics"),
		log.Int("counter", 1),
	).Info("tracked some metrics")
	// Output:
	// {"level":"info","msg":"tracked some metrics","metrics":{"counter":1}}
}

func ExampleNewStdLog() {
	logger := log.NewExample()
	defer logger.Sync()

	std := log.NewStdLog(logger)
	std.Print("standard logger wrapper")
	// Output:
	// {"level":"info","msg":"standard logger wrapper"}
}

func ExampleRedirectStdLog() {
	logger := log.NewExample()
	defer logger.Sync()

	undo := log.RedirectStdLog(logger)
	defer undo()
	// Output:
	// {"level":"info","msg":"redirected standard library"}
}

func ExampleReplaceGlobals() {
	logger := log.NewExample()
	defer logger.Sync()

	undo := log.ReplaceGlobals(logger)
	defer undo()

	log.L().Info("replaced zap's global loggers")
	// Output:
	// {"level":"info","msg":"replaced zap's global loggers"}
}

func ExampleAtomicLevel() {
	atom := log.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := log.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := log.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	logger.Info("info logging enabled")

	atom.SetLevel(log.ErrorLevel)
	logger.Info("info logging disabled")
	// Output:
	// {"level":"info","msg":"info logging enabled"}
}

func ExampleAtomicLevel_config() {
	// The zap.Config struct includes an AtomicLevel. To use it, keep a
	// reference to the Config.
	rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	var cfg log.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("info logging enabled")

	cfg.Level.SetLevel(log.ErrorLevel)
	logger.Info("info logging disabled")
	// Output:
	// {"level":"info","message":"info logging enabled"}
}

func ExampleLogger_Check() {
	logger := log.NewExample()
	defer logger.Sync()

	if ce := logger.Check(log.DebugLevel, "debugging"); ce != nil {
		// If debug-level log output isn't enabled or if zap's sampling would have
		// dropped this log entry, we don't allocate the slice that holds these
		// fields.
		ce.Write(
			log.String("foo", "bar"),
			log.String("baz", "quux"),
		)
	}

	// Output:
	// {"level":"debug","msg":"debugging","foo":"bar","baz":"quux"}
}

func ExampleLogger_Named() {
	logger := log.NewExample()
	defer logger.Sync()

	// By default, Loggers are unnamed.
	logger.Info("no name")

	// The first call to Named sets the Logger name.
	main := logger.Named("main")
	main.Info("main logger")

	// Additional calls to Named create a period-separated path.
	main.Named("subpackage").Info("sub-logger")
	// Output:
	// {"level":"info","msg":"no name"}
	// {"level":"info","logger":"main","msg":"main logger"}
	// {"level":"info","logger":"main.subpackage","msg":"sub-logger"}
}

func ExampleWrapCore_replace() {
	// Replacing a Logger's core can alter fundamental behaviors.
	// For example, it can convert a Logger to a no-op.
	nop := log.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewNopCore()
	})

	logger := log.NewExample()
	defer logger.Sync()

	logger.Info("working")
	logger.WithOptions(nop).Info("no-op")
	logger.Info("original logger still works")
	// Output:
	// {"level":"info","msg":"working"}
	// {"level":"info","msg":"original logger still works"}
}

func ExampleWrapCore_wrap() {
	// Wrapping a Logger's core can extend its functionality. As a trivial
	// example, it can double-write all logs.
	doubled := log.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, c)
	})

	logger := log.NewExample()
	defer logger.Sync()

	logger.Info("single")
	logger.WithOptions(doubled).Info("doubled")
	// Output:
	// {"level":"info","msg":"single"}
	// {"level":"info","msg":"doubled"}
	// {"level":"info","msg":"doubled"}
}
