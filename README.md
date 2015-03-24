# stdslack

Pipe terminal output to Slack.

## Installation

```
$ go get github.com/sjkaliski/stdslack
```

## Usage

```
$ stdslack --token="YOUR_TOKEN"
$ cat file.txt | stdslack --channel=#mychannel
```

Once the command has completed, all output to `stdout` will be posted by user "stdslack" to your #channel.
