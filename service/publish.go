package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// VideoToImage  从视频中提取封面图片保存函数
// ps: 需要额外安装ffmpeg http://ffmpeg.org/download.html
func VideoToImage(videoPath string, toSavePath string) {
	arg := []string{"-hide_banner"}
	arg = append(arg, "-i", videoPath)
	arg = append(arg, "-r", "1")
	arg = append(arg, "-f", "image2")
	arg = append(arg, "-frames:v", "1") // 截取一张
	arg = append(arg, "-q", "8")        // 设置图片压缩等级，越高压缩越大
	arg = append(arg, toSavePath)
	// 通过命令行运行ffmpeg截取视频帧图片保存为封面图
	cmd := exec.Command("ffmpeg", arg...)
	cmd.Stderr = os.Stderr
	fmt.Println("Run", cmd)
	err := cmd.Run()
	if err != nil {
		return
	}
	fmt.Println("提取视频封面图成功！")
}

// GetSavedUrlAddress 获取视频和图片的保存URL地址
func GetSavedUrlAddress(toSaveFilePath string) string {
	// 根据操作系统自动判断分隔符
	sysType := runtime.GOOS
	var sysSpliter string
	if sysType == "windows" {
		sysSpliter = "\\"
	} else {
		sysSpliter = "/"
	}
	// 分割获取文件名
	fileSlice := strings.Split(toSaveFilePath, sysSpliter)
	fileName := fileSlice[len(fileSlice)-1]
	// 拼接返回文件存储的URL地址
	var bt bytes.Buffer
	bt.WriteString("http://")
	bt.WriteString(IpAddress)
	bt.WriteString(":8080/static/")
	bt.WriteString(fileName)
	fileUrlAddress := bt.String()
	return fileUrlAddress
}
