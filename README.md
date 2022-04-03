Take the output of `go mod graph` and print a nicely indented graph:

    % go mod graph | head -n5
    zgo.at/goatcounter/v2 code.soquee.net/otp@v0.0.1
    zgo.at/goatcounter/v2 github.com/BurntSushi/toml@v1.0.0
    zgo.at/goatcounter/v2 github.com/PuerkitoBio/goquery@v1.8.0
    zgo.at/goatcounter/v2 github.com/andybalholm/cascadia@v1.3.1
    zgo.at/goatcounter/v2 github.com/bmatcuk/doublestar/v3@v3.0.0

    % go mod graph | gomodgraph | head -n15
    zgo.at/goatcounter/v2
            code.soquee.net/otp
            github.com/BurntSushi/toml
            github.com/PuerkitoBio/goquery
            │       github.com/andybalholm/cascadia
            github.com/bmatcuk/doublestar/v3
            github.com/boombuler/barcode
            github.com/go-chi/chi/v5
            github.com/google/uuid
            github.com/gorilla/websocket
            github.com/mattn/go-sqlite3
            github.com/monoculum/formam
            github.com/oschwald/geoip2-golang
            │       github.com/oschwald/maxminddb-golang
            │       │       github.com/stretchr/testify


Use `-v` to print the version too, and `-d` to set the maximum depth:

Just makes it easier to see "why is this package included?"
