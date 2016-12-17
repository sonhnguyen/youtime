package youtuber_test

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

var odooConfig map[string]string

// TestMain sets up the entire suite
// TEST HELPERS

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func LoadConfiguration(pwd string) error {
	viper.SetConfigName("youtuber-config")
	viper.AddConfigPath(pwd)
	_, devPath, _, _ := runtime.Caller(1)
	devPath = path.Dir(devPath) + "/cmd/youtuberweb/"
	viper.AddConfigPath(devPath)
	viper.SetDefault("path", devPath)
	return viper.ReadInConfig() // Find and read the config file
}
