package writers

import (
	"bufio"
	"bytes"
	"os"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"fmt"

	"encoding/binary"
)

// FileConfig ...
type FileConfig struct {
	directory   string
	fileName    string
	fileMaxSize int64
	flushTime   time.Duration
}

// FileWriter ...
type FileWriter struct {
	writer        *bufio.Writer
	config        *FileConfig
	queue         IList
	formatHandler FormatHandler
	mux           *sync.Mutex
	outOnEmpty    bool
	quit          chan bool
}

// NewFileWriter ...
func NewFileWriter(options ...FileWriterOption) *FileWriter {
	fileWriter := &FileWriter{
		queue:         NewQueue(WithMode(FIFO)),
		formatHandler: JsonFormatHandler,
		mux:           &sync.Mutex{},
		config:        &FileConfig{},
		quit:          make(chan bool),
	}
	fileWriter.Reconfigure(options...)
	fileWriter.start()

	return fileWriter
}

func (fileWriter *FileWriter) start() error {
	if _, err := os.Stat(fileWriter.config.directory); os.IsNotExist(err) {
		if err = os.Mkdir(fileWriter.config.directory, 0777); err != nil {
			return err
		}
	}

	go func(fileWriter *FileWriter) {
		var tmpLogFileName string
		var logMessage []byte
		for {
			select {
			case <-fileWriter.quit:
				if fileWriter.queue.IsEmpty() {
					return
				} else {
					fileWriter.outOnEmpty = true
				}

			case <-time.After(fileWriter.config.flushTime):
				fileWriter.mux.Lock()
				defer fileWriter.mux.Unlock()

			newFile:
				tmpLogFileName = fmt.Sprintf("%s/%s%s", fileWriter.config.directory, fileWriter.config.fileName, time.Now().Format("2006.01.02 15.04.05.06"))
				file, err := os.OpenFile(tmpLogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
				checkError(err, fmt.Sprintf("error opening file %s: %s", tmpLogFileName, err), file)

				fileSize, _ := file.Stat()
				maxSize := fileWriter.config.fileMaxSize - fileSize.Size()
				buffer := bytes.NewBuffer(make([]byte, 0))

				for fileWriter.queue.Size() > 0 {
					value := fileWriter.queue.Remove()
					switch value.(type) {
					case []byte:
						logMessage = value.([]byte)
					case Message:
						message := value.(Message)
						if bytes, err := fileWriter.formatHandler(message.Prefixes, message.Tags, message.Message, message.Fields, message.Sufixes); err != nil {
							continue
						} else {
							logMessage = bytes
						}
					}

					if int64(binary.Size(buffer.Bytes())+binary.Size(logMessage)) <= maxSize {
						buffer.Write(logMessage)
					} else {
						if _, err := file.Write(buffer.Bytes()); err != nil {
							checkError(err, fmt.Sprintf("error writing file %s: %s", tmpLogFileName, err), file)
						}
						file.Close()
						goto newFile
					}
				}

				if _, err := file.Write(buffer.Bytes()); err != nil {
					checkError(err, fmt.Sprintf("error flushing to file %s: %s", tmpLogFileName, err), file)
				}
				file.Close()

				if fileWriter.queue.IsEmpty() && fileWriter.outOnEmpty {
					return
				}
			}
		}
	}(fileWriter)
	return nil
}

// Write ...
func (fileWriter *FileWriter) Write(message []byte) (n int, err error) {
	id := uuid.NewV4()
	fileWriter.queue.Add(id.String(), message)
	return 0, nil
}

// SWrite ...
func (fileWriter *FileWriter) SWrite(prefixes map[string]interface{}, tags map[string]interface{}, message interface{}, fields map[string]interface{}, sufixes map[string]interface{}) (n int, err error) {
	id := uuid.NewV4()
	fileWriter.queue.Add(id.String(), Message{Prefixes: prefixes, Tags: tags, Message: message, Fields: fields, Sufixes: sufixes})
	return 0, nil
}
