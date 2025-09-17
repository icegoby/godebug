package godebug

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/k0kubun/pp"
)

type Config struct {
	Debug bool
	File  *os.File
}

const MAC_STR = "%02x:%02x:%02x:%02x:%02x:%02x"

func MAC2STR(b []byte) string {
	if len(b) != 6 {
		return "NOT MAC ADDRESS"
	}
	return fmt.Sprintf(MAC_STR, b[0], b[1], b[2], b[3], b[4], b[5])
}

var conf Config

func SetDebug(debug bool) {
	conf.Debug = debug
}

func SetUs(us bool) {
	f := log.Flags()
	if us {
		f |= log.Lmicroseconds
	} else {
		f &= ^log.Lmicroseconds
	}
	log.SetFlags(f)
}

func GetDebug() bool {
	return conf.Debug
}

func SetFile(fname string) error {
	if f, err := os.Create(fname); err != nil {
		return errors.New(fmt.Sprintf("failed to open file '%v'", fname))
	} else {
		conf.File = f
		log.SetOutput(conf.File)
		return nil
	}
}

func SetStdout() {
	if conf.File != nil {
		conf.File.Close()
		conf.File = nil
		log.SetOutput(os.Stdout)
	}
}

func SetStderr() {
	if conf.File != nil {
		conf.File.Close()
		conf.File = nil
		log.SetOutput(os.Stderr)
	}
}

func Deinit() {
	if conf.File != nil {
		conf.File.Close()
		conf.File = nil
		log.SetOutput(os.Stderr)
	}
}

func ___func(level int) string {
	pc, _, _, ok := runtime.Caller(level)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)
	tmp := strings.Split(f.Name(), ".")
	name := tmp[len(tmp)-1]
	tmp = strings.Split(file, "/")
	return fmt.Sprintf("%v:%v[%v]", tmp[len(tmp)-1], name, line)
}

func __func() string {
	return ___func(1)
}

func DPrintf(f string, arg ...interface{}) {
	if conf.Debug {
		s := fmt.Sprintf("%v: ", ___func(2))
		s += fmt.Sprintf(f, arg...)
		log.Printf(s)
	}
}

func DPPrintf(arg interface{}) {
	if conf.Debug {
		s := fmt.Sprintf("%v: ", ___func(2))
		s += pp.Sprintf("%v", arg)
		log.Printf(s)
	}
}

func LPrintf(f string, arg ...interface{}) {
	s := fmt.Sprintf("%v: ", ___func(2))
	s += fmt.Sprintf(f, arg...)
	log.Printf(s)
}

func SPrintf(f string, arg ...interface{}) string {
	s := fmt.Sprintf("%v: ", ___func(2))
	s += fmt.Sprintf(f, arg...)
	return s
}

func EPrintf(f string, arg ...interface{}) error {
	s := fmt.Sprintf("%v: ", ___func(2))
	s += fmt.Sprintf(f, arg...)
	return errors.New(s)
}

func SHexDump(d []byte) string {
	var i int
	s := ""
	for i = 0; i < len(d); i++ {
		if i%16 == 0 {
			if i != 0 {
				s += fmt.Sprintf("\n")
			}
			s += fmt.Sprintf("%08x:", i)
		} else if i%8 == 0 {
			s += fmt.Sprintf(" :")
		}
		s += fmt.Sprintf(" %02x", d[i])
	}
	if i%16 != 0 {
		s += fmt.Sprintf("\n")
	}
	return s
}

func _LHexDump(label string, d []byte) {
	var i int
	log.Printf("%v: %v", ___func(3), label)
	s := ""
	for i = 0; i < len(d); i++ {
		if i%16 == 0 {
			if i != 0 {
				log.Printf(s)
			}
			s = fmt.Sprintf("\t%08x:", i)
		} else if i%8 == 0 {
			s += fmt.Sprintf(" :")
		}
		s += fmt.Sprintf(" %02x", d[i])
	}
	log.Printf(s)
}

func LHexDump(label string, d []byte) {
	_LHexDump(label, d)
}

func DHexDump(label string, d []byte) {
	if !conf.Debug {
		return
	}
	_LHexDump(label, d)
}
