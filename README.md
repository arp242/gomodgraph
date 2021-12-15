Take the output of `go mod graph` and print a nicely indented graph:

    % go mod graph | head -n5
    zgo.at/goatcounter/v2 code.soquee.net/otp@v0.0.1
    zgo.at/goatcounter/v2 github.com/BurntSushi/toml@v0.4.1
    zgo.at/goatcounter/v2 github.com/PuerkitoBio/goquery@v1.8.0
    zgo.at/goatcounter/v2 github.com/andybalholm/cascadia@v1.3.1
    zgo.at/goatcounter/v2 github.com/bmatcuk/doublestar/v3@v3.0.0

    % go mod graph | gomodgraph | head -n10
    zgo.at/goatcounter/v2
            code.soquee.net/otp
            github.com/BurntSushi/toml
            github.com/PuerkitoBio/goquery
                    github.com/andybalholm/cascadia
                            golang.org/x/net
                    golang.org/x/net
            github.com/andybalholm/cascadia
                    golang.org/x/net
            github.com/bmatcuk/doublestar/v3


Use `-v` to print the version too, and `-d` to set the maximum depth:

    % go mod graph | gomodgraph -v -d 2 | head -n10
    zgo.at/goatcounter/v2
            code.soquee.net/otp v0.0.1
            github.com/BurntSushi/toml v0.4.1
            github.com/PuerkitoBio/goquery v1.8.0
                    golang.org/x/net v0.0.0-20210916014120-12bc252f5db8
            github.com/andybalholm/cascadia v1.3.1
                    golang.org/x/net v0.0.0-20210916014120-12bc252f5db8
            github.com/bmatcuk/doublestar/v3 v3.0.0
            github.com/boombuler/barcode v1.0.1
            github.com/fsnotify/fsnotify v1.4.9

Just makes it easier to see "why is this package included?"
