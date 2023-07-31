package download

import (
	"errors"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/u-on/uon"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
)

// Download 下载文件
func Download(url string, savepath string) error {
	_path, _name := filepath.Split(savepath)

	if dd, _ := uon.PathExists(_path); dd == false {
		err := os.MkdirAll(_path, 0777)
		if err != nil {
			fmt.Println(err)
			return errors.New("failed to create directory")
		}
	}
	if _name == "" {
		_, _name = filepath.Split(url)
		savepath = savepath + _name
	}

	// progressbar
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UserAgent)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(savepath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	//bar := progressbar.DefaultBytes(
	//	resp.ContentLength,
	//	"downloading",
	//)

	bar := progressbar.NewOptions(int(resp.ContentLength),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		//progressbar.OptionSetWidth(55),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[ "+_name+" ]:"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[yellow]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("")
	return nil
}
