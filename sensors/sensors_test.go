package phocus_sensors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wolffshots/ha_types/device_classes"
	"github.com/wolffshots/ha_types/state_classes"
	"github.com/wolffshots/ha_types/units"
)

func TestRegister(t *testing.T) {
	sensor := Sensor{
		SensorTopic:   "homeassistant/sensor/phocus/qid_serial/config",
		UniqueId:      "phocus_qid_serial",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QID Serial",
		ValueTemplate: "{{ value_json.SerialNumber }}",
		StateTopic:    "phocus/stats/qid",
		Icon:          "mdi:update",
	}

	sensorDefinition := Format(sensor, "v0.0.0")

	assert.Equal(t, "{\"unique_id\":\"phocus_qid_serial\",\"name\":\"QID Serial\",\"state_topic\":\"phocus/stats/qid\",\"icon\":\"mdi:update\",\"device\":{\"name\":\"phocus\",\"identifiers\":[\"phocus\"],\"model\":\"phocus\",\"manufacturer\":\"phocus\",\"sw_version\":\"v0.0.0\"},\"force_update\":false, \"value_template\":\"{{ value_json.SerialNumber }}\"}", sensorDefinition)

	sensor = Sensor{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_ac_input_frequency/config",
		UniqueId:      "phocus_qpgs2_ac_input_frequency",
		Unit:          units.Frequency,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Frequency,
		Name:          "QPGS2 AC Input Frequency",
		ValueTemplate: "{{ value_json.ACInputFrequency }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:sine-wave",
	}

	sensorDefinition = Format(sensor, "v0.0.0")

	assert.Equal(t, "{\"unique_id\":\"phocus_qpgs2_ac_input_frequency\",\"name\":\"QPGS2 AC Input Frequency\",\"state_topic\":\"phocus/stats/qpgs2\",\"icon\":\"mdi:sine-wave\",\"device\":{\"name\":\"phocus\",\"identifiers\":[\"phocus\"],\"model\":\"phocus\",\"manufacturer\":\"phocus\",\"sw_version\":\"v0.0.0\"},\"force_update\":false, \"unit_of_measurement\":\"Hz\", \"state_class\":\"measurement\", \"device_class\":\"frequency\", \"value_template\":\"{{ value_json.ACInputFrequency }}\"}", sensorDefinition)

}
