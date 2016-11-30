package grpc

import (
	"fmt"

	"github.com/TheThingsNetwork/go-utils/log"
	"google.golang.org/grpc/grpclog"
)

// Wrap grpclog.Logger
func Wrap(logger log.Interface) grpclog.Logger {
	return &grpcWrapper{logger}
}

type grpcWrapper struct {
	log.Interface
}

func (w *grpcWrapper) Fatal(args ...interface{}) {
	w.Interface.Fatal(fmt.Sprint(args...))
}
func (w *grpcWrapper) Fatalln(args ...interface{}) {
	w.Fatal(args...)
}
func (w *grpcWrapper) Print(args ...interface{}) {
	w.Interface.Debug(fmt.Sprint(args...))
}
func (w *grpcWrapper) Printf(format string, args ...interface{}) {
	w.Interface.Debugf(format, args...)
}
func (w *grpcWrapper) Println(args ...interface{}) {
	w.Print(args...)
}
