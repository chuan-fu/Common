package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// Md5 md5加密
func Md5(encodeString string) string {
	h := md5.New()
	h.Write([]byte(encodeString))
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果
}

// TODO 未测试过
// ValidateRemoteAddr 判断ip端口是否合法
func ValidateRemoteAddr(ip string) bool {
	match, err := regexp.MatchString(`^(?:(?:1[0-9][0-9]\.)|(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))\:(([0-9])|([1-9][0-9]{1,3})|([1-6][0-9]{0,4}))$`, ip)
	if err != nil {
		return false
	}
	return match
}

// TODO 未测试过
// ValidateURL 判断ip端口是否合法
func ValidateURL(url string) bool {
	match, err := regexp.MatchString(`^/(([a-zA-Z][0-9a-zA-Z+\-\.]*:)?/{0,2}[0-9a-zA-Z;/?:@&=+$\.\-_!~*'()%]+)?(#[0-9a-zA-Z;/?:@&=+$\.\-_!~*'()%]+)?$`, url)
	if err != nil {
		return false
	}
	return match
}

// TODO 未测试过
// Stop 关闭网关服务，重启读取配置文件
func Stop() bool {
	id := os.Getpid()
	cmd := exec.Command("/bin/bash", "-c", "kill -HUP "+strconv.Itoa(id))
	if _, err := cmd.Output(); err != nil {
		return false
	}
	return true
}

// TODO 未测试过
// GetMac 获取MAC地址
func GetMac() (bool, string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, "Poor soul, here is what you got: " + err.Error()
	}
	for k := range interfaces {
		mac := interfaces[k].HardwareAddr // 获取本机MAC地址
		m := fmt.Sprintf("%s", mac)
		match, err := regexp.MatchString(`[0-9a-f][0-9a-f][:-][0-9a-f][0-9a-f][:-][0-9a-f][0-9a-f][:-][0-9a-f][0-9a-f][:-][0-9a-f][0-9a-f][:-][0-9a-f][0-9a-f]`, m)
		if err != nil {
			return false, ""
		}
		if match {
			return true, m
		}
	}
	return false, ""
}
