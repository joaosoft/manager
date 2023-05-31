package web

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"io"

	"fmt"
	"path"

	"github.com/joaosoft/errors"
)

func GetEnv() string {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	return env
}

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(file string, obj interface{}) ([]byte, error) {
	var err error

	if !Exists(file) {
		return nil, errors.New(errors.LevelError, 0, "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if obj != nil {
		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ReadFileLines(file string) ([]string, error) {
	lines := make([]string, 0)

	if !Exists(file) {
		return nil, errors.New(errors.LevelError, 0, "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(file string, obj interface{}) error {
	if !Exists(file) {
		return errors.New(errors.LevelError, 0, "file don't exist")
	}

	jsonBytes, _ := json.MarshalIndent(obj, "", "    ")
	if err := ioutil.WriteFile(file, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetMimeType(fileName string) (mimeType string) {
	mimeType = mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType
}

func RandomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func NewSimpleConfig(file string, obj interface{}) error {
	dir, _ := os.Getwd()
	if _, err := ReadFile(fmt.Sprintf("%s%s", dir, file), obj); err != nil {
		return err
	}
	return nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}

func readData(obj reflect.Value, data map[string]string) error {
	types := reflect.TypeOf(obj.Interface())

	if !obj.CanInterface() {
		return nil
	}

checkAgain:
	if obj.Kind() == reflect.Ptr && !obj.IsNil() {
		obj = obj.Elem()

		if obj.IsValid() {
			types = obj.Type()
		} else {
			return nil
		}

		goto checkAgain
	}

	switch obj.Kind() {
	case reflect.Struct:
		for i := 0; i < types.NumField(); i++ {
			nextValue := obj.Field(i)
			nextType := types.Field(i)

			if nextValue.Kind() == reflect.Ptr {
				if !nextValue.IsNil() {
					nextValue = nextValue.Elem()
				} else {
					isSlice := nextValue.Kind() == reflect.Slice
					isMap := nextValue.Kind() == reflect.Map
					isMapOfSlices := isMap && nextValue.Type().Elem().Kind() == reflect.Slice

					if isMapOfSlices {
						nextValue = reflectAlloc(nextValue.Type().Elem().Elem())
					} else if isSlice || isMap {
						nextValue = reflectAlloc(nextValue.Type().Elem())
					} else {
						nextValue = reflectAlloc(nextValue.Type())
					}
				}
			}

			if !nextValue.CanInterface() {
				continue
			}

			var tagName string
			jsonName, exists := nextType.Tag.Lookup("json")
			if exists {
				tagName = strings.SplitN(jsonName, ",", 2)[0]
			}

			if value, ok := data[tagName]; ok {
				if nextValue.Kind() == reflect.Ptr {
					obj.Field(i).Set(nextValue)
					nextValue = nextValue.Elem()
				}

				if err := setValue(nextValue.Kind(), nextValue, value); err != nil {
					return err
				}
			}

			if err := readData(nextValue, data); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < obj.Len(); i++ {
			nextValue := obj.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := readData(nextValue, data); err != nil {
				return err
			}
		}
	case reflect.Map:

	default:
		// do nothing ...
	}
	return nil
}
