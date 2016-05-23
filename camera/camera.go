package camera

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

type Camera struct {
	Alias    string
	URL      string
	User     string
	Password string
}

type Direction string

const (
	MoveUp          Direction = "ptzMoveUp"
	MoveDown        Direction = "ptzMoveDown"
	MoveLeft        Direction = "ptzMoveLeft"
	MoveRight       Direction = "ptzMoveRight"
	MoveTopLeft     Direction = "ptzTopLeft"
	MoveTopRight    Direction = "ptzTopRight"
	MoveBottomLeft  Direction = "ptzBottomLeft"
	MoveBottomRight Direction = "ptzBottomRight"
	StopMove        Direction = "ptzStopRun"
	Reset           Direction = "ptzReset"
)

func (camera *Camera) GetMotionDetectConfig() MotionDetectConfig {

	resp, err := http.Get("http://" + camera.URL + "/cgi-bin/CGIProxy.fcgi?usr=" + camera.User + "&pwd=" + camera.Password + "&cmd=getMotionDetectConfig")

	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	v := MotionDetectConfig{}

	err = xml.Unmarshal(b, &v)

	if err != nil {
		log.Fatal(err)
	}

	return v
}

func (camera *Camera) SetMotionDetectConfig(config MotionDetectConfig) bool {
	request := "http://" + camera.URL + "/cgi-bin/CGIProxy.fcgi?usr=" + camera.User + "&pwd=" + camera.Password + "&cmd=setMotionDetectConfig&"

	param := url.Values{}

	param.Add("isEnable", boolConversion(config.IsEnable))
	param.Add("linkage", strconv.FormatUint(uint64(config.Linkage), 10))
	param.Add("snapInterval", strconv.FormatUint(uint64(config.SnapInterval), 10))
	param.Add("sensitivity", strconv.FormatUint(uint64(config.Sensitivity), 10))
	param.Add("triggerInterval", strconv.FormatUint(uint64(config.Triggerinterval), 10))
	param.Add("schedule0", strconv.FormatUint(config.ScheduleMonday, 10))
	param.Add("schedule1", strconv.FormatUint(config.ScheduleTuesday, 10))
	param.Add("schedule2", strconv.FormatUint(config.ScheduleWednesday, 10))
	param.Add("schedule3", strconv.FormatUint(config.ScheduleThursday, 10))
	param.Add("schedule4", strconv.FormatUint(config.ScheduleFriday, 10))
	param.Add("schedule5", strconv.FormatUint(config.ScheduleSaturday, 10))
	param.Add("schedule6", strconv.FormatUint(config.ScheduleSunday, 10))
	param.Add("area0", strconv.FormatUint(uint64(config.Area0), 10))
	param.Add("area1", strconv.FormatUint(uint64(config.Area1), 10))
	param.Add("area2", strconv.FormatUint(uint64(config.Area2), 10))
	param.Add("area3", strconv.FormatUint(uint64(config.Area3), 10))
	param.Add("area4", strconv.FormatUint(uint64(config.Area4), 10))
	param.Add("area5", strconv.FormatUint(uint64(config.Area5), 10))
	param.Add("area6", strconv.FormatUint(uint64(config.Area6), 10))
	param.Add("area7", strconv.FormatUint(uint64(config.Area7), 10))
	param.Add("area8", strconv.FormatUint(uint64(config.Area8), 10))
	param.Add("area9", strconv.FormatUint(uint64(config.Area9), 10))

	request += param.Encode()

	resp, err := http.Get(request)

	if err != nil {
		return false
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false
	}

	v := ReturnValue{}

	err = xml.Unmarshal(b, &v)

	if err != nil {
		return false
	}

	if v.Result == Success {
		return true
	} else {
		return false
	}
}

func (camera *Camera) Move(direction Direction) {
	http.Get("http://" + camera.URL + "/cgi-bin/CGIProxy.fcgi?usr=" + camera.User + "&pwd=" + camera.Password + "&cmd=" + string(direction))
}

