package runcmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/u-on/uon/conver"
	"io"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// RunCommand 可关闭+实时输出
func RunCommand(cmd string, args ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancelFunc context.CancelFunc) {
		time.Sleep(3 * time.Second)
		cancelFunc()
	}(cancel)
	c := exec.CommandContext(ctx, cmd, args...) // mac linux
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != nil {
					fmt.Printf("error: %q", ctx.Err()) //程序出现错误
				} else {
					fmt.Println("Program terminated") //程序被终止
				}
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					return
				}
				if runtime.GOOS == "windows" { //处理windows中文乱码
					readString = conver.Utf8ToGbk(readString)
				}
				fmt.Print(readString)
			}
		}
	}(&wg)
	err = c.Start()
	wg.Wait()
	return err
}

// RunCommand2 实时显示
func RunCommand2(cmd string, str ...string) error {
	//c := exec.Command("cmd", "/C", cmd)   // windows
	c := exec.Command(cmd, str...) // mac or linux
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			if runtime.GOOS == "windows" { //处理windows中文乱码
				readString = conver.Utf8ToGbk(readString)
			}
			fmt.Print(readString)
		}
	}()
	err = c.Start()
	wg.Wait()
	return err
}

func RunCommand3(name string, arg ...string) error {

	cmd := exec.Command(name, arg...)
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {

		return err
	}

	if err = cmd.Start(); err != nil {

		return err
	}
	// 从管道中实时获取输出并打印到终端
	for {

		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		outstr := ""
		if runtime.GOOS == "windows" { //处理windows中文乱码
			outstr = conver.Utf8ToGbk(string(tmp))
		} else {
			outstr = string(tmp)
		}

		fmt.Print(outstr)
		if err != nil {

			break
		}
	}

	if err = cmd.Wait(); err != nil {

		return err
	}
	return nil
}

// Run 运行并等待
func Run(cmd string, args ...string) {
	exec.Command(cmd, args...).Run()
}

// start 运行不等待
func start(cmd string, args ...string) {
	exec.Command(cmd, args...).Start()
}
