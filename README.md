# Introduction
P2P chat is an implementation of peer-to-peer messaging and presence protocol in Golang. Messages are exchanged between the client and server in the JSON format.

# Requirements for Execution
You will need Golang compiler (specific to your operating system) for compiling the program code. A compiler is available for each of Windows, Mac OS and other Linux based operating systems. Following the instructions below to install the compile on your system:

## Windows Machine
- Download the installer using [this link](https://dl.google.com/go/go1.11.2.windows-amd64.msi)
- Running the installer with install the compiler and configure the required environment variables

## MacOS Machine
- Download the installer using [this link](https://dl.google.com/go/go1.11.2.darwin-amd64.pkg)
- Running the installer will install the compiler in `/usr/local/go/bin` and configure the required environment variables

## Linux Machine
- Download the archive using [this link](https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz)
- Run `tar -C /usr/local -xzf go1.11.2.linux-amd64.tar.gz`, to extract the Go compiler in `/usr/local/go` directory
- Add `/usr/local/go/bin` to the PATH environment variable by adding `export PATH=$PATH:/usr/local/go/bin` line in `$HOME/.profile` file

# Execution Steps

## Installing dependancies
The project is dependant on [webview](https://github.com/zserge/webview). Install it using:
`https://github.com/zserge/webview`

## Compiling the code
- Open a terminal/command prompt and `cd` to the project directory
- Run command `go build` to compile the source code and build an executable

## Updating config file
The config file is `config.json` contains three parameters, which are as follows:
- **port**: It is the port on which the server should run. It is recommended to use port numbers greater than 9000, so that it does not conflict with the standard ports used by other applications
- **name**: It is the name with which your instance will be registered in the network
- **knownHosts**: It is an array pre-available list of known hosts. The format of host string should be like <IP Address>:<Port Number>, e.g. 127.0.0.1:9356

If you are not the first one to start the network, then you must add the hostname and port of another running instance. Information about other known peers in the network will be sent by that instance.

## Running the program
- If you are on the windows machine, run p2p.exe
- If you are on MacOS or any other Linux based operating system, the by executing `./p2p`