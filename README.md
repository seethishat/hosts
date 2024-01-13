# A collection of hosts files to block certain domains

## Source
```bash
curl https://rpz-api.iso.vt.edu/bl/domains | jq | grep name | awk '{print $2}' | tr --delete ",\"" | sort -u | awk '{print "0.0.0.0 " $0}' > hosts
```
