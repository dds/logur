package logur

// WithFields returns a new contextual logger instance with context added to it.
func WithFields(logger Logger, fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return logger
	}

	return &ContextualLogger{logger: logger, fields: fields}
}

type ContextualLogger struct {
	logger Logger
	fields map[string]interface{}
}

func (l *ContextualLogger) Trace(args ...interface{}) {
	l.logger.WithFields(l.fields).Trace(args...)
}

func (l *ContextualLogger) Debug(args ...interface{}) {
	l.logger.WithFields(l.fields).Debug(args...)
}

func (l *ContextualLogger) Info(args ...interface{}) {
	l.logger.WithFields(l.fields).Info(args...)
}

func (l *ContextualLogger) Warn(args ...interface{}) {
	l.logger.WithFields(l.fields).Warn(args...)
}

func (l *ContextualLogger) Error(args ...interface{}) {
	l.logger.WithFields(l.fields).Error(args...)
}

func (l *ContextualLogger) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return l
	}

	logger := l.logger

	// Do not add a new layer
	// Create a new logger instead with the parent fields
	if ctxlogger, ok := l.logger.(*ContextualLogger); ok && len(ctxlogger.fields) > 0 {
		_fields := make(map[string]interface{}, len(ctxlogger.fields)+len(fields))

		for key, value := range ctxlogger.fields {
			_fields[key] = value
		}

		for key, value := range fields {
			_fields[key] = value
		}

		fields = _fields
		logger = ctxlogger.logger
	}

	return &ContextualLogger{logger: logger, fields: fields}
}