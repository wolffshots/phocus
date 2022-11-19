package sensors

import (
	"fmt"
	"log"
	"time"
	"wolffshots/phocus/mqtt"
	"wolffshots/phocus/types/device_classes"
	"wolffshots/phocus/types/state_classes"
	"wolffshots/phocus/types/units"
)

// Sensor is the shape of the sensor for the MQTT Home Assistant integration
type Sensor struct {
	SensorTopic   string                     // "homeassistant/sensor/phocus/start_time/config" must end in /config
	UniqueId      string                     // "unique_id": "phocus_qpgs1_ac_output_apparent_power",
	Unit          units.Unit                 // "unit_of_measurement": "VA",
	StateClass    state_classes.StateClass   // "state_class": "measurement",
	DeviceClass   device_classes.DeviceClass // "device_class": "apparent_power",
	Name          string                     // "name": "QPGS1 AC Output Apparent Power",
	ValueTemplate string                     // "value_template": "{{ value_json.ACOutputApparentPower }}",
	StateTopic    string                     // "state_topic": "phocus/stats/qpgs1",
	Icon          string                     // "icon": "mdi:battery",
}

var sensors = []Sensor{
	{
		SensorTopic:   "homeassistant/sensor/phocus/start_time/config",
		UniqueId:      "phocus_start_time",
		Unit:          units.None,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Timestamp,
		Name:          "Start Time",
		ValueTemplate: "{{ value_json }}",
		StateTopic:    "phocus/stats/start_time",
		Icon:          "mdi:clock",
	}, //err := mqtt.Send("homeassistant/sensor/phocus/start_time/config", 0, true, `{"unique_id":"phocus_start_time","name":"phocus - Start Time","state_topic":"phocus/stats/start_time","icon":"mdi:hammer-wrench","device":{"name":"phocus","identifiers":["phocus"],"model":"phocus","manufacturer":"phocus","sw_version":"1.1.0"},"force_update":false}`, 10)

}

// Register adds some sensors to Home Assistant MQTT
func Register() {
	log.Println("Registering sensors")
	for _, input := range sensors {
		log.Printf("Registering %s\n", input.Name)
		err := mqtt.Send(input.SensorTopic, 0, true, fmt.Sprintf(
			"{\""+
				"unique_id\":\"%s\",\""+
				"name\":\"%s\",\""+
				"state_topic\":\"%s\",\""+
				"icon\":\"%s\",\""+
				"value_template\":\"%s\",\""+
				"unit\":\"%s\",\""+
				"state_class\":\"%s\",\""+
				"device_class\":\"%s\",\""+
				"device\":{\"name\":\"phocus\",\""+
				"identifiers\":[\"phocus\"],\""+
				"model\":\"phocus\",\""+
				"manufacturer\":\"phocus\",\""+
				"sw_version\":\"1.1.0\"},\""+
				"force_update\":false"+
				"}",
			input.UniqueId,
			input.Name,
			input.StateTopic,
			input.Icon,
			input.ValueTemplate,
			input.Unit,
			input.StateClass,
			input.DeviceClass,
		), 10)
		if err != nil {
			log.Fatalf("Failed to send initial setup stats to MQTT with err: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

}
