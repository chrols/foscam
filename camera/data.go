package camera

import "encoding/xml"

type ReturnCode uint8

const (
	Success ReturnCode = iota
	CGIRequestStringFormatError
	UsernameOrPasswordError
	AccessDenied
	CGIExecutionError
	Timeout
	Reserved1
	UnknownError
	Reserved2
)

type ReturnValue struct {
	XMLName xml.Name   `xml:"CGI_Result"`
	Result  ReturnCode `xml:"result"`
}

type DeviceInformation struct {
	XMLName         xml.Name   `xml:"CGI_Result"`
	Result          ReturnCode `xml:"result"`
	ProductName     string     `xml:"productName"`
	SerialNumber    string     `xml:"serialNo"`
	DeviceName      string     `xml:"devName"`
	MAC             string     `xml:"mac"`
	Year            string     `xml:"year"`
	Month           string     `xml:"mon"`
	Day             string     `xml:"day"`
	Hour            uint       `xml:"hour"`
	Min             uint       `xml:"min"`
	Sec             uint       `xml:"sec"`
	TimeZone        int        `xml:"timeZone"`
	FirmwareVersion string     `xml:"firmwareVer"`
	HardwareVersion string     `xml:"hardwareVer"`
}

type MotionDetectConfig struct {
	XMLName           xml.Name   `xml:"CGI_Result"`
	Result            ReturnCode `xml:"result"`
	IsEnable          bool       `xml:"isEnable"`
	Linkage           uint8      `xml:"linkage"`
	SnapInterval      int        `xml:"snapInterval"`
	Sensitivity       int        `xml:"sensitivity"`
	Triggerinterval   int        `xml:"triggerInterval"`
	ScheduleMonday    uint64     `xml:"schedule0"`
	ScheduleTuesday   uint64     `xml:"schedule1"`
	ScheduleWednesday uint64     `xml:"schedule2"`
	ScheduleThursday  uint64     `xml:"schedule3"`
	ScheduleFriday    uint64     `xml:"schedule4"`
	ScheduleSaturday  uint64     `xml:"schedule5"`
	ScheduleSunday    uint64     `xml:"schedule6"`
	Area0             uint16     `xml:"area0"`
	Area1             uint16     `xml:"area1"`
	Area2             uint16     `xml:"area2"`
	Area3             uint16     `xml:"area3"`
	Area4             uint16     `xml:"area4"`
	Area5             uint16     `xml:"area5"`
	Area6             uint16     `xml:"area6"`
	Area7             uint16     `xml:"area7"`
	Area8             uint16     `xml:"area8"`
	Area9             uint16     `xml:"area9"`
}
