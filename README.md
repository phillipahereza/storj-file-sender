# Storj File Sender

My solution for the [Storj interview question](https://gist.github.com/jtolds/0cde4aa3e07b20d6a42686ad3bc9cb53).

## Installation

This repository uses no none core packages in the main code, but the tests use [`spew`](https://github.com/davecgh/go-spew). To install these dependencies enter the following command.

```bash
go get -t
```

This repo is made up of 3 main packages and each need installing, to do so enter the following commands.

```bash
go install ./...
```

## Testing

This solution has some tests but given the involved nature of this challenge I didn't have much time to implement significant meaningful tests.

However to run what tests there are please enter the following command.

```bash
go test ./...
```

Optionally to see the outputs of the tests enter the following command.

```bash
go test -v ./...
```

## Usage

After installation and as per the specifications of the test the usage of the applications are as follows:

### Relay

```
relay :<port>
```

example
```
$ relay :9021
```

### Sender

```
sender <relay-host>:<relay-port> <file-to-send>
```

example
```
$ sender localhost:9021 corgis.mp4
```

### Receiver

```
receiver <relay-host>:<relay-port> <secret-code> <output-directory>
```

example
```
$ receiver localhost:9021 this-is-a-secret-code out/
```