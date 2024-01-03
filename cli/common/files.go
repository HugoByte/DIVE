package common

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type diveFileHandler struct{}

// The function returns a new instance of the diveFileHandler struct.
func NewDiveFileHandler() *diveFileHandler {
	return &diveFileHandler{}
}

// The `ReadFile` method is responsible for reading the contents of a file given its file path.
func (df *diveFileHandler) ReadFile(filePath string) ([]byte, error) {

	fileData, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		_, err := df.OpenFile(filePath, "append|write|create", 0644)
		if err != nil {
			return nil, WrapMessageToErrorf(ErrNotFound, "Error While Creating File %s", err.Error())
		}

		return []byte{}, nil

	} else if err != nil {
		return nil, WrapMessageToErrorf(ErrOpenFile, "Error While Reading File %s", err.Error())
	}

	return fileData, nil
}

// The `ReadJson` method is responsible for reading a JSON file and unmarshaling its contents into
// the provided object.
func (df *diveFileHandler) ReadJson(fileName string, obj interface{}) error {

	var filePath string

	if filepath.IsAbs(fileName) {
		filePath = fileName
	} else {
		pwd, err := df.GetPwd()
		if err != nil {
			return WrapMessageToErrorf(ErrPath, "Failed to get present working dir %s", err.Error())
		}
		outputDirPath := filepath.Join(pwd, DiveOutFileDirectory, EnclaveName)

		err = df.MkdirAll(outputDirPath, 0755)
		if err != nil {
			return WrapMessageToError(err, "Failed to Create Output Directory")
		}
		filePath = filepath.Join(outputDirPath, fileName)

	}

	data, err := df.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(data) != 0 {
		if err := json.Unmarshal(data, obj); err != nil {
			return WrapMessageToErrorf(ErrDataUnMarshall, " %s object %v", err.Error(), obj)
		}
	}

	return nil
}

// The `ReadAppFile` method is responsible for reading the contents of a file located in the
// application directory.
func (df *diveFileHandler) ReadAppFile(fileName string) ([]byte, error) {

	appFilePath, err := df.GetAppDirPathOrAppFilePath(fileName)
	if err != nil {
		return nil, WrapMessageToErrorf(ErrPath, "%s. path:%s", err, fileName)
	}

	data, err := df.ReadFile(appFilePath)

	if err != nil {
		return nil, WrapMessageToErrorf(err, "Invalid file path %s", appFilePath)
	}

	return data, nil
}

// The `WriteAppFile` method is responsible for writing data to a file located in the application
// directory.
func (df *diveFileHandler) WriteAppFile(fileName string, data []byte) error {

	appFileDir, err := df.GetAppDirPathOrAppFilePath("")
	if err != nil {
		return WrapMessageToErrorf(ErrPath, "%s. path:%s", err, fileName)
	}

	err = df.MkdirAll(appFileDir, os.ModePerm)

	if err != nil {
		return WrapMessageToErrorf(ErrWriteFile, "%s. path:%s", err, appFileDir)
	}

	appFilePath, err := df.GetAppDirPathOrAppFilePath(fileName)
	if err != nil {
		return WrapMessageToErrorf(err, "Invalid file path %s", appFilePath)
	}

	file, err := df.OpenFile(appFilePath, "append|write|create|truncate", 0644)
	if err != nil {
		return WrapMessageToErrorf(ErrOpenFile, "%s . Failed To Open App File %s for write", err, fileName)
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return WrapMessageToErrorf(ErrWriteFile, "%s . Failed To Write to App File %s", err, fileName)
	}

	return nil
}

// The `WriteFile` method is responsible for writing data to a file. It takes the file name and the
// data to be written as parameters.
func (df *diveFileHandler) WriteFile(fileName string, data []byte) error {

	pwd, err := df.GetPwd()

	if err != nil {
		return WrapMessageToErrorf(ErrWriteFile, "%s .Failed to Write File %s", err, fileName)
	}
	outputDirPath := filepath.Join(pwd, DiveOutFileDirectory, EnclaveName)

	filePath := filepath.Join(outputDirPath, fileName)

	err = df.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return WrapMessageToError(err, "Failed to Create Output Directory")
	}

	file, err := df.OpenFile(filePath, "write|append|create|truncate", 0644)

	if err != nil {
		return WrapMessageToError(err, "Failed")
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}

	return nil
}

// The `WriteJson` method is responsible for serializing the provided data object into JSON format and
// writing it to a file.
func (df *diveFileHandler) WriteJson(fileName string, data interface{}) error {

	serializedData, err := json.Marshal(data)

	if err != nil {
		return ErrDataMarshall
	}

	err = df.WriteFile(fileName, serializedData)
	if err != nil {
		return ErrWriteFile
	}
	return nil
}

// The `GetPwd()` function is a method of the `diveFileHandler` struct. It is responsible for
// retrieving the present working directory (PWD) and returning it as a string. It uses the
// `os.Getwd()` function to get the PWD and returns it along with any error that occurred during the
// process.
func (df *diveFileHandler) GetPwd() (string, error) {

	pwd, err := os.Getwd()
	if err != nil {
		return "", ErrPath
	}
	return pwd, err
}

// The `MkdirAll` function is a method of the `diveFileHandler` struct. It is responsible for creating
// a directory at the specified `dirPath` if it does not already exist.
func (df *diveFileHandler) MkdirAll(dirPath string, permission fs.FileMode) error {

	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, permission); err != nil {
			return WrapMessageToError(ErrWriteFile, err.Error())
		}
	} else if err != nil {

		return WrapMessageToError(ErrPath, "Failed to check directory existence")
	}

	return nil
}

