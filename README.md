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

By running it created `.ltr` folder and inside `/tasks`, `/snippets` and `/notes` folders.

```go
ltr init
```

### Tasks

#### List Tasks

On default it lists just uncompleted tasks.

```go
ltr task list
```

List only completed tasks

```go
ltr task list -c
```

List all tasks

```go
ltr task list -a
```

#### New Task

```go
ltr task new <task-name>
```

#### Delete task

```go
ltr task delete <task-id>
```

#### Open task in neovim

```go
ltr task open <task-id>
```

#### Complete task

```go
ltr task complete <task-id>
```

### Notes

#### List notes

```go
ltr note list
```

#### New note

```go
ltr note new <note-name>
```

#### Open note

```go
ltr note open <note-id>
```

#### Delete note

```go
ltr note delete <note-id>
```

### Snippets

#### List snippets

```go
ltr snippet list
```

#### New snippet

```go
ltr snippet new <snippet-name>
```

#### Open snippet

```go
ltr snippet open <snippet-id>
```

#### Delete snippet

```go
ltr snippet delete <snippet-id>
```
