package disk

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"strings"
)

func Hello() {
	fmt.Println("Hello from mylibrary!")
}

var (
	stdoutBuf bytes.Buffer
	stderrBuf bytes.Buffer
)

func sshd(d string) {
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("123.com")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "192.168.1.30:9200", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
		return
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(fmt.Sprintf("df -k %s | awk 'NR==2 {print $4}'", d)) // 替换为你想要执行的命令
	if err != nil {
		log.Fatalf("执行命令失败: %v", err)
	}

	available, err := strconv.Atoi(strings.TrimSpace(stdoutBuf.String()))
	total := 10 * 1024 * 1024
	if err != nil {
		log.Fatalf("解析错误%v", err)
		return
	}
	if available > total {
		fmt.Println("可用硬盘大于10G")
	} else {
		fmt.Println("可用硬盘小于10G")
	}
	fmt.Printf("可用硬盘大小当前为%dKB\n", available)
}

func main() {
	config := flag.String("disk", "/", "a string")
	flag.Parse()
	sshd(*config)
}
