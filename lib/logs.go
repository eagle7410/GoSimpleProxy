package lib

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logfile     = "development.log"
	logFileName string
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) GetStatusColorString() string {
	if rec.status < 300 {
		return fmt.Sprintf("[0;32m%v[39m", rec.status)
	}

	if rec.status < 500 {
		return fmt.Sprintf("[1;33m%v[39m", rec.status)
	}

	return fmt.Sprintf("[0;31m%v[39m", rec.status)
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func SetLogFileName(name string) {
	logfile = name
}

func OpenLogFile() {
	t := time.Now()
	oldFileName := logFileName

	logFileName = t.Format("2006-01-02") + "_" + logfile

	if FileExists("./logs") == false {
		err := os.Mkdir("./logs", 0777)

		if err != nil {
			log.Fatal("Mkdir logs error %s", err)
		}
	}

	if len(oldFileName) == 0 || oldFileName != logFileName {
		fmt.Printf("Logging to file %v\n", logFileName)

		lf, err := os.OpenFile("./logs/"+logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile: %s", err)
		}

		log.SetOutput(lf)
	}
}

const formatLog = "%v %s %s %s\n"

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		OpenLogFile()

		rec := statusRecorder{w, http.StatusOK}

		handler.ServeHTTP(&rec, r)

		defer func() {
			log.Printf(formatLog, rec.GetStatusColorString(), r.Method, r.URL, r.RemoteAddr)
			fmt.Printf(formatLog, rec.GetStatusColorString(), r.Method, r.URL, r.RemoteAddr)
		}()
	})
}

func LogAppRun(port string) {
	mask := "\n=== App run At (" + time.Now().Format("2006-01-02T15:04:05") + ") in http://localhost:" + port + " ===\n"
	fmt.Printf(mask)
	log.Printf(mask)
}

func LogFatalf(format string, v ...interface{}) {
	fr := "[0;31m" + format + "[39m"
	fmt.Printf(fr, v...)
	log.Fatalf(fr, v...)
}

func Logf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	log.Printf(format, v...)
}

func LogEF(format string, v ...interface{}) {
	Logf("[0;31m"+format+"[39m\n", v...)
}
