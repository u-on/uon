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
	Releases gjson.Result
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

	rt.Releases = result
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

	for _, name := range t.Releases.Array() {
		result, err := regexp.MatchString(Ex, name.String())
		if err != nil {
			fmt.Println(err.Error())
		}
		if result {
			return name.String(), nil
			//fmt.Printf("%s matches\n", tt)
		}
	}
	return "", errors.New("null")
}
