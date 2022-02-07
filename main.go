package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func TrimString(s string) string {
	return strings.Trim(s, "\r\n\t ")
}

// SplitString 封装一个字符串分割函数
func SplitString(str, sep string) (stringArr []string) {
	stringArr = strings.Split(TrimString(str), sep)
	return
}

func ExecCommand(command string, output bool) (out string) {
	cmdArr := SplitString(command, " ")
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	_ = cmd.Run()
	out = string(stdout.Bytes())
	if out == "" && output {
		panic("命令执行失败")
	}
	return

}

func main() {
	flag.Parse()
	// 获取参数
	args := flag.Args()
	// 如果没有命令行为空则放出帮助信息
	if len(args) == 0 {
		fmt.Println("帮助信息：")
		fmt.Println("1. create 创建项目，可以在后面指定项目名称也可以后续手动输入。")
		fmt.Println("2. version 显示当前脚手架版本信息")
		return
	}
	action := args[0]
	switch action {
	case "create":
		var projectName string
		if len(args) == 1 {
			fmt.Print("请输入项目名称：")
			_, err := fmt.Scanln(&projectName)
			if err != nil {
				return
			}
		} else if len(args) == 2 {
			projectName = args[1]
		} else {
			fmt.Println("输入参数不正确！")
			return
		}
		fmt.Println("正在从GitHub上克隆项目，，请确保您的网络连接...")
		// 执行git clone
		ExecCommand("git clone https://github.com/allyxmiko/mvc_for_gin", false)
		// 判断下载
		_, err := os.Stat("mvc_for_gin")
		if err != nil {
			fmt.Println("下载失败！请检查网络连接！")
			return
		}
		err = os.RemoveAll("./mvc_for_gin/.git")
		if err != nil {
			fmt.Println(err)
			return
		}
		projectName = strings.ToLower(projectName)
		GetAllFiles("./mvc_for_gin", projectName)
		err = os.Rename("./mvc_for_gin", projectName)
		if err != nil {
			return
		}
		fmt.Println("项目创建完成！")
	case "vsersion":
		fmt.Println("Mgin Cli Version 1.0.0")
	default:
		fmt.Println("命令错误！不受支持的命令：", action)
	}

}

func GetAllFiles(path, projectName string) {
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		if file.IsDir() {
			GetAllFiles(path+"/"+file.Name(), projectName)
		} else {
			f, err := ioutil.ReadFile(path + "/" + file.Name())
			if err != nil {
				fmt.Println("read fail", err)
				continue
			}
			content := string(f)
			news := strings.ReplaceAll(content, "mvc_for_gin", projectName)
			var byteNew = []byte(news)
			err = ioutil.WriteFile(path+"/"+file.Name(), byteNew, 0755)
			if err != nil {
				return
			}
		}
	}
}
