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

## HTTP/1.1 (Text) - Used here

Because HTTP uses TCP, if the HTTP request or response is too big to fit into a single TCP packet it can be broken up into many packets and reconstructed in the correct order on the other side.

TCP guarantees that the data is in order and complete.

At the heart of HTTP is the `HTTP-message`: the format that the text in an HTTP request or response must use. From [RFC 9112 Section 2.1](https://datatracker.ietf.org/doc/html/rfc9112#name-message-format):

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

At the moment, there are several key RFCs for HTTP/1.1:

- [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231) – An active and widely referenced RFC.
- [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) – Easier to read than RFC 7231, relies on understanding from RFC 9110.
- [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110) – Covers HTTP "semantics."

### Example GET - Request

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

```text
GET /goodies HTTP/1.1       # start-line CRLF
Host: localhost:42069       # *( field-line CRLF )
User-Agent: curl/7.81.0     # *( field-line CRLF )
Accept: */*                 # *( field-line CRLF )
                            # CRLF
                            # [ message-body ] (empty)
```

### Example POST - Request

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

if it's a request, then the start-line is called the `request-line` and has a specific format.

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

- The RFC uses the term `field-line`, but it's basically the same thing. From [5. Field Syntax](https://datatracker.ietf.org/doc/html/rfc9112#name-field-syntax):

"_Each field line consists of a case-insensitive field name followed by a colon ("\:"), optional leading whitespace, the field line value, and optional trailing whitespace._"

```text
field-line   = field-name ":" OWS field-value OWS
```

- There can be an unlimited amount of whitespace before and after the field-value (Header value).
- There must be no spaces betwixt the colon and the field-name. Example:

```text
// valid
'Host: localhost:42069'
' Host: localhost:42069 '

// not valid
Host : localhost:42069
```

- Keys are case insensitive
- `field-name` has an implicit definition of a token as defined by [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110).
- `field-name` must contain only:
  - Uppercase letters: A-Z
  - Lowercase letters: a-z
  - Digits: 0-9
  - Special characters: !, #, $, %, &, ', \*, +, -, ., ^, \_, `, |, ~
  - At least a length of 1.
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
This word, or the adjective "RECOMMENDED", mean that there may exist valid reasons in particular circumstances to ignore a particular item, but the full implications must be understood and carefully weighed before choosing a different course.
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

- About the reason-phrase, from Section 4,

"_A client SHOULD ignore the reason-phrase content because it is not a reliable channel for information (it might be translated for a given locale, overwritten by intermediaries, or discarded when the message is forwarded via other versions of HTTP). A server MUST send the space that separates the status-code from the reason-phrase even when the reason-phrase is absent (i.e., the status-line would end with the space)._"

So, while reason phrases are typically included (and match one to one with the status code), they are not required and should be ignored by clients.

### HTTP/2/3 (Binary)

Frames replace raw text. Example frame types:

- **HEADERS** (compressed HTTP headers)

- **DATA** (response body chunks)

- **SETTINGS** (configuration)

**HPACK/QPACK**: Encodes headers as binary (e.g., **:method:** GET → a tiny integer)
