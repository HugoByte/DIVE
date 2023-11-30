package common

import "os"

type diveFileHandler struct{}

func NewDiveFileHandler() *diveFileHandler {
	return &diveFileHandler{}
}

func (f *diveFileHandler) ReadFile(filePath string) ([]byte, error) {
	return nil, nil
}
func (f *diveFileHandler) ReadJson(filePath string, obj interface{}) (string, error) {
	return "nil", nil
}
func (f *diveFileHandler) WriteFile(filePath string, data []byte) error {
	return nil
}
func (f *diveFileHandler) WriteJson(filePath string, data interface{}) error {
	return nil
}
func (f *diveFileHandler) GetPwd() string {
	return ""
}
func (f *diveFileHandler) MkdirAll(dirPath string, permission string) error {
	return nil
}
func (f *diveFileHandler) OpenFile(filePath string, fileOpenMode string, permission int) (*os.File, error) {
	return nil, nil
}
