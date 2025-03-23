package ssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SSHConfig содержит конфигурацию для SSH-подключения
type SSHConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// AddAuthorizedKey добавляет публичный ключ в authorized_keys на сервере
func AddAuthorizedKey(config SSHConfig, publicKey string) error {
	// Создаем SSH-клиент
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Подключаемся к серверу
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), sshConfig)
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer client.Close()

	// Создаем сессию
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("ошибка создания сессии: %w", err)
	}
	defer session.Close()

	// Добавление ключа на сервер
	cmd := fmt.Sprintf(`
		mkdir -p ~/.ssh && chmod 700 ~/.ssh &&
		touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys &&
		echo '%s' >> ~/.ssh/authorized_keys
	`, publicKey)

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("ошибка выполнения команды: %w", err)
	}

	return nil
}

// RemoveAuthorizedKey удаляет публичный ключ из authorized_keys на сервере
func RemoveAuthorizedKey(config SSHConfig, publicKey string) error {
	// Создаем SSH-клиент
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Подключаемся к серверу
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), sshConfig)
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer client.Close()

	// Создаем сессию
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("ошибка создания сессии: %w", err)
	}
	defer session.Close()

	// Создаем временный файл
	tempFile := "/tmp/authorized_keys.tmp"

	// Читаем текущий файл authorized_keys и удаляем указанный ключ
	removeKeyCmd := fmt.Sprintf(`
		if [ -f ~/.ssh/authorized_keys ]; then
			grep -v "%s" ~/.ssh/authorized_keys > %s
			mv %s ~/.ssh/authorized_keys
			chmod 600 ~/.ssh/authorized_keys
		fi
	`, publicKey, tempFile, tempFile)

	if err := session.Run(removeKeyCmd); err != nil {
		return fmt.Errorf("ошибка удаления ключа: %w", err)
	}

	return nil
}

// ValidatePublicKey проверяет корректность публичного ключа
func ValidatePublicKey(publicKey string) error {
	// Проверяем, что ключ начинается с правильного префикса
	if !strings.HasPrefix(publicKey, "ssh-rsa ") && !strings.HasPrefix(publicKey, "ssh-ed25519 ") {
		return fmt.Errorf("неверный формат публичного ключа")
	}

	// Проверяем, что ключ не пустой
	if len(strings.TrimSpace(publicKey)) == 0 {
		return fmt.Errorf("пустой публичный ключ")
	}

	return nil
}
