package github

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"regexp"
)

type GithubInfo struct {
	Releases []string //ads        gjson.Result
	Tag      string
}

func GetReleasesInfo(rep string) (GithubInfo, error) {
	var rt GithubInfo
	gurl := "https://api.github.com/repos/" + rep + "/releases/latest"
	req, _ := http.NewRequest("GET", gurl, nil)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("[ERROR]: " + err.Error())

		return GithubInfo{}, err // Result.Type == gjson.Null
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// 获取所有下载地址 https://api.github.com/repos/?/?/releases/latest
	result := gjson.Get(string(body), "assets.#.browser_download_url")

	array := result.Array()

	// 创建一个字符串切片
	strSlice := make([]string, len(array))

	// 遍历JSON数组并将元素转换为字符串
	for i, element := range array {
		strSlice[i] = element.String()
	}

	rt.Releases = strSlice
	ghTag := gjson.Get(string(body), "tag_name").String()
	rt.Tag = ghTag

	return rt, nil

}

// GetReleasesEx 获取下载地址  正则匹配
func GetReleasesEx(rep string, Ex string) (string, error) {

	t, err := GetReleasesInfo(rep)
	if err != nil {
		return "", errors.New("error")
	}

	// 使用for循环遍历字符串数组
	for i := 0; i < len(t.Releases); i++ {
		result, err := regexp.MatchString(Ex, t.Releases[i])
		if err != nil {
			fmt.Println(err.Error())
		}
		if result {
			return t.Releases[i], nil
			//fmt.Printf("%s matches\n", tt)
		}

	}

	return "", errors.New("null")
}
