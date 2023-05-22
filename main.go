package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// 获取参数
	args := os.Args

	args1 := args[1]   // 文件夹  "/root"
	args2 := args[2]   // 文件后缀 ".mp4"
	args3 := args[3]   // 推流地址 "srt://"
	// 获取到文件夹内所有文件
	paths, _ := GetDirAllFilePathsFollowSymlink(args1, false)
	// 筛选出特定文件类型
	var listfile []string //获取文件列表
	for _, path := range paths {
		ok := strings.HasSuffix(path, args2)
		if ok {
			listfile = append(listfile, path)
		}
	}
	fmt.Println("发现文件：", listfile)
	// 推流地址
	srtPath := args3
	for {
		// 设置 ffmpeg 命令行参数 循环推流
		for _, file := range listfile {
			fmt.Println("开始推流:", file)
			cmdArguments := []string{"-re", "-i", file, "-c", "copy", "-f", "mpegts", srtPath}
			CmdRun("ffmpeg", cmdArguments)
			fmt.Println("推流完成:", file)
		}
	}
}

// 常规运行，统一输出
func CmdRun(command string, cmdArguments []string) {
	cmd := exec.Command(command, cmdArguments...)
	// combined, _  := cmd.CombinedOutput()
	//执行命令
	cmd.Start()
	//wait 方法会一直阻塞到其所属的命令完全运行结束为止
	cmd.Wait()
}

// 逐行打印
// exec.Command("bash", "-c", "for i in 1 2 3 4;do echo $i;sleep 2;done")
// cmdArguments := []string{"-c", "for i in 1 2 3 4;do echo $i;sleep 2;done"}
// CmdPrintLineByLine("bash", cmdArguments)
func CmdPrintLineByLine(command string, cmdArguments []string)  {
	cmd := exec.Command(command, cmdArguments...)
	//创建获取命令输出管道
	stdout, err1 := cmd.StdoutPipe()
	if err1 != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err1)
		return
	}
	//执行命令
	if err2 := cmd.Start(); err2 != nil {
		fmt.Println("Error:The command is err,", err2)
		return
	}
	//使用带缓冲的读取器
	outputBuf := bufio.NewReader(stdout)
	for {
		//一次获取一行,_ 获取当前行是否被读完
		output, _, err3 := outputBuf.ReadLine()
		if err3 != nil {
			// 判断是否到文件的结尾了否则出错
			if err3.Error() != "EOF" {
				fmt.Printf("Error :%s\n", err3)
			}
			return
		}
		fmt.Printf("out:%s\n", string(output))
	}
	//wait 方法会一直阻塞到其所属的命令完全运行结束为止
	err4 := cmd.Wait()
	if err4 != nil {
		fmt.Println("wait:", err4.Error())
		return
	}
}

func ListDir(dirname string) ([]string, error) {
	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names, nil
}

// GetDirAllEntryPaths gets all the file or dir paths in the specified directory recursively.
// Note that GetDirAllEntryPaths won't follow symlink if the subdir is a symbolic link.
func GetDirAllEntryPaths(dirname string, incl bool) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	// Include current dir.
	if incl {
		paths = append(paths, dirname)
	}

	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetDirAllEntryPaths(path, incl)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// GetDirAllFilePathsFollowSymlink gets all the file or dir paths in the specified directory recursively.
func GetDirAllFilePathsFollowSymlink(dirname string, incl bool) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	// Include current dir.
	if incl {
		paths = append(paths, dirname)
	}

	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if realInfo.IsDir() {
			tmp, err := GetDirAllFilePathsFollowSymlink(path, incl)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

