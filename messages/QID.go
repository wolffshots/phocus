package phocus_messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	phocus_crc "github.com/wolffshots/phocus/v2/crc"
	phocus_mqtt "github.com/wolffshots/phocus/v2/mqtt"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

type QIDResponse struct {
	SerialNumber string
}

func SendQID(port phocus_serial.Port, payload interface{}) (int, error) {
	written, err := port.Write(phocus_crc.Encode("QID"))
	if err != nil {
		return -1, err
	} else {
		return written, nil
	}
}

func ReceiveQID(port phocus_serial.Port, timeout time.Duration) (string, error) {
	// read from port
	response, err := port.Read(timeout)
	log.Println(response)
	// verify
	if err != nil || response == "" {
		log.Printf("Failed to read from serial with: %v\n", err)
		return "", err
	} else {
		if phocus_crc.Verify(response) {
			serialNumber := strings.Trim(response[:len(response)-3], "(")
			// return
			// TODO add check for length
			log.Printf("Serial number queried: %s\n", serialNumber)
			return serialNumber, nil
		} else {
			actual := response[len(response)-3 : len(response)-1]
			remainder := response[:len(response)-3]
			wanted := phocus_crc.Checksum(remainder)
			message := fmt.Sprintf("invalid response from QID: CRC should have been %x but was %x", wanted, actual)
			log.Println(message)
			return "", errors.New(message)
		}
	}
}

func InterpretQID(response string) (*QIDResponse, error) {
	return &QIDResponse{
		SerialNumber: response,
	}, nil
}

func PublishQID(response *QIDResponse) error {
	jsonQIDResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to parse response to json with :%v", err)
	}
	err = phocus_mqtt.Send("phocus/stats/qid", 0, false, string(jsonQIDResponse), 10)
	if err != nil {
		log.Fatalf("MQTT send of %s failed with: %v\ntype of thing sent was: %T", "QID", err, jsonQIDResponse)
	}
	log.Printf("Sent to MQTT:\n%s\n", jsonQIDResponse)
	return err
}
