package log4go

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RecordLogWriter struct {
	rec chan *LogRecord
	rot chan bool

	// The opened file
	filename string
	file     *os.File

	// Rotate at linecount
	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	daily_opendate int

	// Keep old logfiles (.001, .002, etc)
	rotate    bool
	maxbackup int
}

// This is the FileLogWriter's output method
func (w *RecordLogWriter) LogWrite(rec *LogRecord) {
	w.rec <- rec
}

func (w *RecordLogWriter) Close() {
	close(w.rec)
	w.file.Sync()
}



// NewRecordLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewRecordLogWriter(fname string, rotate bool) *RecordLogWriter {

	createFileLog(fname)

	w := &RecordLogWriter{
		rec:       make(chan *LogRecord, LogBufferLength),
		rot:       make(chan bool),
		filename:  fname,
		rotate:    false,
		maxbackup: 999,
	}

	// open the file for the first time
	if err := w.intRotate(); err != nil {
		Error("FileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}
	w.SetRotate(rotate)

	go func() {
		defer func() {
			if w.file != nil {
				w.file.Close()
			}
		}()

		for {
			select {
			case <-w.rot:
				if err := w.intRotate(); err != nil {
					Error("NewRecordLogWriter(%q): %s\n", w.filename, err)
					return
				}
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				now := time.Now()
				if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
					(w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
					(w.daily && now.Day() != w.daily_opendate) {
					if err := w.intRotate(); err != nil {
						Error("NewRecordLogWriter(%q): %s\n", w.filename, err)
						return
					}
				}

				data, err := json.Marshal(rec.Infos)
				if err != nil {
					Error("json errorr(%q): %s",w.filename,err)
					return
				}
				// Perform the write
				n, err := fmt.Fprint(w.file, string(data),"\n")
				if err != nil {
					Error("NewRecordLogWriter(%q): %s\n", w.filename, err)
					return
				}

				// Update the counts
				w.maxlines_curlines++
				w.maxsize_cursize += n
			}
		}
	}()

	return w
}


// Request that the logs rotate
func (w *RecordLogWriter) Rotate() {
	w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *RecordLogWriter) intRotate() error {
	// Close any log file that may be open
	if w.file != nil {
		w.file.Close()
	}

	// If we are keeping log files, move it to the next available number
	if w.rotate {
		_, err := os.Lstat(w.filename)
		if err == nil { // file exists
			// Find the next available number
			num := 1
			fname := ""
			if w.daily && time.Now().Day() != w.daily_opendate {
				yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

				fname = w.filename + fmt.Sprintf(".%s", yesterday)
				_, err = os.Lstat(fname)

				if err == nil{
					for ; err == nil && num <= w.maxbackup; num++ {
						fname = w.filename + fmt.Sprintf(".%s.%03d", yesterday, num)
						_, err = os.Lstat(fname)
					}
				}

				// return error if the last file checked still existed
				if err == nil {
					return Error("Rotate: Cannot find free log number to rename %s\n", w.filename)
				}
			} else {
				num = w.maxbackup - 1
				for ; num >= 1; num-- {
					fname = w.filename + fmt.Sprintf(".%d", num)
					nfname := w.filename + fmt.Sprintf(".%d", num+1)
					_, err = os.Lstat(fname)
					if err == nil {
						os.Rename(fname, nfname)
					}
				}
			}

			w.file.Close()
			// Rename the file to its newfound home
			err = os.Rename(w.filename, fname)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}
	}

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.file = fd

	now := time.Now()

	// Set the daily open date to the current date
	w.daily_opendate = now.Day()

	// initialize rotation values
	w.maxlines_curlines = 0
	w.maxsize_cursize = 0

	return nil
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *RecordLogWriter) SetRotateLines(maxlines int) *RecordLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
	w.maxlines = maxlines
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *RecordLogWriter) SetRotateSize(maxsize int) *RecordLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *RecordLogWriter) SetRotateDaily(daily bool) *RecordLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// Set max backup files. Must be called before the first log message
// is written.
func (w *RecordLogWriter) SetRotateMaxBackup(maxbackup int) *RecordLogWriter {
	w.maxbackup = maxbackup
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *RecordLogWriter) SetRotate(rotate bool) *RecordLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}