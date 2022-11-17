package logger

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type SrvLogger struct {
	Mode    string
	srvLogs *SrvLogs
}

type AlertStr struct {
	Message  string
	Domain   string
	Function string
	Err      error
}

type WarningStr struct {
	Message  string
	Domain   string
	Function string
	Kv       []string
}

type InfoStr struct {
	Message string
	Domain  string
	Kv      []string
}

type DebugStr struct {
	Message string
	Kv      []string
}

func NewSrvLogger(ctx context.Context, md string) *SrvLogger {

	a := make(chan AlertStr, 10)
	w := make(chan WarningStr, 10)
	i := make(chan InfoStr, 10)
	d := make(chan DebugStr, 10)

	srvLogs := &SrvLogs{
		ctx: ctx,
		A:   a,
		W:   w,
		I:   i,
		D:   d,
	}

	var mode string
	switch md {
	case "prod":
		mode = "prod"
	case "Development":
		mode = "development"
	case "debug":
		mode = "debug"
	default:
		mode = "prod"
	}

	go srvLogs.runLogger()

	return &SrvLogger{mode, srvLogs}
}

func (s *SrvLogger) Alert(a *AlertStr) {
	s.srvLogs.A <- *a
}

func (s *SrvLogger) Warning(w *WarningStr) {
	s.srvLogs.W <- *w
}

func (s *SrvLogger) Info(i *InfoStr) {
	if s.Mode != "prod" {
		s.srvLogs.I <- *i
	}
}

func (s *SrvLogger) Debug(d *DebugStr) {
	if s.Mode == "debug" {
		s.srvLogs.D <- *d
	}
}

type SrvLogs struct {
	ctx context.Context
	A   chan AlertStr
	W   chan WarningStr
	I   chan InfoStr
	D   chan DebugStr
}

func (s *SrvLogs) runLogger() {
	for {
		select {
		case e := <-s.A:
			zap.S().DPanicw("Alert", e.Message, e.Domain, e.Function, e.Err)
		case e := <-s.W:
			zap.S().Warnw("Warning", e.Message, e.Domain, e.Kv)
		case e := <-s.I:
			zap.S().Infow("Info", e.Message, e.Domain, e.Kv)
		case e := <-s.D:
			zap.S().Debugw("Debug", e.Message, e.Kv)
		case <-s.ctx.Done():
			// TODO unbuffer the channels ??
			return
		}
	}
}

// package logger
//
// import (
// 	"go.uber.org/zap"
// )
//
// type SrvLogger struct {
// 	Mode string
// 	Log  *SrvLogs
// }
//
// func NewSrvLogger(md string) *SrvLogger {
//
// 	var mode string
// 	switch md {
// 	case "prod":
// 		mode = "prod"
// 	case "development":
// 		mode = "development"
// 	case "debug":
// 		mode = "debug"
// 	default:
// 		mode = "prod"
// 	}
//
// 	return &SrvLogger{mode, &SrvLogs{}}
// }
//
// func (s *SrvLogger) Alert(msg, domain, function string, err error) {
// 	s.Log.alert(msg, domain, function, err)
// }
//
// func (s *SrvLogger) Warning(msg, domain, function string, kv ...string) {
// 	s.Log.warning(msg, domain, function, kv)
// }
//
// func (s *SrvLogger) Info(msg, domain string, kv ...string) {
// 	if s.Mode != "prod" {
// 		s.Log.info(msg, domain, kv)
// 	}
// }
//
// func (s *SrvLogger) Debug(msg string, kv ...string) {
// 	if s.Mode == "debug" {
// 		s.Log.debug(msg, kv)
// 	}
// }
//
// type SrvLogs struct{}
//
// func (s *SrvLogs) alert(msg, domain, function string, err error) {
// 	zap.S().DPanicw("Alert", msg, domain, function, err)
// }
//
// func (s *SrvLogs) warning(msg, domain, function string, kv []string) {
// 	zap.S().Warnw("Warning", msg, domain, kv)
// }
//
// func (s *SrvLogs) info(msg, domain string, kv []string) {
// 	zap.S().Infow("Info", msg, domain, kv)
// }
//
// func (s *SrvLogs) debug(msg string, kv []string) {
// 	zap.S().Debugw("Debug", msg, kv)
// }
