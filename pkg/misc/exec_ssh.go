package misc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SSHConfig defines SSH connection parameters.
type SSHConfig struct {
	User     string
	Host     string
	Port     int
	KeyPath  string // Path to private key file, empty for default (~/.ssh/id_rsa)
	Password string // Optional, if using password auth
}

// DefaultSSHKeyPath returns the default SSH key path (~/.ssh/id_rsa) if it exists.
func DefaultSSHKeyPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir failed: %v", err)
	}
	keyPath := filepath.Join(homeDir, ".ssh", "id_rsa")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return "", fmt.Errorf("default key path %s does not exist", keyPath)
	}
	return keyPath, nil
}

// ExecSSHCommand executes a command over SSH and returns its exit code, stdout, stderr, and any error.
func ExecSSHCommand(config SSHConfig, command string, timeout int) (int, string, string, error) {
	// 处理密钥路径
	var keyPath string
	if config.KeyPath != "" {
		keyPath = config.KeyPath
	} else {
		var err error
		keyPath, err = DefaultSSHKeyPath()
		if err != nil {
			return 125, "", "", err
		}
	}

	// SSH 认证
	var authMethod ssh.AuthMethod
	if keyPath != "" {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return 125, "", "", fmt.Errorf("read key file failed: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return 125, "", "", fmt.Errorf("parse private key failed: %v", err)
		}
		authMethod = ssh.PublicKeys(signer)
	} else if config.Password != "" {
		authMethod = ssh.Password(config.Password)
	} else {
		return 125, "", "", fmt.Errorf("no authentication method provided")
	}

	// SSH 客户端配置
	clientConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 超时控制
	baseCtx := context.Background()
	ctx := baseCtx
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(baseCtx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	// 连接 SSH
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), clientConfig)
	if err != nil {
		return 125, "", "", fmt.Errorf("ssh dial failed: %v", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		return 125, "", "", fmt.Errorf("create session failed: %v", err)
	}
	defer session.Close()

	// 设置命令，使用 setsid 创建进程组并记录 PID
	wrappedCmd := fmt.Sprintf("bash -c 'setsid %s; echo $! > /tmp/ssh_cmd_pid_%d'", command, os.Getpid())
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return 125, "", "", fmt.Errorf("capture stdout pipe failed: %v", err)
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return 125, "", "", fmt.Errorf("capture stderr pipe failed: %v", err)
	}

	// 捕获输出
	var stdoutBuf, stderrBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(&stdoutBuf, stdoutPipe)
		if err != nil {
			logrus.Errorf("copy stdout failed: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(&stderrBuf, stderrPipe)
		if err != nil {
			logrus.Errorf("copy stderr failed: %v", err)
		}
	}()

	// 启动命令
	if err := session.Start(wrappedCmd); err != nil {
		return 125, "", "", fmt.Errorf("start command failed: %v", err)
	}

	// 超时清理
	pidFile := fmt.Sprintf("/tmp/ssh_cmd_pid_%d", os.Getpid())
	if timeout > 0 {
		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded {
				session.Close()
				cleanupRemoteProcess(clientConfig, config, pidFile)
			}
		}()
	}

	// 等待命令完成
	waitErr := session.Wait()
	wg.Wait()

	// 处理退出码和错误
	var exitCode int
	var retErr error
	if waitErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			exitCode = 124
			retErr = fmt.Errorf("command timed out")
		} else if exitErr, ok := waitErr.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
			retErr = nil
		} else {
			exitCode = 125
			retErr = waitErr
		}
	} else {
		exitCode = 0
		retErr = nil
	}

	// 清理 PID 文件
	cleanupPidFile(clientConfig, config, pidFile)

	return exitCode, stdoutBuf.String(), stderrBuf.String(), retErr
}

// cleanupRemoteProcess kills the remote process group based on PID file.
func cleanupRemoteProcess(clientConfig *ssh.ClientConfig, config SSHConfig, pidFile string) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), clientConfig)
	if err != nil {
		logrus.Errorf("cleanup: ssh dial failed: %v", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logrus.Errorf("cleanup: create session failed: %v", err)
		return
	}
	defer session.Close()

	// 读取 PID 并清理进程组
	cleanupCmd := fmt.Sprintf("if [ -f %s ]; then pid=$(cat %s); kill -TERM -$pid 2>/dev/null || kill -KILL -$pid 2>/dev/null; fi", pidFile, pidFile)
	if err := session.Run(cleanupCmd); err != nil {
		logrus.Errorf("cleanup: failed to kill process group: %v", err)
	}
}

// cleanupPidFile removes the temporary PID file.
func cleanupPidFile(clientConfig *ssh.ClientConfig, config SSHConfig, pidFile string) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), clientConfig)
	if err != nil {
		logrus.Errorf("cleanup pid file: ssh dial failed: %v", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logrus.Errorf("cleanup pid file: create session failed: %v", err)
		return
	}
	defer session.Close()

	rmCmd := fmt.Sprintf("rm -f %s", pidFile)
	if err := session.Run(rmCmd); err != nil {
		logrus.Errorf("cleanup pid file: failed to remove %s: %v", pidFile, err)
	}
}
