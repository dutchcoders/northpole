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

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

**Uvis Grinfelds**

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.
Code released under [the MIT license](LICENSE).
