package logger

import (
	"log"
	"sync"
	"time"

	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
)

type BufferedLogger struct {
	Writer   LogWriter
	Buffer   []db.CreateJobLogParams
	Mutex    sync.Mutex
	Interval time.Duration
	MaxSize  int
	quit     chan struct{}
}

func NewBufferedLogger(writer LogWriter, flushInterval time.Duration, maxBufferSize int) *BufferedLogger {
	bl := &BufferedLogger{
		Writer:   writer,
		Buffer:   make([]db.CreateJobLogParams, 0, maxBufferSize),
		Interval: flushInterval,
		MaxSize:  maxBufferSize,
		quit:     make(chan struct{}),
	}
	go bl.runFlusher()
	return bl
}

func (bl *BufferedLogger) Log(taskID, logLevel, message string) {
	bl.Mutex.Lock()
	jobId, err := utils.ParseUUID(taskID)
	if err != nil {
		log.Println(err)
	}
	bl.Buffer = append(bl.Buffer, db.CreateJobLogParams{
		JobID:    jobId,
		LogLevel: logLevel,
		Message:  message,
	})
	if len(bl.Buffer) >= bl.MaxSize {
		bufferCopy := bl.Buffer
		bl.Buffer = make([]db.CreateJobLogParams, 0, bl.MaxSize)
		go bl.Writer.WriteLogs(bufferCopy)
	}
	bl.Mutex.Unlock()
}

func (bl *BufferedLogger) runFlusher() {
	ticker := time.NewTicker(bl.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bl.Mutex.Lock()
			if len(bl.Buffer) > 0 {
				bufferCopy := bl.Buffer
				bl.Buffer = make([]db.CreateJobLogParams, 0, bl.MaxSize)
				go bl.Writer.WriteLogs(bufferCopy)
			}
			bl.Mutex.Unlock()
		case <-bl.quit:
			return
		}
	}
}

func (bl *BufferedLogger) Stop() {
	close(bl.quit)
	bl.Mutex.Lock()
	if len(bl.Buffer) > 0 {
		_ = bl.Writer.WriteLogs(bl.Buffer)
	}
	bl.Mutex.Unlock()
}
