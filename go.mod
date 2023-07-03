module github.com/wolffshots/phocus/v2

go 1.19

require (
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.3.0
	github.com/sigurn/crc16 v0.0.0-20211026045750-20ab5afb07e3
	github.com/stretchr/testify v1.8.3
	github.com/wolffshots/ha_types v1.1.1
	github.com/wolffshots/phocus_messages v1.1.1
	go.bug.st/serial v1.5.0
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/creack/goselect v0.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/wolffshots/phocus_crc v0.0.0-20230426105240-c3f33a1eb597 // indirect
	github.com/wolffshots/phocus_mqtt v0.0.0-20230426105149-7f39b6f21cd2 // indirect
	github.com/wolffshots/phocus_serial v0.0.0-20230426105436-e8a7e9fff872 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// // Helpful to override packages either for development or
// // to use drop in replacements for packages such as the messages
// // eg. replace github.com/wolffshots/phocus_messages => github.com/your_name/some_other_messages
//
// replace github.com/wolffshots/phocus_messages => ./messages

//
// replace github.com/wolffshots/phocus_mqtt => ../phocus_mqtt
//
// replace github.com/wolffshots/phocus_sensors => ../phocus_sensors
//
// replace github.com/wolffshots/phocus_serial => ../phocus_serial
//
// replace github.com/wolffshots/phocus_api => ../phocus_api
