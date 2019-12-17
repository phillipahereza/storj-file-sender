# Storj File Sender

My solution for the [Storj interview question](https://gist.github.com/jtolds/0cde4aa3e07b20d6a42686ad3bc9cb53).

## Contents

- [Installation](#installation)
- [Testing](#testing)
- [Usage](#usage)
  - [Relay](#relay)
  - [Sender](#sender)
  - [Receiver](#receiver)
- [Justifications](#justifications)
  - [codegen](#codegen)
  - [Checksuming](#checksuming)
  - [TCP](#tcp)
  - [Data Terminator](#data-terminator)
  - [Headers](#headers)
  - [Concurrency](#concurrency)
- [Improvements](#Improvements)

## Installation

This repository uses no none core packages in the main code, but the tests use [`spew`](https://github.com/davecgh/go-spew). To install this dependency enter the following command.

```bash
$ go get -t
```

This repo is made up of 3 main packages and each need installing, to do so enter the following commands.

```bash
$ go install ./...
```

## Testing

This solution has some tests but given the involved nature of this challenge I didn't have much time to implement significant meaningful tests.

However to run what tests there are please enter the following command.

```bash
$ go test ./...
```

Optionally to see the outputs of the tests enter the following command.

```bash
$ go test -v ./...
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
delicious-bisque-jaguar-382
```

### Receiver

```
receiver <relay-host>:<relay-port> <secret-code> <output-directory>
```

example
```
$ receiver localhost:9021 delicious-bisque-jaguar-382 out/
```

## Justifications

In this repo there are a number of choices that I've made and I wish to highlight them and justify their existence.

### codegen
A package of the `Sender` application is `codegen`. This little package's purpose is to create a randomly generated string that is sufficiently unique and also sufficiently easy to convey over the phone.

>The user of laptop A wants to send a file to the user of laptop B. The users are talking on the phone, so they can exchange information, but not the file itself.

Because the specification states that the users are talking on the phone I tried to make the secret code something that could be easy to remember and easy to convey. So I chose a [passphrase-esque approach](https://xkcd.com/936/).

In order to get as random result as possible I sum the `BigEndian` sum of the file hash with the unix time, to form the random seed. This means that the chance of two simultaneous senders getting the same secret code for different files is incredibly small. From the pool of 420 animals, 852 colours, 1133 adjectives plus a number between 0 and 10,000 there are 4,054,732,632,720‬‬ possible combinations of secret codes.

The kinds of secret codes you will get are things like, `magical-liver-mouse-deer-311`, `mundane-chestnut-tarantula-989` and `political-lemon-glacier-lynx-348`.

### Checksuming
Although TCP has native checksuming and Go is a super language and during my testing and developing I haven't seen any data that got mashed up or eaten on my local machine, but some networks could do this. So from the beginning I developed this system with the understanding that part of the data transfer would include a checksum of the original file, to ensure that the integrity of the file was determinable at the `Receiver's` side.  

### TCP
>Your relay program should not use much memory. It should not use more than 4MB of memory or storage per transfer, regardless of the size of the file being transfered.

This stood out to me in the specifications and so I immediately thought of using TCP for transferring the data between the remote parties. TCP is an excellent data streaming protocol and for this reason I chose TCP to send files of an unknown size. 

### Data Terminator
Something that I spent a worrying amount of time on was the copying of `net.TCPConns` via the `io.Copy` function. My problem was that `io.Copy` will never end its internal `for{}` unless one of the connections closed. This causes, or at least in my development, the `io.Copy(conn, sConn)` to hang indefinitely because neither connection can be closed with both are in the `io.Copy`.

To fix this I implemented a data terminator into the protocol to notify the `Relay` when the body of the data had reached its end. As stated in the comments most of the implementation is just copy and paste of the core `io.Copy` func, the can code be [seen here](https://github.com/Samyoul/storj-file-sender/blob/master/relay/main.go#L156).

### Headers
This challenge required more than just the raw file data to be transferred it also required that additional meta data be transferred. To facilitate the sending of additional meta data I implemented functionality for the creation and handling of data headers.

The code for this can be [seen here](https://github.com/Samyoul/storj-file-sender/blob/master/common/header.go).

**Requests :**

The basic request header has the following structure:

|Name|Bytes|
|:---:|:---:|
|`Type`|1|
|`Code`|64|

A `send` request header has the following additional fields

|Name|Bytes|
|:---:|:---:|
|`Checksum`|32|
|`Filename`|Until data terminator|

Because filename was the only piece of meta data that was required to be sent from sender to receiver and was the only variable that could be of any length the filenames are appended to the end of the header and are delimited from the data stream via the use of the data terminator.

**Responses**

The `Receiver` is the only application that receives a data header in its response as it requires the `Filename` and the `Checksum`. The structure of this header is the same as the additional header data from the send request.

|Name|Bytes|
|:---:|:---:|
|`Checksum`|32|
|`Filename`|Until data terminator|
 
### Concurrency
This test requires a few concurrency techniques because the specifications require:

>Your relay program should support multiple people sending files at the same time.

Concurrency would be required even if multiple people sending files wasn't a requirement as both the `Sender` and `Receiver` have concurrent TCP connections with the `Relay`. So my approach of **pushing the connection handler into a go routine** can allow for `n` number of concurrent connections.

Now that we have multiple concurrent connections we need a way to ensure data can be transferred between them. To allow this I implemented a simple `stream` struct that was mapped to a `streamMap`. The `stream` struct holds a `chan net.Conn` to allow for data streaming between concurrent processes.

I added a `sync.WaitGroup` to the `stream` struct to allow the `Relay` to know when was appropriate to close the send connection.

## Improvements
Given the time restrictions on this task there are a lot of this that could be improved on.

- Better test coverage
- More in depth tests
- Better error handling between connections. Functionality to allow the `Relay` to inform the `Sender` and `Receiver` specifically why something went wrong on the `Relay` and then display this error to the user of either the `Sender` or the `Receiver`. This could be acheived with expanding the Header protocol.
- Research into potentially more efficient methods of transferring data between two concurrent connections.
- Add a `-h` flag to the CLI of each application, to allow the user to discover what arguments are expect without failing.
- Lower chance of secret code collision via higher entropy.
