# ltr

Last Trace CLI version.

The point of this project is to have a quick and easy tool to track my tasks, notes, snippets, and useful links
across all of my projects in a quick and easy way via the terminal / cli and later on TUI.

The core idea is to save all of the data in markdown files so you can use your favourite editor to write and edit them, or later on via TUI, and to have the ability of "cloud storage" by just pushing it through git.

## Development

### Build

```go
go build -o ltr main.go
go build -o ~/bin/ltr .
```

### Run

```go
ltr
```

### Config

```go
ltr config set <key> <value>
```

### Init

To use this tool, you have to init it in the
desired project / directory / repository by running:

```go
ltr init
```

### Tasks

#### New Task

```go
ltr task new <task-name>
```
