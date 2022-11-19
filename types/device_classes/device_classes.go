// Package device_classes
package device_classes

// DeviceClass refers to the type of device from Home Assistant's perspective
//
// [Source]: https://www.home-assistant.io/integrations/sensor#device-class
type DeviceClass string

const (
	ApparentPower            = "apparent_power"             // ApparentPower in VA.
	AirQualityIndex          = "aqi"                        // AirQualityIndex
	Battery                  = "battery"                    // Percentage of Battery that is left
	CarbonDioxide            = "carbon_dioxide"             // CarbonDioxide in CO2 (Smoke)
	CarbonMonoxide           = "carbon_monoxide"            // CarbonMonoxide in CO (Gas CNG/LPG)
	Current                  = "current"                    // Current in A
	Date                     = "date"                       // Date string (ISO 8601)
	Distance                 = "distance"                   // Generic Distance in km, m, cm, mm, mi, yd, or in
	Duration                 = "duration"                   // Duration in days, hours, minutes or seconds
	Energy                   = "energy"                     // Energy in Wh, kWh or MWh
	Frequency                = "frequency"                  // Frequency in Hz, kHz, MHz or GHz
	Gas                      = "gas"                        // Gas volume in m³ or ft³
	Humidity                 = "humidity"                   // Percentage of Humidity in the air
	Illuminance              = "illuminance"                // The current light level (Illuminance) in lx or lm
	Moisture                 = "moisture"                   // Percentage of water in a substance (Moisture)
	Monetary                 = "monetary"                   // The Monetary value
	NitrogenDioxide          = "nitrogen_dioxide"           // Concentration of NitrogenDioxide in µg/m³
	NitrogenMonoxide         = "nitrogen_monoxide"          // Concentration of NitrogenMonoxide in µg/m³
	NitrousOxide             = "nitrous_oxide"              // Concentration of NitrousOxide in µg/m³
	Ozone                    = "ozone"                      // Concentration of Ozone in µg/m³
	ParticulateMatter1       = "pm1"                        // Concentration of ParticulateMatter1 less than 1 micrometer in µg/m³
	ParticulateMatter10      = "pm10"                       // Concentration of ParticulateMatter10 less than 10 micrometers in µg/m³
	ParticulateMatter25      = "pm25"                       // Concentration of ParticulateMatter25 less than 2.5 micrometers in µg/m³
	PowerFactor              = "power_factor"               // PowerFactor in %
	Power                    = "power"                      // Power in W or kW
	PrecipitationIntensity   = "precipitation_intensity"    // PrecipitationIntensity in in/d, in/h, mm/d, or mm/h
	Pressure                 = "pressure"                   // Pressure in Pa, kPa, hPa, bar, cbar, mbar, mmHg, inHg, or psi
	ReactivePower            = "reactive_power"             //  ReactivePower in var
	SignalStrength           = "signal_strength"            //  SignalStrength in dB or dBm
	Speed                    = "speed"                      //  Generic Speed in ft/s, in/d, in/h, km/h, kn, m/s, mph, or mm/d
	SulphurDioxide           = "sulphur_dioxide"            //  Concentration of SulphurDioxide in µg/m³
	Temperature              = "temperature"                //  Temperature in °C or °F
	Timestamp                = "timestamp"                  //  Datetime object or Timestamp string (ISO 8601)
	VolatileOrganicCompounds = "volatile_organic_compounds" //  Concentration of VolatileOrganicCompounds in µg/m³
	Voltage                  = "voltage"                    //  Voltage in V
	Volume                   = "volume"                     //  Generic Volume in L, mL, gal, fl. oz., m³, or ft³
	Water                    = "water"                      //  Water consumption in L, gal, m³, or ft³
	Weight                   = "weight"                     //  Generic mass in kg, g, mg, µg, oz, or lb (Weight)
	WindSpeed                = "wind_speed"                 //  WindSpeed in ft/s, km/h, kn, m/s, or mph
	None                     = ""                           //  No class
)
