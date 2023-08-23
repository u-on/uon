package uon

import (
	"fmt"
	"io"
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
	s, err := os.Stat(path)
	if err != nil {

		return false
	}
	return s.IsDir()

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
		fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFile(fullname)
			if err != nil {
				//fmt.Printf("读取文件目录失败,fullname=%v, err=%v", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}

// ReadFileToStr
// @description: 返回文件内容str
// @param {string} filepath - 文件路径
// @return {string} - 读取的字符串
func ReadFileToStr(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {

		return "", err
	}
	defer f.Close()

	fd, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fd), nil
}
