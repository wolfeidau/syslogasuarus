package main

import (
	"fmt"
	"time"

	"log/syslog"

	"github.com/jeromer/syslogparser"
	"github.com/wolfeidau/syslogasuarus/syslogd"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

func main() {

	channel := make(chan syslogparser.LogParts, 1)

	svr := syslogd.NewServer()

	svr.ListenUDP(":10514")

	svr.Start(channel)

	for {
		logparts := <-channel

		fmt.Printf("%s %s %s %s %s\n", ts(logparts["timestamp"]), severity(logparts["severity"]), logparts["hostname"], logparts["tag"], logparts["content"])

	}

}

func ts(s interface{}) string {

	if s == nil {
		return time.Now().Format(time.RFC3339)
	}

	return s.(time.Time).Local().Format(time.RFC3339)

}

func severity(s interface{}) string {

	if s == nil {
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m ", blue, "UNKN")
	}

	sev := syslog.Priority(s.(int))

	switch sev {
	case syslog.LOG_EMERG:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", red, "EMERG")
	case syslog.LOG_ALERT:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", red, "ALERT")
	case syslog.LOG_CRIT:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", red, "CRIT")
	case syslog.LOG_ERR:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", red, "ERROR")
	case syslog.LOG_WARNING:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", yellow, "WARN")
	case syslog.LOG_NOTICE:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", blue, "NOTIC")
	case syslog.LOG_INFO:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", blue, "INFO")
	case syslog.LOG_DEBUG:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", green, "DEBUG")
	default:
		return fmt.Sprintf("\x1b[%dm[%-5s]\x1b[0m", blue, "UNKN")
	}

}
