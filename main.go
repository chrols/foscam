package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chrols/foscam/camera"
	"github.com/urfave/cli"
)

// Functions: i.e TODO

func main() {
	cam := camera.ReadConfiguration()

	app := cli.NewApp()
	app.Name = "foscam"
	app.Usage = "cli utility to control foscam cameras"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "motion",
			Usage: "motion alarm commands",
			Subcommands: []cli.Command{
				{
					Name:  "enable",
					Usage: "enable motion alarm",
					Action: func(c *cli.Context) error {
						cam.SetMotionDetectionEnabled(true)
						return nil
					},
				},
				{
					Name:  "disable",
					Usage: "disable motion alarm",
					Action: func(c *cli.Context) error {
						cam.SetMotionDetectionEnabled(false)
						return nil
					},
				},
				{
					Name:  "status",
					Usage: "query motion alarm status",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "code",
							Usage: "Use exit code to indicate motion alarm status",
						},
					},
					Action: func(c *cli.Context) error {
						if cam.GetMotionDetectConfig().IsEnable {
							fmt.Println("motion detection is enabled")
						} else {
							if c.Bool("code") {
								return cli.NewExitError("motion detection is disabled", -1)
							} else {
								fmt.Println("motion detection is disabled")
							}
						}
						return nil
					},
				},
			},
		},
		{
			Name:  "move",
			Usage: "move camera",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "duration",
					Value: "1s",
					Usage: "`DURATION` to use for movement",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:  "left",
					Usage: "move camera left",
					Action: func(c *cli.Context) error {
						cam.Move(camera.MoveLeft)
						time.Sleep(c.GlobalDuration("duration"))
						cam.Move(camera.StopMove)
						return nil
					},
				},
				{
					Name:  "right",
					Usage: "move camera right",
					Action: func(c *cli.Context) error {
						cam.Move(camera.MoveRight)
						time.Sleep(c.GlobalDuration("duration"))
						cam.Move(camera.StopMove)
						return nil
					},
				},
				{
					Name:  "up",
					Usage: "move camera up",
					Action: func(c *cli.Context) error {
						cam.Move(camera.MoveUp)
						time.Sleep(c.GlobalDuration("duration"))
						cam.Move(camera.StopMove)
						return nil
					},
				},
				{
					Name:  "down",
					Usage: "move camera down",
					Action: func(c *cli.Context) error {
						cam.Move(camera.MoveDown)
						cam.Move(camera.StopMove)
						return nil
					},
				},
				{
					Name:  "abort",
					Usage: "abort camera movement",
					Action: func(c *cli.Context) error {
						cam.Move(camera.StopMove)
						return nil
					},
				},
				{
					Name:  "reset",
					Usage: "reset camera position",
					Action: func(c *cli.Context) error {
						cam.Move(camera.Reset)
						return nil
					},
				},
			},
		},

		{
			Name:  "rtsp",
			Usage: "open RTSP-stream",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "command, c",
					Value: "vlc",
					Usage: "`COMMAND` to use for viewing RTSP",
				},
			},
			Action: func(c *cli.Context) error {
				cam.OpenRTSP()
				return nil
			},
		},
		{
			Name:  "mjpeg",
			Usage: "open mjpeg-stream",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "command, c",
					Value: "vlc",
					Usage: "`COMMAND` to use for viewing MJPEG stream",
				},
			},
			Action: func(c *cli.Context) error {
				cam.OpenMJPEG()
				return nil
			},
		},
		{
			Name:  "snapshot",
			Usage: "take a snapshot",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Write to `FILE` instead of stdout",
				},
			},
			Action: func(c *cli.Context) error {
				cam.OpenRTSP()
				return nil
			},
		},
		{
			Name:  "raw",
			Usage: "raw command",
			Action: func(c *cli.Context) error {
				cam.RawCommand(c.Args().Get(0))
				return nil
			},
		},
	}

	app.Run(os.Args)
	return
}
