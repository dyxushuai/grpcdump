# gRPC dump tool

[中文文档](README-zh.md)

```sh
» sudo go run cmd/grpcdump/main.go -d lo0 -port 8085 -ip 127.0.0.1 -proto ./grpc_example/helloworld/helloworld/helloworld.proto
2018/06/06 21:18:02 Starting capture on device "lo0"
2018/06/06 21:18:02 reading in packets

REQUEST(STREAM=1) > 2018-06-06T21:18:04.921128144+08:00: 127.0.0.1:56327 ---> 127.0.0.1:8085
  HEADERS:
    :method = "POST"
    :scheme = "http"
    :path = "/helloworld.Greeter/SayHello"
    :authority = "127.0.0.1:8085"
    content-type = "application/grpc"
    user-agent = "grpc-go/1.13.0-dev"
    te = "trailers"
  BODY:
    {"name":"world","i18":"zh"}

RESPONSE(STREAM=1) > 2018-06-06T21:18:04.921691304+08:00: 127.0.0.1:56327 <--- 127.0.0.1:8085
  HEADERS:
    :status = "200"
    content-type = "application/grpc"
  BODY:
    {"message":"Hello, 你好"}

RESPONSE(STREAM=1) > 2018-06-06T21:18:04.921865543+08:00: 127.0.0.1:56327 <--- 127.0.0.1:8085
  HEADERS:
    grpc-status = "0"
    grpc-message = ""
```

## Requirement

Build httpdump requires libpcap-dev and cgo enabled.

### libpcap

for ubuntu/debian:

```sh
sudo apt install libpcap-dev
```

for centos/redhat/fedora:

```sh
sudo yum install libpcap-devel
```

for osx:

Libpcap and header files already installed.

## Install

```sh
go get -u github.com/dyxushuai/grpcdump/cmd/grpcdump
```

## Usage

```sh
-assembly_debug_log
        If true, the github.com/google/gopacket/tcpassembly library will log verbose debugging information (at least one line per packet)
  -assembly_memuse_log
        If true, the github.com/google/gopacket/tcpassembly library will log information regarding its memory use every once in a while.
  -d string
        Interface to get packets from (default "eth0")
  -ip string
        Filter by ip, if either source or target ip is matched, the packet will be processed
  -port uint
        Filter by port, if either source or target port is matched, the packet will be processed.
  -proto string
        Protobuf spec file
  -v    Logs every packet in great detail
```

## Wanted

- [x] Dump bytes from Network Interface, e.g. "eth0".
- [x] Filter by host ip and host port.
- [x] Parse bytes as HTTP2 protocol.
- [ ] Decrypt TLS sessions by `key log`.
- [x] Dynamic reflect the protobuf files at runtime.
- [x] Pretty print.