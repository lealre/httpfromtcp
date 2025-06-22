# httpfromtcp

## Communication Mechanisms for Processes

| Feature              | TCP                                | Unix Domain Socket                    | UDP                                     |
| -------------------- | ---------------------------------- | ------------------------------------- | --------------------------------------- |
| **Type**             | Network Protocol (Transport Layer) | IPC Mechanism (Filesystem-based)      | Network Protocol (Transport Layer)      |
| **Speed**            | Medium (reliability overhead)      | ⚡ **Fastest** (kernel bypass)        | **Fast** (no connection setup)          |
| **Reliability**      | ✅ Guaranteed (in-order, no loss)  | ✅ Guaranteed (local only)            | ❌ Best-effort (possible loss/disorder) |
| **Connection Model** | Connection-oriented (stateful)     | Connection-oriented (file-based)      | Connectionless (stateless)              |
| **Addressing**       | IP + Port (e.g., `192.168.1.1:80`) | Filesystem Path (e.g., `/tmp/sock`)   | IP + Port (e.g., `10.0.0.1:53`)         |
| **Error Handling**   | Retransmits lost packets           | N/A (local, no packet loss)           | No retransmission                       |
| **Flow Control**     | ✅ (sliding window)                | N/A                                   | ❌ (no built-in control)                |
| **Use Cases**        | Web (HTTP/1-2), SSH, email (SMTP)  | Databases, Containers, Local Services | VoIP, DNS, Gaming, HTTP/3 (QUIC)        |
| **Overhead**         | High (headers, handshakes, ACKs)   | Minimal (no network stack)            | Low (8-byte header)                     |
| **OSI Layer**        | Layer 4 (Transport)                | Layer 0 (Bypasses network stack)      | Layer 4 (Transport)                     |

- [TCP listener Example](cmd/tcplistener/main.go)
- [UDP sender Example](cmd/udpsender/main.go)

## HTTP Version Comparison

| Version  | Transport       | Format        | Header Compression | Multiplexing | Key Innovation                              |
| -------- | --------------- | ------------- | ------------------ | ------------ | ------------------------------------------- |
| HTTP/1.1 | TCP             | Plain text    | ❌ None            | ❌ No        | Standardized persistent connections         |
| HTTP/2   | TCP             | Binary frames | ✅ HPACK           | ✅ Yes       | Framing + multiplexing over TCP             |
| HTTP/3   | QUIC (over UDP) | Binary frames | ✅ QPACK           | ✅ Yes       | Fixes TCP's HOL blocking at transport layer |

### HTTP/1.1 (Text)

```http
start-line CRLF
*( field-line CRLF )
CRLF
[ message-body ]
```

- Human-readable, line-based

- Headers repeated for each request (no compression)

##### Example GET

1 - Run in one terminal

```shell
go run ./cmd/tcplistener | tee /tmp/rawget.http
```

2 - Run in another terminal

```shell
curl http://localhost:42069/coffee
```

3 - Kill the second process (`curl`) and then the first process (`tcplistener`)

4 - Check the content in `/tmp/rawget.http`

```shell
cat /tmp/rawget.http
```

Expected

```http
GET /goodies HTTP/1.1       # start-line CRLF
Host: localhost:42069       # *( field-line CRLF )
User-Agent: curl/7.81.0     # *( field-line CRLF )
Accept: */*                 # *( field-line CRLF )
                            # CRLF
                            # [ message-body ] (empty)
```

##### Example POST

1 - Run in one terminal

```shell
go run ./cmd/tcplistener | tee /tmp/rawpost.http
```

2 - Run in another terminal

```shell
curl -X POST -H "Content-Type: application/json" -d '{"flavor":"dark mode"}' http://localhost:42069/coffee
```

3 - Kill the second process (`curl`) and then the first process (`tcplistener`)

4 - Check the content in `/tmp/rawpost.http`

```shell
cat /tmp/rawpost.http
```

Expected

```http
POST /coffee HTTP/1.1            # start-line CRLF
Host: localhost:42069            # *( field-line CRLF )
User-Agent: curl/8.6.0           # *( field-line CRLF )
Accept: */*                      # *( field-line CRLF )
Content-Type: application/json   # *( field-line CRLF )
Content-Length: 22               # *( field-line CRLF )
                                 # CRLF
{"flavor":"dark mode"}          # [ message-body ]
```

### HTTP/2/3 (Binary)

Frames replace raw text. Example frame types:

- **HEADERS** (compressed HTTP headers)

- **DATA** (response body chunks)

- **SETTINGS** (configuration)

**HPACK/QPACK**: Encodes headers as binary (e.g., **:method:** GET → a tiny integer)