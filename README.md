Northpole
=========

Northpole: an experimental and work in progress sync server for Santa.

[github.com/google/santa]

## Configuration
Update the **/var/db/santa/config.plist** file with the value for SyncBaseURL.

```
<key>SyncBaseURL</key>
<string>http://127.0.0.1:8080/</string>
```

## Run server
```
go run *.go
```

## Test
Now you can run Santa sync to connect with Nortpole.

```
sudo santa sync
```

## Contributions

Contributions are welcome.

## Generate certificates
http://golang.org/src/crypto/tls/generate_cert.go?m=text

go run contrib/generate_cert.go -host="127.0.0.1" -ca=false
openssl x509 -in /var/db/santa/cert.pem  -text

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.
Code released under [the MIT license](LICENSE).
