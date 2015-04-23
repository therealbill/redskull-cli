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

var (
	client  *rpcclient.Client
	app     *cli.App
	timeout = time.Second * 2
)

//var rpc_addr string

func main() {
	app = cli.NewApp()
	app.Name = "redskull-cli"
	app.Version = "0.5.1"
	app.EnableBashCompletion = true
	author := cli.Author{Name: "Bill Anderson", Email: "therealbill@me.com"}
	app.Authors = append(app.Authors, author)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "rpcaddr, r",
			Value:  "localhost:8001",
			Usage:  "Redskull RCP address in form 'ip:port'",
			EnvVar: "REDSKULL_RPCADDR",
		},
	}
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
					Name:   "authcheck",
					Usage:  "check auth status of pod",
					Action: CheckPodAuth,
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
				{
					Name:   "addslave",
					Usage:  "Add slave to a pod",
					Action: AddSlaveToPod,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "name of the pod",
						},
						cli.StringFlag{
							Name:  "slaveauth, s",
							Value: "",
							Usage: "auth token for the slave",
						},
						cli.StringFlag{
							Name:  "ip, i",
							Value: "",
							Usage: "ip of the slave",
						},
						cli.IntFlag{
							Name:  "port, p",
							Value: 6379,
							Usage: "port of the slave",
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
}

func RemovePod(c *cli.Context) {
	client, err := rpcclient.NewClient(c.GlobalString("rpcaddr"), timeout)
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
	Sentinels   []string
	HasErrors   bool
}

var templateFuncs = template.FuncMap{"rangeStruct": RangeStructer}

func ShowPod(c *cli.Context) {
	client, err := rpcclient.NewClient(c.GlobalString("rpcaddr"), timeout)
	t := template.Must(template.New("podinfo").Parse(PodInfoTemplate)).Funcs(templateFuncs)
	if err != nil {
		log.Fatal(err)
	}
	pod, err := client.GetPod(c.String("name"))
	if err != nil {
		log.Printf("Pod pull failed with error '%s'", err.Error())
		return
	}
	pod.Master.LastUpdateValid = false
	pod.Master.UpdateData()
	pod.CanFailover()
	scount, sentinels, err := client.GetSentinelsForPod(c.String("name"))
	if err != nil {
		log.Printf("Error on GSFP call: %s", err.Error())
	} else {
		pod.SentinelCount = scount
		log.Printf("Found %d Sentinels: %v", scount, sentinels)
	}
	data := PodData{Pod: pod, CanFailover: pod.CanFailover(), HasErrors: pod.HasErrors(), Sentinels: sentinels}
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func AddPod(c *cli.Context) {
	client, err := rpcclient.NewClient(c.GlobalString("rpcaddr"), timeout)
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
	client, err := rpcclient.NewClient(c.GlobalString("rpcaddr"), timeout)
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

func AddSlaveToPod(c *cli.Context) {
	client, err := rpcclient.NewClient(c.GlobalString("rpcaddr"), timeout)
	if err != nil {
		log.Fatal(err)
	}
	added, err := client.AddSlaveToPod(c.String("name"), c.String("ip"), c.Int("port"), c.String("slaveauth"))
	if err != nil {
		log.Fatal(err)
	}
	if !added {
		log.Print("No error, but not reported as added. Check the server logs for why.")
		return
	}
	log.Printf("Enslaved %s:%d to %s", c.String("ip"), c.Int("port"), c.String("name"))
}

func CheckPodAuth(c *cli.Context) {
	client, err := rpcclient.NewClient(c.GlobalString("rpcaddr"), timeout)
	authmap, err := client.CheckPodAuth(c.String("name"))
	if err != nil {
		log.Print("Unable to get pod auth. Err: ", err)
		return
	}
	allgood := true
	for s, r := range authmap {
		if !r {
			log.Print("%s can not be athenticated to using the pod's auth token", s)
			allgood = false
		}
	}
	if !allgood {
		log.Print("Pod is not considered to have a valid authentication config because at least one node has a different authentication setting")
		return
	}
	log.Print("Auth Valid")
}
