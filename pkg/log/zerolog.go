package log

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type logger struct {
	log zerolog.Logger
}

type loggerEvent struct {
	event *zerolog.Event
}

func NewZeroLog(programName, programVersion string, level Level) *logger {
	var zerolevel = zerolog.ErrorLevel
	if level != "" {
		switch level {
		case Error:
			zerolevel = zerolog.ErrorLevel
		case Warn:
			zerolevel = zerolog.WarnLevel
		case Info:
			zerolevel = zerolog.InfoLevel
		case Debug:
			zerolevel = zerolog.DebugLevel
		}
	}
	zerolog.SetGlobalLevel(zerolevel)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	ctx := zerolog.New(os.Stdout).With().Timestamp()

	ctx = ctx.Dict("program", zerolog.Dict().Fields(map[string]interface{}{
		"name":    programName,
		"version": programVersion,
	}))

	return &logger{ctx.Logger()}
}

func (l *logger) Fatal() LoggerEvent {
	return &loggerEvent{l.log.Fatal()}
}

func (l *logger) Error() LoggerEvent {
	return &loggerEvent{l.log.Error()}
}

func (l *logger) Warn() LoggerEvent {
	return &loggerEvent{l.log.Warn()}
}

func (l *logger) Info() LoggerEvent {
	return &loggerEvent{l.log.Info()}
}

func (l *logger) Debug() LoggerEvent {
	return &loggerEvent{l.log.Debug()}
}

func (le *loggerEvent) Trace(ID string) LoggerEvent {
	le.event = le.event.Str("traceId", ID)
	return le
}

func (le *loggerEvent) Org(clientID, applicationID string) LoggerEvent {
	le.event = le.event.Dict("org", zerolog.Dict().Fields(map[string]interface{}{
		"clientId":      clientID,
		"applicationId": applicationID,
	}))
	return le
}

func (le *loggerEvent) Req(ID, IP, host, scheme, method, URL, body string, headers map[string]string) LoggerEvent {
	req := zerolog.Dict().Fields(map[string]interface{}{
		"id":      ID,
		"ip":      IP,
		"host":    host,
		"scheme":  scheme,
		"method":  method,
		"url":     URL,
		"headers": headers,
	})

	if body != "" {
		req.Str("body", body)
	}

	le.event = le.event.Dict("req", req)
	return le
}

func (le *loggerEvent) Res(status int, elapsedTime time.Duration, body string, bodyByteLength int, headers map[string]string) LoggerEvent {
	res := zerolog.Dict().Fields(map[string]interface{}{
		"status":         status,
		"elapsedTime":    elapsedTime,
		"bodyByteLength": bodyByteLength,
		"headers":        headers,
	})

	if body != "" {
		res.Str("body", body)
	}

	le.event = le.event.Dict("res", res)
	return le
}

func (le *loggerEvent) Err(err error) LoggerEvent {
	return le.ErrWithStack(err, "")
}

func (le *loggerEvent) ErrWithStack(err error, stacktrace string) LoggerEvent {
	if err == nil {
		return le
	}

	fields := map[string]interface{}{
		"message": err.Error(),
	}

	if stacktrace != "" {
		fields["stacktrace"] = stacktrace
	}

	le.event = le.event.Dict("error", zerolog.Dict().Fields(fields))
	return le
}

func (le *loggerEvent) Send(message string) {
	le.event.Msg(message)
}

func (le *loggerEvent) Sendf(message string, args ...interface{}) {
	if len(args) == 0 {
		le.Send(message)
		return
	}

	le.event.Msg(fmt.Sprintf(message, args...))
}
