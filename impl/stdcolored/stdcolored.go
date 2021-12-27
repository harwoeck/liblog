package stdcolored

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func (l *logger) logColored(level Level, msg string, fields fieldCollection) {
	var s string

	t := time.Now()
	s += fmt.Sprintf("\033[0;37m%s\033[0m", t.Format("2006-01-02 15:04:"))
	s += fmt.Sprintf("\033[0;97m%s\033[0m", t.Format("05.999"))
	tz := t.Format("Z07:00")
	if len(tz) > 0 {
		s += fmt.Sprintf(" \033[0;37m%s\033[0m ", tz)
	}

	switch level {
	case DebugLevel:
		s += "\033[1;95mDEBUG\033[0m  " // purple
	case InfoLevel:
		s += "\033[1;94mINFO\033[0m   " // cyan
	case WarnLevel:
		s += "\033[1;103mWARN\033[0m   " // yellow background
	case ErrorLevel:
		s += "\033[1;101mERROR\033[0m  " // red background
	case DPanicLevel:
		s += "\033[1;101mDPANIC\033[0m " // red background
	case FatalLevel:
		s += "\033[1;101mFATAL\033[0m  " // red background
	}

	if len(l.name) > 0 {
		s += fmt.Sprintf("\033[0;37m%s\033[0m ", l.name)
	}

	_, frameF, frameL, defined := runtime.Caller(3)
	if !defined {
		_, _ = fmt.Fprintf(l.errWriter, "liblog.std: failed to get caller\n")
	} else {
		if len(l.wd) > 0 {
			frameF = strings.TrimPrefix(frameF, l.wd)
			frameF = strings.TrimPrefix(frameF, "/")
		}

		s += fmt.Sprintf("%s:%d ", frameF, frameL)
	}

	s += fmt.Sprintf("\033[0;97m%s\033[0m", msg)

	if len(fields) > 0 {
		s += fmt.Sprintf("\033[0;37m%s\033[0m", " (")

		for i, f := range fields {
			if i > 0 {
				s += ", "
			}

			s += fmt.Sprintf("\033[0;37m%s\033[0m", f.Key())
			s += fmt.Sprintf("\033[0;37m%s\033[0m", "=")

			var p string
			switch t := f.Value().(type) {
			case int:
				p = fmt.Sprintf("%d", t)
			case float64:
				p = fmt.Sprintf("%f", t)
			default:
				p = fmt.Sprintf("%q", fmt.Sprintf("%v", t))
			}

			s += fmt.Sprintf("\033[0;94m%s\033[0m", p)

		}

		s += fmt.Sprintf("\033[0;37m%s\033[0m", ")")
	}

	s += "\n"

	if _, err := fmt.Fprintf(l.outWriter, s); err != nil {
		_, _ = fmt.Fprintf(l.errWriter, "liblog.std: writing of message %q failed due to: %v", s, err)
	}
}
