# httpx

[![Licence](https://img.shields.io/github/license/krostar/httpx.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/httpx.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/httpx/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/httpx)
[![Code quality](https://img.shields.io/codacy/grade/128a195d311c47db872c0dca0555ad5c/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/httpx/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/128a195d311c47db872c0dca0555ad5c.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/httpx/dashboard)

Useful set of functions and middlewares to use on top of net/http.

## Motivations

On every project that requires a http server, I find myself copy/pasting the same
boiloerplate over and over again.
While I agree with [Rob Pike](https://go-proverbs.github.io/) saying that a little
copying is better than a little dependency, I think it's time to release and use
only one copy of these helpers / middleware functions.

## Examples

The code below ...

-   creates a http server with sane timeouts
-   sets a modern tls config [described by mozilla](https://wiki.mozilla.org/Security/Server_Side_TLS) for the listener
-   sets a keepalive, and configures a period for it, [see this blog on cloudflare](https://blog.cloudflare.com/exposing-go-on-the-internet/)
    to get why it can be important)
-   start the server and set two signals to gracefully stop it
-   is heavility tested ;)

```go
func main() {
    server, err := httpx.NewServer(initRoutes(usecases))
    if err != nil {
        panic("unable to create server")
    }
    server.ErrorLog = logger.StdLog(
        log.WithField("source", "http-error"), logger.LevelWarn,
    )

    listener, err := httpx.NewListener(":8080",
        httpx.ListenWithKeepAlive(time.Second * 45),
        httpx.ListenWithModernTLSConfig("./domain.crt", "./domain.key"),
    )
    if err != nil {
        panic("unable to create listener")
    }

    if err := httpx.StartAndStopWithSignal(server, listener, time.Second * 10, syscall.SIGINT, syscall.SIGTERM); err != nil {
        panic(err)
    }
}
```

## License

This project is under the MIT licence, please see the LICENCE file.
