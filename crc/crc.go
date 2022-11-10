package crc

import (
	"github.com/sigurn/crc16"
    "log"
)

func Checksum(input string) (uint16, error) {
	table := crc16.MakeTable(crc16.CRC16_XMODEM)
	result := crc16.Checksum([]byte(input), table)
	return result, nil
}

func Encode(input string) (string, error){
    checksum, err := Checksum(input)
	log.Printf("crc for encode: %x\n", checksum)
    if err != nil {
		log.Fatal(err)
	}
    result := input + string([]byte{byte((checksum >> 8) & 0xff), byte(checksum & 0xff)}) + "\r"
    
    log.Printf("result to be written: %s\n", result)
    
    return result, err
}

func Verify(input string) (bool, error){
    /* TODO 
    - cut off CRC  
    - put remainder into Encode
    - compare input to result
    - return true if input (with original CRC) matches the Encode-d result 
    */
    return true, nil
}
