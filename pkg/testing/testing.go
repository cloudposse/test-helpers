package testing

import (
	"fmt"
	"os"
)

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Parallel()
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

type CustomT struct {
	TestingT
}

func (m *CustomT) Error(args ...interface{}) {
	fmt.Println(append([]interface{}{"Error:"}, args...)...)
}
func (m *CustomT) Errorf(format string, args ...interface{}) {
	fmt.Printf("Error: "+format+"\n", args...)
}
func (m *CustomT) Fail()                     {}
func (m *CustomT) FailNow()                  { os.Exit(1) }
func (m *CustomT) Failed() bool              { return false }
func (m *CustomT) Fatal(args ...interface{}) { fmt.Print("Fatal: "); fmt.Println(args...); os.Exit(1) }
func (m *CustomT) Fatalf(format string, args ...interface{}) {
	fmt.Printf("Fatal: "+format+"\n", args...)
	os.Exit(1)
}
func (m *CustomT) Helper()                                  {}
func (m *CustomT) Log(args ...interface{})                  { fmt.Println(append([]interface{}{"Log:"}, args...)...) }
func (m *CustomT) Logf(format string, args ...interface{})  { fmt.Printf("Log: "+format+"\n", args...) }
func (m *CustomT) Name() string                             { return "CustomT" }
func (m *CustomT) Skip(args ...interface{})                 {}
func (m *CustomT) SkipNow()                                 {}
func (m *CustomT) Skipf(format string, args ...interface{}) {}
func (m *CustomT) Skipped() bool                            { return false }
