# Reload: Hot Reloading for Go Applications

## Introduction
Reload is a dynamic tool designed to enhance the development workflow for Go applications. It enables hot reloading of Go applications, automatically rebuilding your project whenever changes are detected. This means immediate feedback and a more efficient development process.

## Features
- **Hot Reloading**: Automatically rebuilds your Go application when a change is detected in any file within the current working directory (CWD).
- **Fault Tolerant**: Simply put your application wont stop working if you have unfinished and files. If there is an error in your current code, reload will keep serving your last working binary unless all complier error are fixed.
- **Recursive Monitoring**: Watches for changes in all files recursively within the current working directory.
- **Easy to Use**: Simple command-line interface, requiring minimal setup.

## Getting Started

### Prerequisites
- Go installed on your system (visit [Go's official website](https://golang.org/dl/) for installation instructions).

### Installation
To install Reload, you can use `go get`:

```bash
go get -u github.com/hhacker1999/reload

go build -o reload cmd/server/main.go

sudo mv reload /usr/local/bin/ 
```

### Usage
To use Reload, navigate to your project directory and run:
```bash
reload -p path/to/main.go
```
Where `path/to/main.go` is the path to the main file of your Go application.

Reload will now monitor for file changes in your project directory and its subdirectories. When a change is detected, Reload will automatically rebuild and restart your application.

## Configuration
Currently, Reload requires only the `-p` flag to specify the path to the main file. Future versions may include additional configuration options.


## Upcoming Features
- **Skip certain files**: Reload currently skips .git repository but you may want certain files to be skipped in your project.
- **Pre build step**: Users should be able to specify a pre build command like generating protoc files if a proto file is changed.

## Contributing
Contributions are welcome! Feel free to fork the repository and submit pull requests.

## License
This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details.
