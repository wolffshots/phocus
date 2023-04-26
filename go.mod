module wolffshots/phocus

go 1.19

require (
	github.com/gin-gonic/gin v1.8.1
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.1
	github.com/wolffshots/phocus_messages v0.0.0-20230426105100-58bb669225b5
	github.com/wolffshots/phocus_mqtt v0.0.0-20230426105149-7f39b6f21cd2
	github.com/wolffshots/phocus_sensors v0.0.0-20230426105352-8c5b8ada4c22
	github.com/wolffshots/phocus_serial v0.0.0-20230426105436-e8a7e9fff872
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sigurn/crc16 v0.0.0-20211026045750-20ab5afb07e3 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/wolffshots/ha_types v0.0.0-20230426105513-f6f90d3ea64f // indirect
	github.com/wolffshots/phocus_crc v0.0.0-20230426105240-c3f33a1eb597 // indirect
	go.bug.st/serial v1.4.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// // Helpful to override packages either for development or
// // to use drop in replacements for packages such as the messages
// // eg. replace github.com/wolffshots/phocus_messages => github.com/your_name/some_other_messages
//
//replace github.com/wolffshots/phocus_messages => ../phocus_messages
//
//replace github.com/wolffshots/phocus_mqtt => ../phocus_mqtt
//
//replace github.com/wolffshots/phocus_sensors => ../phocus_sensors
//
//replace github.com/wolffshots/phocus_serial => ../phocus_serial
