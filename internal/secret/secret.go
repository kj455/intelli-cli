package secret

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/kj455/intelli-cli/internal/prompt"
	"github.com/zalando/go-keyring"
)

const (
	service = "openai-api-key"
	user    = "intelli-cli"
)

func GetApiKey() (string, error) {
	return keyring.Get(service, user)
}

func DeleteApiKey() error {
	return keyring.Delete(service, user)
}

func setApiKey(value string) error {
	err := keyring.Set(service, user, value)
	if err != nil {
		return fmt.Errorf("failed to set secret key: %w", err)
	}
	return nil
}

func openApiKeySettingPage() error {
	const URL = "https://platform.openai.com/account/api-keys"

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", URL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", URL).Start()
	case "darwin":
		err = exec.Command("open", URL).Start()
	default:
		err = errors.New("unsupported platform")
	}
	return err
}

func SetupSecretIfNeeded() error {
	_, err := GetApiKey()

	if err != nil {
		err := openApiKeySettingPage()
		if err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}

		key, err := prompt.PromptApiKey()
		if err != nil {
			return fmt.Errorf("failed to get api key: %w", err)
		}
		fmt.Println("🚀 Your API key was successfully set!")

		err = setApiKey(key)
		if err != nil {
			return fmt.Errorf("failed to set api key: %w", err)
		}
	}

	return nil
}
