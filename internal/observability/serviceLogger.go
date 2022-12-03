package observability

import (
	"go.uber.org/zap"
)

type SrvLogger struct {
	Mode string
}

func NewSrvLogger(md string) *SrvLogger {

	var mode string
	switch md {
	case "prod":
		mode = "prod"
	case "development":
		mode = "development"
	case "debug":
		mode = "debug"
	default:
		mode = "prod"
	}

	return &SrvLogger{mode}
}

// display in any mode
func (s *SrvLogger) Error(msg string, err error) {
	zap.S().Errorw(msg, "error_report", err)
}

// display in any mode
func (s *SrvLogger) Warning(msg, various string) {
	zap.S().Warnw(msg, "various", various)
}

// display in debug and development mode only
func (s *SrvLogger) Info(msg, various string) {
	if s.Mode != "prod" {
		zap.S().Infow(msg, "various", various)
	}
}

// display in debug mode only(see setup() for the setting)
func (s *SrvLogger) Debug(msg, various string) {
	zap.S().Debugw(msg, "various", various)
}
