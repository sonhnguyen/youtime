package youtuber_test

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"youtuber"

	"github.com/boltdb/bolt"
	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

var odooConfig map[string]string

// TestMain sets up the entire suite
func TestMain(m *testing.M) {
	// Global db setup
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("cannot retrieve present working directory: %s", err)
	}
	db, err := bolt.Open(path.Join(pwd, "test.db"), 0600, nil)
	if err != nil {
		log.Fatalf("unable to open bolt db: %s", err)
	}

	err = LoadConfiguration(pwd)
	if err != nil && viper.GetBool("isProduction") {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	odooConfig = viper.GetStringMapString("odoo")

	testDB = &youtuber.DB{DB: db}
	err = testDB.CreateAllBuckets()
	if err != nil {
		log.Fatalf("unable to create all bucket: %s", err)
	}

	// Check in the first time
	exists := testDB.DoesAnyUserExist()
	if exists {
		log.Fatal("some users currently exists before test begins")
	}

	user1, _ = testDB.CreateUser("testuser", "demotoken")
	retCode := m.Run()
	err = db.Close()
	if err != nil {
		log.Fatalf("error on closing db %s", err)
	}
	os.Exit(retCode)
}

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