func (camera *Camera) GetDeviceInformation() DeviceInformation {
	devInfo := DeviceInformation{}

	request := url.URL{}
	request.Scheme = "http"
	request.Host = camera.URL
	request.Path = "cgi-bin/CGIProxy.fcgi"

	v := url.Values{}
	v.Add("usr", camera.User)
	v.Add("pwd", camera.Password)
	v.Add("cmd", "getDevInfo")
	request.RawQuery = v.Encode()

	fmt.Println(request.String())

	response, err := http.Get(request.String())

	if err != nil {
		return devInfo
	}

	b, err := ioutil.ReadAll(response.Body)

	fmt.Println(string(b))

	if err != nil {
		return devInfo
	}

	err = xml.Unmarshal(b, &devInfo)

	return devInfo
}

func (camera *Camera) SetMotionDetectionEnabled(enabled bool) {
	v := camera.GetMotionDetectConfig()

	if enabled != v.IsEnable {
		v.IsEnable = enabled
		if camera.SetMotionDetectConfig(v) {
			fmt.Println("State updated")
		} else {
			fmt.Println("State update failed. Check credentials?")
		}
	} else {
		fmt.Println("Already in desired state")
	}
}

func (camera *Camera) HandleArgument(arg string) {
	switch arg {
	case "isEnabled":
		config := camera.GetMotionDetectConfig()

		if config.IsEnable == true {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func (camera *Camera) OpenRTSP() {
	cmd := exec.Command("/usr/bin/vlc", "rtsp://"+camera.User+":"+camera.Password+"@"+camera.URL+"/videoSub")
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}
}

func (camera *Camera) OpenMJPEG() {
	request := url.URL{}
	request.Scheme = "http"
	request.Host = camera.URL
	request.Path = "cgi-bin/CGIStream.cgi"

	v := url.Values{}
	v.Add("usr", camera.User)
	v.Add("pwd", camera.Password)
	v.Add("cmd", "GetMJStream")
	request.RawQuery = v.Encode()

	cmd := exec.Command("vlc", request.String())
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}
}

func (camera *Camera) RawCommand(cmd string) {
	request := url.URL{}
	request.Scheme = "http"
	request.Host = camera.URL
	request.Path = "cgi-bin/CGIProxy.fcgi"

	v := url.Values{}
	v.Add("usr", camera.User)
	v.Add("pwd", camera.Password)
	v.Add("cmd", cmd)
	request.RawQuery = v.Encode()

	resp, _ := http.Get(request.String())
	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Print(string(b))
}

func boolConversion(boolean bool) string {
	if boolean {
		return "1"
	} else {
		return "0"
	}
}

func GetCameraConfiguration() (Camera, error) {
	config := Camera{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Camera IP[:PORT]: ")
	url, err := reader.ReadString('\n')

	if err != nil {
		return config, err
	}

	url = strings.TrimRight(url, "\n")
	config.URL = url

	fmt.Print("User: ")
	user, err := reader.ReadString('\n')

	if err != nil {
		return config, err
	}

	user = strings.TrimRight(user, "\n")
	config.User = user

	fmt.Print("Password: ")
	bytes, err := terminal.ReadPassword(0)

	if err != nil {
		return config, err
	}

	config.Password = string(bytes)

	return config, err
}

func ReadConfiguration() Camera {
	var camera Camera

	filePath := os.Getenv("HOME") + "/.foscamrc"
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Creating new configuration file: " + filePath)

		configFile, err = os.Create(filePath)

		if err != nil {
			log.Fatal(err)
		}

		camera, err = GetCameraConfiguration()

		if err != nil {
			log.Fatal(err)
		}

		encoder := json.NewEncoder(configFile)
		err = encoder.Encode(&camera)

		if err != nil {
			log.Fatal(err)
		}

		configFile.Seek(0, 0)
	}

	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&camera)

	if err != nil {
		log.Fatal(err)
	}

	return camera
}
