package libconfig

/*
#cgo pkg-config: libconfig
#include <libconfig.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	CONFIG_TRUE  C.int = 1
	CONFIG_FALSE C.int = 0
)

type Lookupable interface {
	Lookup(string) Setting
	LookupInt(string) int64
	LookupFloat(string) float64
	LookupBool(string) bool
	LookupString(string) string
}

type Config struct {
	cconf C.struct_config_t
}

type Setting struct {
	csetting *C.struct_config_setting_t
}

func NewConfig() Config {
	var conf Config
	C.config_init(&conf.cconf)
	return conf
}

func (config *Config) Destroy() {
	C.config_destroy(&config.cconf)
}

func (config *Config) Clear() {
	C.config_clear(&config.cconf)
}

// cgo will copy the whole content - which i do very much not like,
func (config *Config) ReadString(str string) error {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	rc := C.config_read_string(&config.cconf, cstr)
	if rc == CONFIG_FALSE {
		return config.error("load")
	}
	return nil
}

func (config *Config) ReadFile(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	rc := C.config_read_file(&config.cconf, cfilename)
	if rc == CONFIG_FALSE {
		return config.error("load")
	}
	return nil
}

func (config *Config) WriteFile(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	rc := C.config_write_file(&config.cconf, cfilename)
	if rc == CONFIG_FALSE {
		return config.error("load")
	}
	return nil
}

// TODO All the options stuff

func (config *Config) SetIncludeDir(dir string) {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	C.config_set_include_dir(&config.cconf, cdir)
}

func (config *Config) GetIncludeDir(dir string) string {
	// include dir is in macros...
	//cdir := C.config_get_include_dir(unsafe.Pointer(config.cconf))
	//return C.GoString(cdir)
	return C.GoString(config.cconf.include_dir)
}

func (config *Config) Lookup(path string) (setting *Setting) {
	setting = new(Setting)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	setting.csetting = C.config_lookup(&config.cconf, cpath)
	return
}

func (config *Config) LookupInt(path string) (int, error) {
	var result C.longlong
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_lookup_int64(
		&config.cconf,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return 0, fmt.Errorf("'%s' Not Found", path)
	}
	return int(result), nil
}

func (config *Config) LookupFloat(path string) (float64, error) {
	var result C.double
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_lookup_float(&config.cconf,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return 0, fmt.Errorf("'%s' Not Found", path)
	}
	return float64(result), nil
}

func (config *Config) LookupBool(path string) (bool, error) {
	var result C.int
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_lookup_bool(&config.cconf,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return false, fmt.Errorf("'%s' Not Found", path)
	}
	return result != 0, nil
}

func (config *Config) LookupString(path string) (string, error) {
	// todo: maybe this segfaults?
	var result *C.char
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_lookup_string(&config.cconf,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return "", fmt.Errorf("'%s' Not Found", path)
	}
	return C.GoString(result), nil
}

//  SETTING stuff

func (setting *Setting) Lookup(path string) (subsetting *Setting) {
	subsetting = new(Setting)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	subsetting.csetting = C.config_setting_lookup(setting.csetting, cpath)
	return
}

func (setting *Setting) LookupInt(path string) (int, error) {
	var result C.longlong
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_setting_lookup_int64(
		setting.csetting,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return 0, fmt.Errorf("'%s' Not Found", path)
	}
	return int(result), nil
}

func (setting *Setting) LookupFloat(path string) (float64, error) {
	var result C.double
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_setting_lookup_float(
		setting.csetting,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return 0, fmt.Errorf("'%s' Not Found", path)
	}
	return float64(result), nil
}

func (setting *Setting) LookupBool(path string) (bool, error) {
	var result C.int
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_setting_lookup_bool(
		setting.csetting,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return false, fmt.Errorf("'%s' Not Found", path)
	}
	return result != 0, nil
}

func (setting *Setting) LookupString(path string) (string, error) {
	// todo: maybe this segfaults?
	var result *C.char
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	rc := C.config_setting_lookup_string(
		setting.csetting,
		cpath,
		&result)
	if rc == CONFIG_FALSE {
		return "", fmt.Errorf("'%s' Not Found", path)
	}
	return C.GoString(result), nil
}

func (config *Config) error(op string) error {
	// This ***** library implements those as macros, so we have to hack....
	//error_text := C.conf.error_text(unsafe.Pointer(config.cconf))
	//error_file := C.conf.error_file(unsafe.Pointer(config.cconf))
	//error_line := C.conf.error_line(unsafe.Pointer(config.cconf))

	error_text := C.GoString(config.cconf.error_text)
	error_file := C.GoString(config.cconf.error_file)
	error_line := int(config.cconf.error_line)

	return fmt.Errorf("Conf.Error: Operation: %s %s in %s:%d",
		op, error_text, error_file, error_line)
}
