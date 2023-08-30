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
	"strconv"
	"sync"
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

func DownloadX(url string, savepath string, numWorkers int) error {

	_path, _name := filepath.Split(savepath)

	if dd, _ := uon.PathExists(_path); dd == false {
		err := os.MkdirAll(_path, 0777)
		if err != nil {
			return err
		}
	}
	if _name == "" {
		_, _name = filepath.Split(url)
		savepath = savepath + _name
	}

	resp, err := http.Head(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// 创建文件
	file, err := os.Create(savepath)
	if err != nil {

		return err
	}
	defer file.Close()

	bar := progressbar.NewOptions(size,
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

	// 开始下载
	start := 0
	partSize := size / numWorkers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		end := start + partSize
		if i == numWorkers-1 {
			end = size
		}

		// 开启一个协程下载一部分文件
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			_download(start, end, file, bar, url)
		}(start, end)

		start = end + 1
	}

	// 等待所有协程结束
	wg.Wait()

	return nil

}

// 下载指定区间的文件数据，并写入文件中
func _download(start, end int, file *os.File, bar *progressbar.ProgressBar, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		file.WriteAt(buf[:n], int64(start))
		start += n
		bar.Add(n)
	}
}
