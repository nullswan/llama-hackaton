```
make build-dev
./dist/cli -m ollama3.2:latest
./dist/cli -m ollama3.1:latest
./dist/cli -m codellama:7b


# This is a special case, the model is flaky from the CLI
ollama run deepseek-coder-v2:latest
./dist/cli -m deepseek-coder-v2:latest
```
