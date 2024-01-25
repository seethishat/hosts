# A collection of hosts files to block certain domains

## Tor and VT Sources

```bash
$ curl https://rpz-api.iso.vt.edu/bl/domains | jq | grep name | awk '{print $2}' | tr --delete ",\"" | sort -u | awk '{print "0.0.0.0 " $0}' > vt.rpz
$ curl https://www.seethishat.com/static/tor-hosts.txt | sort -u | awk '{print "0.0.0.0 " $0}' > sth.tor
```

## Misc

* [hosts](hosts) is compressed and should work on Windows, Linux and Mac systems.
* This hosts file intentionally blocks Google, YouTube, TikTok and many other popular domains.
