# 🤖 Intelli-CLI

Intelli-CLI is a command-line tool that leverages the power of ChatGPT to help you find the exact command you want to execute.

![intelli-cli](https://user-images.githubusercontent.com/38521709/223257260-f283fa74-7487-4a35-8768-fa51771d8089.gif)

## 🌅 Installation

To install Intelli-CLI, run the following command:

```bash
$ brew install kj455/tap/intelli-cli
```

You can then use the -h flag to view the help information:

```bash
$ intelli-cli -h
```

## ⚙️ Setup

Before you can use Intelli-CLI, you will need to obtain an API key from [OpenAI](https://platform.openai.com/account/api-keys). Once you have your API key, run the following command and enter your key:

```bash
$ intelli-cli auth login
```

Your API key will be stored locally using [go-keyring](https://github.com/zalando/go-keyring).

## 🚀 Usage

To use Intelli-CLI, provide a description of the task you would like to accomplish:

```bash
$ intelli-cli q <description>

# Example:
$ intelli-cli q "go command to format all files in this directory"
```

Intelli-CLI will then provide you with a list of suggestions. Select the one that best fits your needs.

## License

Intelli-CLI is available under the [MIT](https://choosealicense.com/licenses/mit/) license.

