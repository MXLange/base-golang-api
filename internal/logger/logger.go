package logger

import (
	"context"
	"fmt"
)

type LoggerIF interface {
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, args ...any)
	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, format string, args ...any)
	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, format string, args ...any)
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, format string, args ...any)
	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, format string, args ...any)
	Panic(ctx context.Context, msg string)
	Panicf(ctx context.Context, format string, args ...any)
}

type Logger struct {
}

func NewLogger() LoggerIF {
	return &Logger{}
}

func (l *Logger) Info(ctx context.Context, msg string) {
	// Implement your logging logic here
	fmt.Println(msg)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...any) {
	// Implement your logging logic here
	fmt.Printf(format, args...)
	fmt.Println()
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	// Implement your logging logic here
	fmt.Println(msg)
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...any) {
	// Implement your logging logic here
	fmt.Printf(format, args...)
	fmt.Println()
}

func (l *Logger) Error(ctx context.Context, msg string) {
	// Implement your logging logic here
	fmt.Println(msg)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...any) {
	// Implement your logging logic here
	fmt.Printf(format, args...)
	fmt.Println()
}

func (l *Logger) Debug(ctx context.Context, msg string) {
	// Implement your logging logic here
	fmt.Println(msg)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...any) {
	// Implement your logging logic here
	fmt.Printf(format, args...)
	fmt.Println()
}

func (l *Logger) Fatal(ctx context.Context, msg string) {
	// Implement your logging logic here
	fmt.Println(msg)
}

func (l *Logger) Fatalf(ctx context.Context, format string, args ...any) {
	// Implement your logging logic here
	fmt.Printf(format, args...)
	fmt.Println()
}

func (l *Logger) Panic(ctx context.Context, msg string) {
	// Implement your logging logic here
	fmt.Println(msg)
}

func (l *Logger) Panicf(ctx context.Context, format string, args ...any) {
	// Implement your logging logic here
	fmt.Printf(format, args...)
	fmt.Println()
}
