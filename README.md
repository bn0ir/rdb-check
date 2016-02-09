# rdb-check
Utility to check RDB file header, version and crc64 checksum

###WARNING
Redis has redis-check-dump utility to accomplish this task, please use it instead

### Compile
```
go get github.com/bn0ir/rdb/crc64
go build
```

### Usage
```
./rdb-check /path/to/file.rdb
```
