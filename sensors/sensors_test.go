package phocus_sensors

import (
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	log.Printf(
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
			"sw_version\":\"%s\"},\""+
			"force_update\":false"+
			"}",
		sensors[0].UniqueId,
		sensors[0].Name,
		sensors[0].StateTopic,
		sensors[0].Icon,
		sensors[0].ValueTemplate,
		sensors[0].Unit,
		sensors[0].StateClass,
		sensors[0].DeviceClass,
		"1.1.1",
	)
	log.Println(sensors)
}
