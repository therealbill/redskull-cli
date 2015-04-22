package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/therealbill/redskull/actions"
	rpcclient "github.com/therealbill/redskull/rpcclient"
)

var client *rpcclient.Client

var timeout = time.Second * 2

func main() {
	app := cli.NewApp()
	app.Name = "redskull-cli"
	app.Version = "0.5.0"
	app.EnableBashCompletion = true
	author := cli.Author{Name: "Bill Anderson", Email: "therealbill@me.com"}
	app.Authors = append(app.Authors, author)

	app.Commands = []cli.Command{
		{
			Name:  "pod",
			Usage: "Pod specific actions",
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "show pod info",
					Action: ShowPod,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "name of the pod",
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "Remove pod from sentinels",
					Action: RemovePod,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "name of the pod",
						},
					},
				},
				{
					Name:   "add",
					Usage:  "Add a pod to be managed",
					Action: AddPod,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "name of the pod",
						},
						cli.StringFlag{
							Name:  "auth, a",
							Value: "",
							Usage: "auth token for the pod",
						},
						cli.StringFlag{
							Name:  "ip, i",
							Value: "",
							Usage: "ip of the pod",
						},
						cli.IntFlag{
							Name:  "quorum, q",
							Value: 2,
							Usage: "sentinel quorum required",
						},
						cli.IntFlag{
							Name:  "port, p",
							Value: 6379,
							Usage: "port of the pod",
						},
					},
				},
			},
		},
		{
			Name:  "sentinel",
			Usage: "sentinel specific actions",
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Add a sentinel to the cluster",
					Action: AddSentinel,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "name of the sentinel",
						},
						cli.StringFlag{
							Name:  "ip, i",
							Value: "",
							Usage: "ip address of the sentinel",
						},
						cli.IntFlag{
							Name:  "port, p",
							Value: 26379,
							Usage: "port for the snetinel",
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)

	/*

		pod, err = client.GetPod("pod1")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("pod: %+v", pod)

		log.Print("Initiating Re-balance\n")
		err = client.BalancePod("pod1")
		if err != nil {
			log.Printf("Unable to request rebalance of pod1")
		} else {
			log.Print("Rebalance initiated")
		}
	*/
}

func RemovePod(c *cli.Context) {
	client, err := rpcclient.NewClient("127.0.0.1:8001", timeout)
	name := c.String("name")
	log.Printf("Removing pod %s", name)
	err = client.RemovePod(name)
	if err != nil {
		log.Printf("Unable to remove pod")
	} else {
		log.Print("removed")
	}

}

type PodData struct {
	Pod         actions.RedisPod
	CanFailover bool
	HasErrors   bool
}

var templateFuncs = template.FuncMap{"rangeStruct": RangeStructer}

func ShowPod(c *cli.Context) {
	client, err := rpcclient.NewClient("127.0.0.1:8001", timeout)
	t := template.Must(template.New("podinfo").Parse(PodInfoTemplate)).Funcs(templateFuncs)
	if err != nil {
		log.Fatal(err)
	}
	pod, err := client.GetPod(c.String("name"))
	pod.Master.LastUpdateValid = false
	pod.Master.UpdateData()
	pod.CanFailover()
	data := PodData{Pod: pod, CanFailover: pod.CanFailover(), HasErrors: pod.HasErrors()}
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func AddPod(c *cli.Context) {
	client, err := rpcclient.NewClient("127.0.0.1:8001", timeout)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.AddPod(c.String("name"), c.String("ip"), c.Int("port"), c.Int("quorum"), c.String("auth"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Res: %+v", res)
}

func AddSentinel(c *cli.Context) {
	client, err := rpcclient.NewClient("127.0.0.1:8001", timeout)
	if err != nil {
		log.Fatal(err)
	}
	addr := fmt.Sprintf("%s:%d", c.String("ip"), c.Int("port"))
	ok, err := client.AddSentinel(addr)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Print("server reports unable to add sentinel, no error given. :(")
		return
	}
	log.Printf("Added sentinel %s:%d", c.String("ip"), c.Int("port"))
}
