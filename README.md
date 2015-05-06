[![Build Status](https://travis-ci.org/therealbill/redskull-cli.svg?branch=master)](https://travis-ci.org/therealbill/redskull-cli)


# redskull-cli
This is a CLI Tool for interacting with Redskull via it's TCP-RPC channel. This
port is 1 higher than the bind port, so the default is 8001.

Like several tools it is a command/subcommand based utility (such as git).


#Usage

Here is the output of `redskull-cli -h`:

```
NAME:
   redskull-cli - A new cli application

USAGE:
   redskull-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR(S): 
   Bill Anderson <therealbill@me.com> 
   
COMMANDS:
   pod          Pod specific actions
   sentinel     sentinel specific actions
   help, h      Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --rpcaddr, -r "localhost:8001"       Redskull RCP address in form 'ip:port' [$REDSKULL_RPCADDR]
   --help, -h                   show help
   --generate-bash-completion
   --version, -v                print the version
```   

## Selecting your RedSkull RPC Address

If running this on a redskull node running all defaults, no work needs to be done.

If however you are running this elsewhere or have changed the RPC port Redskull
listens on, you can either set the environment variable "REDSKULL_RPCADDR" to
the `IP:PORT` address of the RCP endpoint (which can be a TCP load balancer,
DNS RR, etc. ) or pass the option as listed in the `help` output.

## Example

```shell
redskull-cli pod show pod1
# Name: pod1
RunID: cde80aa822975a4067efaa1550fed2f1225e8f78
Quorum: 2
Config Epoch: 0
DownAfter: 30000ms
Current Master: 127.0.0.1:6501
Can AUTH master: true
SentinelCount: 0
Has Errors: true

# Replication
Role: master

## Slaves
IP                  PORT     STATE      OFFSET       LAG
127.0.0.1           6505    online      830455         0


# Stats
EvictedKeys:                              0
ExpiredKeys:                              0
InstanteousInputKbps:                     0
InstanteousOpsPerSecond:                  5
InstanteousOutputKbps:                    0
KeyspaceHits:                             0
KeyspaceMisses:                           0
LatestForkUsec:                         342
PubSubChannels:                           1
PubSubPatterns:                           0
RejectedConnections:                      0
SyncFill:                                 1
SyncPartialErr:                           0
SyncPartialOk:                            0
TotalCommandsProcessed:               34507
TotalConnectionsRecevied:               148
TotalNetInputBytes:                 1664771
TotalNetOutputBytes:                8416161
```


## Adding a Sentinel

Useful for puppet recipes, you can add a sentinel via the CLI without needing
to add a pod with it. For example, to add a sentinel in the default port on IP
1.2.3.5: `redis-cli sentinel add -n 1.2.3.5`

# Bash Completion

To enable bash completion, copy redskull-cli-complete.sh into your completions directory and source it.

# TODO/BUGS
 * Still need several functions via the Server API exposed such as cloning, and more detailed RS determined data
 * More/Better output in `pod show`
