package gomanager

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(fileName string, obj interface{}) ([]byte, error) {
	var err error

	if !Exists(fileName) {
		fileName = global["path"].(string) + fileName
	}

	log.Infof("loading file [ %s ]", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if obj != nil {
		log.Infof("unmarshalling file [ %s ] to struct", fileName)
		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ReadFileLines(fileName string) ([]string, error) {
	lines := make([]string, 0)

	if !Exists(fileName) {
		fileName = global["path"].(string) + fileName
	}

	log.Infof("loading file [ %s ]", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(fileName string, obj interface{}) error {
	if !Exists(fileName) {
		fileName = global["path"].(string) + fileName
	}

	log.Infof("writing file [ %s ]", fileName)
	jsonBytes, _ := json.MarshalIndent(obj, "", "    ")
	if err := ioutil.WriteFile(fileName, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}
