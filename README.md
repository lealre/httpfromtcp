# httpfromtcp

This repository is the implementation of an HTTP/1.1 server using TCP connections, inspired by the course [Learn the HTTP Protocol in Go, from boot.dev](https://www.boot.dev/courses/learn-http-protocol-golang).

The main implementation is located inside the `internal` folder, while the `cmd/httpserver` directory demonstrates its usage.

Inside the `cmd` folder, there are also:

- [TCP listener Example](cmd/tcplistener/main.go)
- [UDP sender Example](cmd/udpsender/main.go)

## Table of Contents

1. [HTTP/1.1 (Text) - Used here](#http11-text---used-here)
2. [Example GET - Request](#example-get---request)
3. [Example POST - Request](#example-post---request)
4. [Request Structure](#request)
5. [Headers in HTTP](#headers)
6. [Body in HTTP](#body)
7. [Response Structure](#response)
8. [Chunked Encoding](#chunked-encoding)
   - [Trailers](#trailers)
9. [Binary Data Handling](#binary-data)
10. [HTTP/2 Overview](#http2)
11. [HTTP/3 Overview](#http3)
12. [HTTP Version Comparison](#http-version-comparison)
13. [Communication Mechanisms for Processes](#communication-mechanisms-for-processes)

## HTTP/1.1 (Text) - Used here

Since HTTP uses TCP, an HTTP request or response may be broken into multiple TCP packets if it is too large to fit into a single packet. These packets are then reassembled in the correct order by the receiver.

TCP guarantees that the data is delivered in order and intact.

At the heart of HTTP is the `HTTP-message`, which defines the format that the text in an HTTP request or response must have. From [RFC 9112 Section 2.1](https://datatracker.ietf.org/doc/html/rfc9112#name-message-format):

```text
start-line CRLF
*( field-line CRLF )
CRLF
[ message-body ]
```

- Human-readable, line-based

- Headers repeated for each request (no compression)

[CRLF](https://developer.mozilla.org/en-US/docs/Glossary/CRLF) (written in plain text as `\r\n`) is a carriage return followed by a line feed. It's a Microsoft Windows (and HTTP) style newline character.

| Part                  | Example                   | Description                                                    |
| --------------------- | ------------------------- | -------------------------------------------------------------- |
| start-line CRLF       | `POST /users HTTP/1.1`    | The request (for a request) or status (for a response) line    |
| \*( field-line CRLF ) | `Host: google.com`        | Zero or more lines of HTTP headers. These are key-value pairs. |
| CRLF                  |                           | A blank line that separates the headers from the body.         |
| [ message-body ]      | `{"name": "TheHTTPuser"}` | The body of the message. This is optional.                     |

Currently, there are several key RFCs for HTTP/1.1:

- [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231) – An active and widely referenced RFC.
- [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) – Easier to read than RFC 7231, relies on understanding from RFC 9110.
- [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110) – Covers HTTP "semantics."

### Example GET - Request

1. Run in one terminal:

```shell
go run ./cmd/tcplistener | tee /tmp/rawget.http
```

2. Run in another terminal:

```shell
curl http://localhost:42069/coffee
```

3. Terminate the second process (`curl`) and then the first process (`tcplistener`).

4. Check the content in `/tmp/rawget.http`:

```shell
cat /tmp/rawget.http
```

Expected output:

```text
GET /goodies HTTP/1.1       # start-line CRLF
Host: localhost:42069       # *( field-line CRLF )
User-Agent: curl/7.81.0     # *( field-line CRLF )
Accept: */*                 # *( field-line CRLF )
                            # CRLF
                            # [ message-body ] (empty)
```

### Example POST - Request

1. Run in one terminal:

```shell
go run ./cmd/tcplistener | tee /tmp/rawpost.http
```

2. Run in another terminal:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"flavor":"dark mode"}' http://localhost:42069/coffee
```

3. Terminate the second process (`curl`) and then the first process (`tcplistener`).

4. Check the content in `/tmp/rawpost.http`:

```shell
cat /tmp/rawpost.http
```

Expected output:

```text
POST /coffee HTTP/1.1            # start-line CRLF
Host: localhost:42069            # *( field-line CRLF )
User-Agent: curl/8.6.0           # *( field-line CRLF )
Accept: */*                      # *( field-line CRLF )
Content-Type: application/json   # *( field-line CRLF )
Content-Length: 22               # *( field-line CRLF )
                                 # CRLF
{"flavor":"dark mode"}          # [ message-body ]
```

### Request

When it's a request, the start-line is referred to as the `request-line` and follows a specific format.

```text
HTTP-version  = HTTP-name "/" DIGIT "." DIGIT
HTTP-name     = %s"HTTP"
request-line  = method SP request-target SP HTTP-version
```

For `HTTP/1.1`, an example request-line looks like this:

```text
GET /coffee HTTP/1.1
```

### Headers

- The RFC uses the term `field-line`, which essentially refers to the same concept as headers. From [5. Field Syntax](https://datatracker.ietf.org/doc/html/rfc9112#name-field-syntax):

"_Each field line consists of a case-insensitive field name followed by a colon (":"), optional leading whitespace, the field line value, and optional trailing whitespace._"

```text
field-line   = field-name ":" OWS field-value OWS
```

- An unlimited amount of whitespace can be present before and after the field-value (Header value).
- There must be no spaces between the colon and the field-name. Example:

```text
// valid
'Host: localhost:42069'
' Host: localhost:42069 '

// not valid
Host : localhost:42069
```

- Field names (keys) are case-insensitive
- `field-name` has an implicit definition of a token as defined by [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110).
- `field-name` must contain only:
  - Uppercase letters: A-Z
  - Lowercase letters: a-z
  - Digits: 0-9
  - Special characters: !, #, $, %, &, ', \*, +, -, ., ^, \_, `, |, ~
  - Must have at least a length of 1.
- It's valid (see [RFC 9110 5.2](https://datatracker.ietf.org/doc/html/rfc9110#name-field-lines-and-combined-fi)) to have multiple values for a single header key. For example:

```text
Set-Person: user1;
Set-Person: user2;
Set-Person: user3;

// Should combine the values into a single string, separated by a comma
Set-Person: user1, user2, user3
```

### Body

According to [RFC9110 8.6](https://datatracker.ietf.org/doc/html/rfc9110#section-8.6):

```text
A user agent SHOULD send Content-Length in a request...
```

And "should" has a specific meaning in RFCs per [RFC2119](https://datatracker.ietf.org/doc/html/rfc2119#section-3):

```text
This word, or the adjective "RECOMMENDED", means that there may exist valid reasons in particular circumstances to ignore a particular item, but the full implications must be understood and carefully weighed before choosing a different course.
```

### Response

- HTTP responses follow the same HTTP message format:

```text
 HTTP-message   = start-line CRLF
                  *( field-line CRLF )
                  CRLF
                  [ message-body ]
```

The only difference is that the `start-line` is a `status-line` instead of a `request-line`. From [RFC 9112](https://www.rfc-editor.org/rfc/rfc9112.html):

```text
status-line = HTTP-version SP status-code SP [ reason-phrase ]
```

- Example:

```text
HTTP/1.1 200 OK

// or...
HTTP/1.1 404 Not Found
```

- Regarding the reason-phrase, from Section 4,

"_A client SHOULD ignore the reason-phrase content because it is not a reliable channel for information (it might be translated for a given locale, overwritten by intermediaries, or discarded when the message is forwarded via other versions of HTTP). A server MUST send the space that separates the status-code from the reason-phrase even when the reason-phrase is absent (i.e., the status-line would end with the space)._"

So, while reason phrases are typically included (and match one to one with the status code), they are not required and should be ignored by clients.

### Chunked Encoding

The `[ message-body ]` can be a bit deceiving; it is a rather flexible field that can contain a variable length of data, known only as it is sent using the [`Transfer-Encoding`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Transfer-Encoding) header rather than the [`Content-Length`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Length) header. Here's the format:

```text
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
... repeat ...
0\r\n
\r\n
```

Where `<n>` is a hexadecimal number indicating the size of the chunk in bytes and `<data of length n>` is the actual data for that chunk. This pattern can be repeated as many times as necessary to send the entire message body. Here's a concrete example with plain text:

```text
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

1E
I could go for a cup of coffee
B
But not Java
12
Never go full Java
0
```

Chunked encoding is most often used for:

- Streaming large amounts of data (like big files)
- Real-time updates (like a chat-style application)
- Sending data of unknown size (like a live feed)

#### Trailers

Additional headers can be appended at the end of chunked encoding, called [Trailers](https://datatracker.ietf.org/doc/html/rfc9112#section-7.1.2). They function similarly to headers but with one important distinction: the trailer names must be specified in a Trailer header. For example:

```text
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked
Trailer: Lane, Prime, TJ

1E
I could go for a cup of coffee
B
But not Java
12
Never go full Java
0

0\r\n
Lane: goober
Prime: chill-guy
TJ: 1-indexer
\r\n
```

Trailers are often used to send information about the message body that cannot be determined until the message body is fully transmitted. For example, the hash of the message body.

### Binary Data

HTTP is a text-based protocol, but it is quite proficient at transmitting binary data. This is achieved through the use of the Content-Type header to specify the type of data being sent. For instance, if you're transmitting an image, you might use `Content-Type: image/png`, or for a video, `Content-Type: video/mp4`.

This helps the client interpret the body, allowing it to expect a specific format (like video data) rather than interpreting it as raw text.

### HTTP/2

Some key differences include:

- It is a binary protocol rather than text-based. This makes it more efficient and less error-prone, but it also requires more steps to debug, typically.
- HTTP/2 uses multiplexing, allowing multiple requests and responses to be sent over a single connection simultaneously. This reduces latency and improves performance.
- It employs header compression to save on header bandwidth.
- It supports server push, enabling the server to send resources to the client before they are requested.

### HTTP/3

Built on QUIC instead of TCP, QUIC is a transport layer protocol that operates over UDP, allowing for faster connection establishment and improved performance.
HTTP/3 mandates encryption at the HTTP protocol level (HTTP/1.1 is unencrypted by default, meaning it requires HTTPS to be secure).

### HTTP Version Comparison

| Version  | Transport       | Format        | Header Compression | Multiplexing | Key Innovations                             |
| -------- | --------------- | ------------- | ------------------ | ------------ | ------------------------------------------- |
| HTTP/1.1 | TCP             | Plain text    | ❌ None            | ❌ No        | Standardized persistent connections         |
| HTTP/2   | TCP             | Binary frames | ✅ HPACK           | ✅ Yes       | Framing + multiplexing over TCP             |
| HTTP/3   | QUIC (over UDP) | Binary frames | ✅ QPACK           | ✅ Yes       | Fixes TCP's HOL blocking at the transport layer |

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
