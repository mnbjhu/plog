# Pretty Logger (plog)

A TUI to watch application logs in the terminal

## Usage

```
plog <command>
```

E.g.

```bash
plog java -jar app.jar
```

## Install

### Using Go Package Manager

```bash
go install github.com/mnbjhu/plog
```

### Download The Binary

Binaries for each target can be found on the releases page [here](https://github.com/mnbjhu/plog/releases/latest).

## Custom Logging Format

Plog will expect logs in the default log4j format, however you can configure it to use a custom format.
Create a file in the project directory named `.plog.json` with your config like so:

```json
{
  "regex": "^(?<date>[^ ]+)  (?<level>[^ ]+) (?<pkg>[^ ]+): (?<message>.*)$",
  "columns": [
    {
      "name": "date",
      "width": 10
    },
    {
      "name": "level",
      "width": 5
    },
    {
      "name": "pkg",
      "width": 8
    },
    {
      "name": "message",
      "width": 20
    }
  ],
  "input": "stderr"
}
```
