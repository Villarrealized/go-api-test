package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const dataPath string = "data/"

func SaveData(data interface{}, filename string) error {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	dirExists, err := exists(dataPath)
	if err != nil {
		return err
	}

	if !dirExists {
		err = os.Mkdir(dataPath, 0755)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(dataPath+filename, jsonBytes, 0755)
	if err != nil {
		return err
	}

	return nil
}

func ReadData(filename string) ([]byte, error) {
	path := dataPath + filename

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	fmt.Println(err)
	return false, err
}
