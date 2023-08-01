package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func SHA256File(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建SHA256哈希对象
	sha256Hash := sha256.New()

	// 逐块读取文件并更新哈希值
	if _, err := io.Copy(sha256Hash, file); err != nil {
		return "", err
	}

	// 返回十六进制格式的哈希值
	return fmt.Sprintf("%x", sha256Hash.Sum(nil)), nil
}

func SHA256Str(str string) string {
	// 创建SHA256哈希对象
	sha256Hash := sha256.New()

	// 将字符串转换为字节数组并计算哈希值
	sha256Hash.Write([]byte(str))

	// 返回十六进制格式的哈希值
	return fmt.Sprintf("%x", sha256Hash.Sum(nil))
}
