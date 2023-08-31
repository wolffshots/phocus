// Define the type for the JSON data
interface InverterData {
    InverterNumber: number;
    OtherUnits: boolean;
    SerialNumber: string;
    OperationMode: string;
    FaultCode: string;
    ACInputVoltage: string;
    ACInputFrequency: string;
    ACOutputVoltage: string;
    ACOutputFrequency: string;
    ACOutputApparentPower: string;
    ACOutputActivePower: string;
    PercentageOfNominalOutputPower: string;
    BatteryVoltage: string;
    BatteryChargingCurrent: string;
    BatteryStateOfCharge: string;
    PVInputVoltage: string;
    TotalChargingCurrent: string;
    TotalACOutputApparentPower: string;
    TotalACOutputActivePower: string;
    TotalPercentageOfNominalOutputPower: string;
    InverterStatus: {
        MPPT: string;
        ACCharging: string;
        SolarCharging: string;
        BatteryStatus: string;
        ACInput: string;
        ACOutput: string;
        Reserved: string;
    };
    ACOutputMode: string;
    BatteryChargerSourcePriority: string;
    MaxChargingCurrentSet: string;
    MaxChargingCurrentPossible: string;
    MaxACChargingCurrentSet: string;
    PVInputCurrent: string;
    BatteryDischargeCurrent: string;
    Checksum: string;
    [key: string]: any;
};

export default InverterData;