package godebug

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/k0kubun/pp"
)

type GoDebug struct {
	Debug  bool
	File   *os.File
	isFile bool
}

func NewGoDebug() *GoDebug {
	return &GoDebug{
		Debug:  false,
		File:   os.Stdout,
		isFile: false,
	}
}

const MAC_STR = "%02x:%02x:%02x:%02x:%02x:%02x"

func MAC2STR(b []byte) string {
	if len(b) != 6 {
		return "NOT MAC ADDRESS"
	}
	return fmt.Sprintf(MAC_STR, b[0], b[1], b[2], b[3], b[4], b[5])
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

func SetUs(us bool) {
	f := log.Flags()
	if us {
		f |= log.Lmicroseconds
	} else {
		f &= ^log.Lmicroseconds
	}
	log.SetFlags(f)
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

var global_gdb *GoDebug

func Init() {
	if global_gdb == nil {
		global_gdb = NewGoDebug()
	}
}

func SetDebug(debug bool) {
	Init()
	global_gdb.SetDebug(debug)
}

func GetDebug() bool {
	Init()
	return global_gdb.GetDebug()
}

func SetFile(fname string) error {
	Init()
	return global_gdb.SetFile(fname)
}

func SetStdout() {
	Init()
	global_gdb.SetStdout()
}

func SetStderr() {
	Init()
	global_gdb.SetStderr()
}

func Deinit() {
	SetStdout()
}

func DPrintf(f string, arg ...interface{}) {
	Init()
	global_gdb.DPrintf(f, arg...)
}

func DPPrintf(arg interface{}) {
	Init()
	global_gdb.DPPrintf(arg)
}

func LPrintf(f string, arg ...interface{}) {
	Init()
	global_gdb.LPrintf(f, arg...)
}

func _LHexDump(label string, d []byte) {
	Init()
	global_gdb._LHexDump(label, d)
}

func LHexDump(label string, d []byte) {
	_LHexDump(label, d)
}

func DHexDump(label string, d []byte) {
	Init()
	if !global_gdb.Debug {
		return
	}
	_LHexDump(label, d)
}

//
//		methods of GoDebug
//

func (gdb *GoDebug) Print(s string) {
	ss := strings.Split(s, "\n")
	for i, v := range ss {
		if len(v) == 0 {
			continue
		}
		vv := "                           " + v + "\n"
		if i == 0 {
			t := time.Now()
			vv = fmt.Sprint("%04d-%02d-%02d %02d:%02d:%02d.%06d %v\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000, v)
		}
		b := []byte(vv)
		if n, err := gdb.File.Write(b); err != nil {
			LPrintf("failed to write log (%v) (%v)", err.Error(), gdb)
		} else if n != len(b) {
			LPrintf("%v / %v bytes written (%v)", n, len(b), gdb)
		}
	}
}

func (gdb *GoDebug) SetDebug(debug bool) {
	gdb.Debug = debug
}

func (gdb *GoDebug) GetDebug() bool {
	return gdb.Debug
}

func (gdb *GoDebug) SetFile(fname string) error {
	gdb.SetStdout()
	if f, err := os.Create(fname); err != nil {
		return errors.New(fmt.Sprintf("failed to open file '%v'", fname))
	} else {
		gdb.File = f
		gdb.isFile = true
		return nil
	}
}

func (gdb *GoDebug) SetStdout() {
	if gdb.isFile {
		gdb.File.Close()
	}
	gdb.File = os.Stdout
	gdb.isFile = false
}

func (gdb *GoDebug) SetStderr() {
	if gdb.isFile {
		gdb.File.Close()
	}
	gdb.File = os.Stderr
	gdb.isFile = false
}

func (gdb *GoDebug) Deinit() {
	gdb.SetStdout()
}

func (gdb *GoDebug) DPrintf(f string, arg ...interface{}) {
	if gdb.Debug {
		s := fmt.Sprintf("%v: ", ___func(2))
		s += fmt.Sprintf(f, arg...)
		gdb.Print(s)
	}
}

func (gdb *GoDebug) DPPrintf(arg interface{}) {
	if gdb.Debug {
		s := fmt.Sprintf("%v: ", ___func(2))
		s += pp.Sprintf("%v", arg)
		gdb.Print(s)
	}
}

func (gdb *GoDebug) LPrintf(f string, arg ...interface{}) {
	s := fmt.Sprintf("%v: ", ___func(2))
	s += fmt.Sprintf(f, arg...)
	gdb.Print(s)
}

func (gdb *GoDebug) _LHexDump(label string, d []byte) {
	var i int

	s := fmt.Sprintf("%v: %v", ___func(3), label)
	gdb.Print(s)
	s = ""
	for i = 0; i < len(d); i++ {
		if i%16 == 0 {
			if i != 0 {
				gdb.Print(s)
			}
			s = fmt.Sprintf("\t%08x:", i)
		} else if i%8 == 0 {
			s += fmt.Sprintf(" :")
		}
		s += fmt.Sprintf(" %02x", d[i])
	}
	gdb.Print(s)
}

func (gdb *GoDebug) LHexDump(label string, d []byte) {
	gdb._LHexDump(label, d)
}

func (gdb *GoDebug) DHexDump(label string, d []byte) {
	if !gdb.Debug {
		return
	}
	gdb._LHexDump(label, d)
}
