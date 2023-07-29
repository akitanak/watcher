# watcher

 Watcher CLI Tool is a command-line utility designed to monitor a specified directory and its subdirectories for any file modifications. Once a change is detected, Watcher automatically executes a user-specified command, providing a flexible and automated solution for various tasks.

## Install

```
go install github.com/akitanak/watcher
```

## Usage

If you want to watch current directory and run `go test` when some code was changed, run a command below.
```
watcher go test
```

If you want to monitor a specific directory, specify the directory with `--directory`.
```
watcher --directory ./src go test
```
