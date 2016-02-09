package main

import (
	"bytes"
	"encoding/binary"
	"github.com/bn0ir/rdb/crc64"
	"hash"
	"log"
	"os"
	"strconv"
)

func main() {
	// check dump file name is specified
	if len(os.Args) < 2 {
		log.Printf("Please specify path to dump file: checkrdb /path/to/file.rdb")
		os.Exit(1)
	}

	filepath := os.Args[1]

	// check dump file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Printf("RDB file %v does not exist", filepath)
		os.Exit(1)
	}

	// open rdb file
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Can't open RDB file %v", filepath)
		os.Exit(1)
	}
	defer file.Close()

	// check rdb file header
	header := make([]byte, 9)
	_, err = file.ReadAt(header, 0)

	if err != nil {
		log.Printf("Can't read RDB file %v header", filepath)
		os.Exit(1)
	}

	if !bytes.Equal(header[:5], []byte("REDIS")) {
		log.Printf("Invalid RDB file format")
		os.Exit(1)
	}

	rdbversion, _ := strconv.ParseInt(string(header[5:]), 10, 64)
	if rdbversion < 1 || rdbversion > 6 {
		log.Printf("Invalid RDB file version number %d", rdbversion)
	}

	log.Printf("RDB file version: %v", rdbversion)

	// get crc64 from rdb file
	footer := make([]byte, 8)
	stat, err := os.Stat(filepath)
	startfooter := stat.Size() - 8
	_, err = file.ReadAt(footer, startfooter)

	if err != nil {
		log.Printf("Can't read RDB file %v footer", filepath)
		os.Exit(1)
	}

	log.Printf("CRC64 from RDB file: %v", binary.LittleEndian.Uint64(footer))

	// calculate crc64 of file
	bufLen := int64(10 * 1024 * 1024)
	buf := make([]byte, bufLen)
	var position int64
	position = 0

	crc64sum := hash.Hash64(crc64.New())

	for position < startfooter {
		if (position + bufLen) > startfooter {
			bufLen = startfooter - position
			buf = make([]byte, bufLen)
		}
		_, err = file.ReadAt(buf, position)
		crc64sum.Write(buf)
		position = position + bufLen
	}

	log.Printf("CRC64 calculated: %v", crc64sum.Sum64())

	if crc64sum.Sum64() != binary.LittleEndian.Uint64(footer) {
		log.Printf("Wrong RDB file checksum: %v != %v", binary.LittleEndian.Uint64(footer), crc64sum.Sum64())
		os.Exit(1)
	}

	log.Printf("RDB file check success")
}
