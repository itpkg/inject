package ioc_test

import (
	"fmt"
	"log/syslog"
	"testing"
	"time"

	"github.com/itpkg/ioc"
)

func Hello(mod *Model, logger *syslog.Writer, version string) string {
	s := fmt.Sprintf("Model=>%v, Version=>%s", mod, version)
	logger.Info(s)
	return s
}

type Model struct {
	Now     *time.Time `inject:""`
	Version int        `inject:"version"`
}

func TestInjector(t *testing.T) {
	inj := ioc.New()

	wrt, _ := syslog.New(syslog.LOG_DEBUG, "itpkg")
	now := time.Now()
	inj.Provide(
		&ioc.Object{Value: &Model{}},
		&ioc.Object{Value: wrt},
		&ioc.Object{Value: &now},
		&ioc.Object{Name: "version", Value: 20150922},
		&ioc.Object{Value: 1.1},
		&ioc.Object{Name: "hello", Value: "Hello, it-package!"},
	)

	if err := inj.Populate(); err == nil {
		t.Logf(inj.String())
	} else {
		t.Errorf("error on populate: %v", err)
	}

	if vls, err := inj.Run(Hello, "v20150923"); err == nil {
		t.Logf(vls[0].(string))
	} else {
		t.Errorf("error on run: %v", err)
	}

}
