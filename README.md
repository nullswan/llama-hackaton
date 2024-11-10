# Nomi

Nowadays, people lose a lot of time trying to make the computer understand what they want to achieve.

So we created Nomi.
Nomi is an AI assistant that enables people to chat with their computer using natural language and get Nomi to execute actions for them. It gets user input, generates code from it, and then executes the newly generated code on the computer. This way, we created an assistant that can help with daily computer use.

It can open any application, launch websites, manage system settings, or interact with other software.
We are targeting companies with more than 500 employees that want to improve productivity at work.

What we are selling is a ChatGPT-like product plus the executable component on the computer.
The unique feature is natural-language-to-action.
The benefit of using Nomi is a drastic reduction in the time spent on time-consuming actions people perform every day on their personal computers.

We believe this will become the new norm and simplify people’s lives, as it’s already simplifying ours.<br />

Thank you for supporting,<br />
Swan and Ethan.


> **Note:** This project is under active development and isn't ready for full use yet. We're working hard to make it stable and reliable.
>
> We welcome any feedback, suggestions, or contributions. Thank you for trying Nomi!

- [✨ Introduction](#nomi)  
  - [🚀 Features](#-features)  
  - [🤔 Why Nomi?](#-why-nomi)  
- [🛠️ Get Started](#%EF%B8%8F-get-started)  
  - [💻 Linux & MacOS](#-linux--macos)  
  - [📟 Windows](#-windows)  
  - [🔧 Compile from Source](#-compile-from-source)  
- [🔌 Enable Providers](#-enable-providers)  
  - [🌐 Ollama](#-ollama)  
- [🗺️ Roadmap](#%EF%B8%8F-roadmap)  
- [📜 License](#-license)
- [⚙️ How to run](#-How-to-run)


### 🚀 Features

- **Versatile AI Runtime:** Lightweight and highly configurable for seamless integration.
- **Privacy-Focused:** Maintains local archives of your data, ensuring you stay in control.
- **Multi-Modal Interface:** Accepts text inputs (image support coming soon).
- **Provider Integration:** Connects with AI service Ollama.
- **Conversation Management:** Create, load, and organize conversations.
- **Prompt Engineering:** Add, edit, and manage system prompts.
- **Code Interpreter:** Run code on the fly within Nomi.
- **Terminal Experience:** Enjoy markdown-formatted output and easy command-line usage.

Explore additional features and use cases in the [Roadmap](#roadmap) section.

### 🤔 Why Nomi?

In a world where data ownership is challenging and AI is changing how we communicate, Nomi acts as a bridge between your private data and AI capabilities. It supports local provider.

While external providers involve sending data externally, Nomi works with local providers like Ollama, ensuring you retain control over your data. Our aim is to democratize AI by making it more accessible and user-friendly for everyone.

**Looking Ahead**

We're building the Nomi runtime quickly, but our journey doesn't stop there. Soon, we'll expand Nomi into a full AI platform designed to bridge the gap for non-technical users. Our goal is to make advanced AI accessible and easy to use for everyone, enabling you to benefit from AI without the need for technical expertise.

## 🛠️ Get Started

### Supported Platforms

- **Linux**: x86_64, ARM64, i686
- **MacOS**: ARM64
- **Windows**: x86_64, i686

### 🌐 Llama 3.2

You can install Ollama from [https://ollama.com/download](https://ollama.com/download) or it will be automatically installed with Nomi.

For now, we support text LLM through Ollama.

## 🗺️ Roadmap

These features are planned for future updates. They may be partially or not implemented yet.

- **Full AI Platform Development**
  - Intuitive interfaces for non-technical users
  - Expanded use case library
- **CLI Enhancements**
  - Auto-update (Update command is already available)
  - Editor mode
  - Sound on completion
- **Actions**
  - Easy transcription command
  - Presets/Projects
  - Memory tools for scripted decisions
  - Memory tools for general decisions
- **Conversation Features**
  - Markdown backup
  - New conversation types
- **Memory Enhancements**
  - Integrations
  - Use of embeddings API
- **Interpreter Updates**
  - Ask for feedback
  - Machine safety
- **File Management**
  - Real-time file management

## 📜 License

This project is licensed under the MIT License.

> See the [LICENSE](LICENSE) file for details. We believe in the power and fairness of open-source software.

## ⚙️ How to run

```
make build-dev
./dist/cli -m ollama3.2:latest
./dist/cli -m ollama3.1:latest
./dist/cli -m codellama:7b


# This is a special case, the model is flaky from the CLI
ollama run deepseek-coder-v2:latest
./dist/cli -m deepseek-coder-v2:latest
```
