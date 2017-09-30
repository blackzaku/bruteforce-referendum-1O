package o1st

import (
	"os"
	"fmt"
	"bufio"
	"encoding/hex"
	"encoding/binary"
)

var data [256*256][128][206]byte

func toHex(num uint32) string {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, num)
	return hex.EncodeToString(bs)[0:2]
}

func ReadData(baseFolder string)  {
	var filename string
	var d uint32 = 0
	var f uint32 = 0
	l := 0
	for d < 256 {
		f = 0
		for f < 256 {
			// Read file and write to the buffer
			filename = baseFolder + "/" + toHex(d) + "_"+toHex(f)+".db"
			file, err := os.Open(filename)
			if err != nil {
				fmt.Println("Could not read", filename, err)
			} else {
				scanner := bufio.NewScanner(file)
				l = 0
				for scanner.Scan() {
					line := scanner.Bytes()
					hex.Decode(data[d*256+f][l][:], line)
					l += 1
				}
			}
			file.Close()
			f += 1
		}
		d += 1
	}
}