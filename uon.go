package uon

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// PauseExit Press Ctrl+c to exit
func PauseExit() {
	fmt.Printf("Press Ctrl+c to exit...")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	for {
		if sig.String() == "interrupt" {
			break
		}
	}
	return

}

// Sleep 延迟 单位ms
func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// SelfDir 获取程序自身所在目录 不含\
func SelfDir() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return exPath
}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//nonexistent来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, _ := os.Stat(path)
	return s != nil && s.IsDir()

}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {

	return !IsDir(path)

}

// GetAllFile 递归获取指定目录下的所有文件名
func GetAllFile(pathname string) ([]string, error) {
	var result []string

	fis, err := os.ReadDir(pathname)
	if err != nil {
		//fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		if fi.IsDir() {
			// 使用递归获取子目录下的文件
			subResult, err := GetAllFile(fullname)
			if err != nil {
				return result, err
			}
			result = append(result, subResult...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}

// ReadFileToStr
// @description: 读取文件内容str
// @param {string} filepath - 文件路径
// @return {string} - 读取的字符串
func ReadFileToStr(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteStringToFile
// @description: 写入
// @param {string} str - 写入的字符串
// @param {string} fileName - 文件名
// @return {error} - 返回
func WriteStringToFile(str string, fileName string) error {
	err := os.WriteFile(fileName, []byte(str), 0666)
	return err
}
