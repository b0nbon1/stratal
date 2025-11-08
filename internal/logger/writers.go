package logger

import "io"

// logWriter implements io.Writer for job-level logging
type logWriter struct {
	logger *JobRunLogger
	stream string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	if w.stream == "stderr" {
		w.logger.LogJob(ErrorLevel, message, nil)
	} else {
		w.logger.LogJob(InfoLevel, message, nil)
	}
	return len(p), nil
}

// taskLogWriter implements io.Writer for task-level logging
type taskLogWriter struct {
	logger    *JobRunLogger
	stream    string
	taskRunID string
}

func (w *taskLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	if w.stream == "stderr" {
		w.logger.LogTask(w.taskRunID, ErrorLevel, message, w.stream, nil)
	} else {
		w.logger.LogTask(w.taskRunID, InfoLevel, message, w.stream, nil)
	}
	return len(p), nil
}

// GetWriter returns an io.Writer for job-level logging
func (jrl *JobRunLogger) GetWriter(stream string) io.Writer {
	return &logWriter{
		logger: jrl,
		stream: stream,
	}
}

// GetWriterForTaskRun returns an io.Writer for task-level logging
func (jrl *JobRunLogger) GetWriterForTaskRun(taskRunID, stream string) io.Writer {
	return &taskLogWriter{
		logger:    jrl,
		stream:    stream,
		taskRunID: taskRunID,
	}
}
