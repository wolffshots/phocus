// Package phocus_sensors defines and sensors for phocus and
// registers them on the MQTT broker
package phocus_sensors

import (
	"fmt"
	"log"
	"time"

	"github.com/wolffshots/ha_types/device_classes"
	"github.com/wolffshots/ha_types/state_classes"
	"github.com/wolffshots/ha_types/units"
	mqtt "github.com/wolffshots/phocus/v2/mqtt"
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
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.Timestamp,
		Name:          "Start Time",
		ValueTemplate: "",
		StateTopic:    "phocus/stats/start_time",
		Icon:          "mdi:clock",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/error/config",
		UniqueId:      "phocus_last_error",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "Last Reported Error",
		ValueTemplate: "",
		StateTopic:    "phocus/stats/error",
		Icon:          "mdi:hammer-wrench",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_serial/config",
		UniqueId:      "phocus_qpgs1_serial",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS1 Serial",
		ValueTemplate: "{{ value_json.SerialNumber }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:update",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_serial/config",
		UniqueId:      "phocus_qpgs2_serial",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS2 Serial",
		ValueTemplate: "{{ value_json.SerialNumber }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:update",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_battery_voltage/config",
		UniqueId:      "phocus_qpgs1_battery_voltage",
		Unit:          units.Voltage,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Voltage,
		Name:          "QPGS1 Battery Voltage",
		ValueTemplate: "{{ value_json.BatteryVoltage }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:battery",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_battery_voltage/config",
		UniqueId:      "phocus_qpgs2_battery_voltage",
		Unit:          units.Voltage,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Voltage,
		Name:          "QPGS2 Battery Voltage",
		ValueTemplate: "{{ value_json.BatteryVoltage }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:battery",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_battery_state_of_charge/config",
		UniqueId:      "phocus_qpgs1_battery_state_of_charge",
		Unit:          units.Battery,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Battery,
		Name:          "QPGS1 Battery SoC",
		ValueTemplate: "{{ value_json.BatteryStateOfCharge }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:battery",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_battery_state_of_charge/config",
		UniqueId:      "phocus_qpgs2_battery_state_of_charge",
		Unit:          units.Battery,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Battery,
		Name:          "QPGS2 Battery SoC",
		ValueTemplate: "{{ value_json.BatteryStateOfCharge }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:battery",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_operation_mode/config",
		UniqueId:      "phocus_qpgs1_operation_mode",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS1 Operation Mode",
		ValueTemplate: "{{ value_json.OperationMode }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:meter-electric",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_operation_mode/config",
		UniqueId:      "phocus_qpgs2_operation_mode",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS2 Operation Mode",
		ValueTemplate: "{{ value_json.OperationMode }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:meter-electric",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_ac_output_active_power/config",
		UniqueId:      "phocus_qpgs1_ac_output_active_power",
		Unit:          units.Power,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Power,
		Name:          "QPGS1 AC Output Active Power",
		ValueTemplate: "{{ value_json.ACOutputActivePower }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_ac_output_active_power/config",
		UniqueId:      "phocus_qpgs2_ac_output_active_power",
		Unit:          units.Power,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Power,
		Name:          "QPGS2 AC Output Active Power",
		ValueTemplate: "{{ value_json.ACOutputActivePower }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_ac_output_apparent_power/config",
		UniqueId:      "phocus_qpgs1_ac_output_apparent_power",
		Unit:          units.ApparentPower,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.ApparentPower,
		Name:          "QPGS1 AC Output Apparent Power",
		ValueTemplate: "{{ value_json.ACOutputApparentPower }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_ac_output_apparent_power/config",
		UniqueId:      "phocus_qpgs2_ac_output_apparent_power",
		Unit:          units.ApparentPower,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.ApparentPower,
		Name:          "QPGS2 AC Output Apparent Power",
		ValueTemplate: "{{ value_json.ACOutputApparentPower }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_pv_input_voltage/config",
		UniqueId:      "phocus_qpgs1_pv_input_voltage",
		Unit:          units.Voltage,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Voltage,
		Name:          "QPGS1 PV Input Voltage",
		ValueTemplate: "{{ value_json.PVInputVoltage }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_pv_input_voltage/config",
		UniqueId:      "phocus_qpgs2_pv_input_voltage",
		Unit:          units.Voltage,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Voltage,
		Name:          "QPGS2 PV Input Voltage",
		ValueTemplate: "{{ value_json.PVInputVoltage }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_pv_input_current/config",
		UniqueId:      "phocus_qpgs1_pv_input_current",
		Unit:          units.Current,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Current,
		Name:          "QPGS1 PV Input Current",
		ValueTemplate: "{{ value_json.PVInputCurrent }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:current",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_pv_input_current/config",
		UniqueId:      "phocus_qpgs2_pv_input_current",
		Unit:          units.Current,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Current,
		Name:          "QPGS2 PV Input Current",
		ValueTemplate: "{{ value_json.PVInputCurrent }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:current",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_battery_discharge_current/config",
		UniqueId:      "phocus_qpgs1_battery_discharge_current",
		Unit:          units.Current,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Current,
		Name:          "QPGS1 Battery Discharge Current",
		ValueTemplate: "{{ value_json.BatteryDischargeCurrent }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:current",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_battery_discharge_current/config",
		UniqueId:      "phocus_qpgs2_battery_discharge_current",
		Unit:          units.Current,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Current,
		Name:          "QPGS2 Battery Discharge Current",
		ValueTemplate: "{{ value_json.BatteryDischargeCurrent }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:current",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_battery_charge_current/config",
		UniqueId:      "phocus_qpgs1_battery_charge_current",
		Unit:          units.Current,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Current,
		Name:          "QPGS1 Battery Charge Current",
		ValueTemplate: "{{ value_json.BatteryChargingCurrent }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:current",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_battery_charge_current/config",
		UniqueId:      "phocus_qpgs2_battery_charge_current",
		Unit:          units.Current,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Current,
		Name:          "QPGS2 Battery Charge Current",
		ValueTemplate: "{{ value_json.BatteryChargingCurrent }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:current",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_ac_input_mode/config",
		UniqueId:      "phocus_qpgs1_ac_input_mode",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS1 AC Input Mode",
		ValueTemplate: "{{ value_json.InverterStatus.ACInput }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:meter-electric",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_ac_input_mode/config",
		UniqueId:      "phocus_qpgs2_ac_input_mode",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS2 AC Input Mode",
		ValueTemplate: "{{ value_json.InverterStatus.ACInput }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:meter-electric",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_total_ac_output_active_power/config",
		UniqueId:      "phocus_qpgs1_total_ac_output_active_power",
		Unit:          units.Power,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Power,
		Name:          "QPGS1 Total AC Output Active Power",
		ValueTemplate: "{{ value_json.TotalACOutputActivePower }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_total_ac_output_active_power/config",
		UniqueId:      "phocus_qpgs2_total_ac_output_active_power",
		Unit:          units.Power,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Power,
		Name:          "QPGS2 Total AC Output Active Power",
		ValueTemplate: "{{ value_json.TotalACOutputActivePower }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_total_ac_output_apparent_power/config",
		UniqueId:      "phocus_qpgs1_total_ac_output_apparent_power",
		Unit:          units.ApparentPower,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.ApparentPower,
		Name:          "QPGS1 Total AC Output Apparent Power",
		ValueTemplate: "{{ value_json.TotalACOutputApparentPower }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_total_ac_output_apparent_power/config",
		UniqueId:      "phocus_qpgs2_total_ac_output_apparent_power",
		Unit:          units.ApparentPower,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.ApparentPower,
		Name:          "QPGS2 Total AC Output Apparent Power",
		ValueTemplate: "{{ value_json.TotalACOutputApparentPower }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_ac_input_voltage/config",
		UniqueId:      "phocus_qpgs1_ac_input_voltage",
		Unit:          units.Voltage,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Voltage,
		Name:          "QPGS1 AC Input Voltage",
		ValueTemplate: "{{ value_json.ACInputVoltage }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_ac_input_voltage/config",
		UniqueId:      "phocus_qpgs2_ac_input_voltage",
		Unit:          units.Voltage,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Voltage,
		Name:          "QPGS2 AC Input Voltage",
		ValueTemplate: "{{ value_json.ACInputVoltage }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:lightning-bolt",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_ac_input_frequency/config",
		UniqueId:      "phocus_qpgs1_ac_input_frequency",
		Unit:          units.Frequency,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Frequency,
		Name:          "QPGS1 AC Input Frequency",
		ValueTemplate: "{{ value_json.ACInputFrequency }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:sine-wave",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_ac_input_frequency/config",
		UniqueId:      "phocus_qpgs2_ac_input_frequency",
		Unit:          units.Frequency,
		StateClass:    state_classes.Measurement,
		DeviceClass:   device_classes.Frequency,
		Name:          "QPGS2 AC Input Frequency",
		ValueTemplate: "{{ value_json.ACInputFrequency }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:sine-wave",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs1_checksum/config",
		UniqueId:      "phocus_qpgs1_checksum",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS1 Checksum",
		ValueTemplate: "{{ value_json.Checksum }}",
		StateTopic:    "phocus/stats/qpgs1",
		Icon:          "mdi:check",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qpgs2_checksum/config",
		UniqueId:      "phocus_qpgs2_checksum",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QPGS2 Checksum",
		ValueTemplate: "{{ value_json.Checksum }}",
		StateTopic:    "phocus/stats/qpgs2",
		Icon:          "mdi:check",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/generic_response/config",
		UniqueId:      "phocus_generic_response",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "Generic Response",
		ValueTemplate: "{{ value_json.Result }}",
		StateTopic:    "phocus/stats/generic",
		Icon:          "fab:readme",
	},
	{
		SensorTopic:   "homeassistant/sensor/phocus/qid_serial/config",
		UniqueId:      "phocus_qid_serial",
		Unit:          units.None,
		StateClass:    state_classes.None,
		DeviceClass:   device_classes.None,
		Name:          "QID Serial",
		ValueTemplate: "{{ value_json.SerialNumber }}",
		StateTopic:    "phocus/stats/qid",
		Icon:          "mdi:update",
	},
}

