package writers

import "os"

func checkError(err error, message string, file *os.File) {
	if err != nil {
		if file != nil {
			file.Close()
		}
		panic(message)
	}
}
