Take the output of `go mod graph` and print a nicely indented graph:

    $ go mod graph | head -n10
    zgo.at/goatcounter github.com/PuerkitoBio/goquery@v1.5.1
    zgo.at/goatcounter github.com/arp242/geoip2-golang@v1.4.0
    zgo.at/goatcounter github.com/go-chi/chi@v4.1.1+incompatible
    zgo.at/goatcounter github.com/jmoiron/sqlx@v1.2.0
    zgo.at/goatcounter github.com/lib/pq@v1.5.2
    zgo.at/goatcounter github.com/mattn/go-sqlite3@v2.0.3+incompatible
    zgo.at/goatcounter github.com/monoculum/formam@v0.0.0-20200527175922-6f3cce7a46cf
    zgo.at/goatcounter github.com/teamwork/reload@v1.3.2
    zgo.at/goatcounter golang.org/x/crypto@v0.0.0-20200510223506-06a226fb4e37
    zgo.at/goatcounter golang.org/x/sync@v0.0.0-20200317015054-43a5402ce75a

    $ go mod graph | gomodgraph | head -n10
    zgo.at/goatcounter
            github.com/PuerkitoBio/goquery
                    github.com/andybalholm/cascadia
                            golang.org/x/net
                                    golang.org/x/crypto
                    golang.org/x/net
                            golang.org/x/crypto
            github.com/arp242/geoip2-golang
                    github.com/arp242/maxminddb-golang
            github.com/jmoiron/sqlx