func Format(sensor Sensor, version string) string {
	log.Printf("Registering %s\n", sensor.Name)

	sensorDefinition := fmt.Sprintf(
		"{\""+
			"unique_id\":\"%s\",\""+
			"name\":\"%s\",\""+
			"state_topic\":\"%s\",\""+
			"icon\":\"%s\",\""+
			"device\":{\"name\":\"phocus\",\""+
			"identifiers\":[\"phocus\"],\""+
			"model\":\"phocus\",\""+
			"manufacturer\":\"phocus\",\""+
			"sw_version\":\"%s\"},\""+
			"force_update\":false",
		sensor.UniqueId,
		sensor.Name,
		sensor.StateTopic,
		sensor.Icon,
		version,
	)
	if sensor.Unit != "" {
		sensorDefinition += fmt.Sprintf(", \"unit_of_measurement\":\"%s\"", sensor.Unit)
	}
	if sensor.StateClass != "" {
		sensorDefinition += fmt.Sprintf(", \"state_class\":\"%s\"", sensor.StateClass)
	}
	if sensor.DeviceClass != "" {
		sensorDefinition += fmt.Sprintf(", \"device_class\":\"%s\"", sensor.DeviceClass)
	}
	if sensor.ValueTemplate != "" {
		sensorDefinition += fmt.Sprintf(", \"value_template\":\"%s\"", sensor.ValueTemplate)
	}
	sensorDefinition += "}"
	return sensorDefinition
}

// Register adds some sensors to Home Assistant MQTT
// version is the current version of the system, added in 1.1.1
func Register(version string) error {
	log.Println("Registering sensors")
	for _, sensor := range sensors {

		sensorDefinition := Format(sensor, version)

		err := mqtt.Send(sensor.SensorTopic, 0, true, sensorDefinition, 10)
		if err != nil {
			log.Printf("Failed to send initial setup stats to MQTT with err: %v", err)
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
