package disk

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

var (
	stdoutBuf bytes.Buffer
	stderrBuf bytes.Buffer
	total     int = 10 * 1024 * 1024
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Hard     string
}

type Directory struct {
	Size int
}

func InitSSH(c *Config) (Directory, error) {
	var tech Directory
	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, &c.Port), config)
	if err != nil {
		return Directory{}, err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return Directory{}, err
	}
	defer session.Close()
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(fmt.Sprintf("df -k %s | awk 'NR==2 {print $4}'", c.Hard)) // 替换为你想要执行的命令
	if err != nil {
		return Directory{}, err
	}

	available, err := strconv.Atoi(strings.TrimSpace(stdoutBuf.String()))
	if err != nil {
		return Directory{}, err
	}
	tech = Directory{Size: available}
	return tech, nil
	//if available > total {
	//	fmt.Println("可用硬盘大于10G")
	//} else {
	//	fmt.Println("可用硬盘小于10G")
	//}
	//fmt.Printf("可用硬盘大小当前为%dKB\n", available)
}

func (t Directory) GB() int {
	return t.Size / 1024 / 1024
}

func (t Directory) KB() int {
	return t.Size
}

func (t Directory) MB() int {
	return t.Size / 1024
}