// The `OpenFile` method is responsible for opening a file given its file path, file open mode, and
// permission. It uses the `os.OpenFile` function to open the file with the specified mode and
// permission. If there is an error during the file opening process, it returns an error with a wrapped
// message. Otherwise, it returns the opened file.
func (df *diveFileHandler) OpenFile(filePath string, fileOpenMode string, permission int) (*os.File, error) {
	mode := parseFileOpenMode(fileOpenMode)
	file, err := os.OpenFile(filePath, mode, fs.FileMode(permission))
	if err != nil {
		return nil, WrapMessageToErrorf(ErrOpenFile, "%s. Failed to Open File %s", err, filePath)
	}

	return file, nil

}

// The `GetHomeDir()` function is a method of the `diveFileHandler` struct. It is responsible for
// retrieving the user's home directory and returning it as a string. It uses the `os.UserHomeDir()`
// function to get the home directory and returns it along with any error that occurred during the
// process.
func (df *diveFileHandler) GetHomeDir() (string, error) {

	uhd, err := os.UserHomeDir()
	if err != nil {
		return "", WrapMessageToError(ErrPath, err.Error())
	}
	return uhd, err
}

// The function `parseFileOpenMode` takes a string representing file open modes separated by "|" and
// returns the corresponding integer value.
func parseFileOpenMode(fileOpenMode string) int {
	modes := strings.Split(fileOpenMode, "|")

	var mode int
	for _, m := range modes {
		switch strings.TrimSpace(m) {
		case "append":
			mode |= os.O_APPEND
		case "create":
			mode |= os.O_CREATE
		case "truncate":
			mode |= os.O_TRUNC
		case "write":
			mode |= os.O_WRONLY
		case "readwrite":
			mode |= os.O_RDWR
		case "read":
			mode |= os.O_RDONLY
		}

	}

	return mode
}

// The `RemoveFile` function is a method of the `diveFileHandler` struct. It is responsible for
// removing a file from the file system.
func (df *diveFileHandler) RemoveFile(fileName string) error {

	pwd, err := df.GetPwd()

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Remove File")
	}

	filePath := filepath.Join(pwd, fileName)

	_, err = os.Stat(filePath)
	if err != nil {
		return WrapMessageToErrorf(ErrNotExistsFile, "%s. PATH:%s", err, filePath)
	}

	err = os.Remove(filePath)
	if err != nil {
		return WrapMessageToErrorf(ErrPath, "%s.Failed To Remove File %s", err, filePath)
	}
	return nil
}

// The `RemoveFiles` function is a method of the `diveFileHandler` struct. It is responsible for
// removing multiple files from the file system.
func (df *diveFileHandler) RemoveFiles(fileNames []string) error {

	pwd, err := df.GetPwd()

	if err != nil {
		return WrapMessageToErrorf(ErrPath, "Failed To Remove File")
	}
	for _, fileName := range fileNames {
		filePath := filepath.Join(pwd, fileName)

		_, err = os.Stat(filePath)
		if err == nil {
			err = os.Remove(filePath)
			if err != nil {
				return WrapMessageToErrorf(ErrInvalidFile, "%s Failed To Remove File %s", err, filePath)
			}
		}

	}
	return nil
}

// The `RemoveDir` function is a method of the `diveFileHandler` struct. It is responsible for
// removing output directories from the file system.
func (df *diveFileHandler) RemoveDir(enclaveName string) error {

	pwd, err := df.GetPwd()

	if err != nil {
		return WrapMessageToErrorf(ErrPath, "Failed To Remove Directory")
	}
	dirPath := filepath.Join(pwd, DiveOutFileDirectory, enclaveName)

	_, err = os.Stat(dirPath)
	if err == nil {
		err = os.RemoveAll(dirPath)
		if err != nil {
			return WrapMessageToErrorf(ErrInvalidFile, "%s Failed To Remove Directory %s", err, enclaveName)
		}
	}

	return nil
}

// The `RemoveAllDir` function is a method of the `diveFileHandler` struct. It is responsible for
// removing all output directories from the file system.
func (df *diveFileHandler) RemoveAllDir() error {

	pwd, err := df.GetPwd()

	if err != nil {
		return WrapMessageToErrorf(ErrPath, "Failed To Remove Directory")
	}
	dirPath := filepath.Join(pwd, DiveOutFileDirectory)

	_, err = os.Stat(dirPath)
	if err == nil {
		err = os.RemoveAll(dirPath)
		if err != nil {
			return WrapMessageToErrorf(ErrInvalidFile, "%s Failed To Remove Output Directory", err)
		}
	}

	return nil
}

// The `GetAppDirPathOrAppFilePath` function is a method of the `diveFileHandler` struct. It is
// responsible for returning the file path of a file located in the application directory.
func (df *diveFileHandler) GetAppDirPathOrAppFilePath(fileName string) (string, error) {

	var path string
	uhd, err := df.GetHomeDir()
	if err != nil {
		return "", WrapMessageToErrorf(err, "Failed To Write App File %s", fileName)
	}
	if fileName == "" {
		path = filepath.Join(uhd, DiveAppDir)
	} else {
		path = filepath.Join(uhd, DiveAppDir, fileName)
	}

	return path, nil
}
