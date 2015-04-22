package main

import "reflect"

var PodInfoTemplate = `
# Name: {{.Pod.Name}}
RunID: {{.Pod.Info.Runid}}
Quorum: {{.Pod.Info.Quorum}}
Config Epoch: {{.Pod.Info.ConfigEpoch}}
DownAfter: {{.Pod.Info.DownAfterMilliseconds}}ms
Current Master: {{.Pod.Master.Name}}
Can AUTH master: {{ .Pod.Master.HasValidAuth}}
SentinelCount: {{.Pod.SentinelCount}}
Has Errors: {{.HasErrors}}

# Replication
Role: {{.Pod.Info.RoleReported}}

## Slaves
{{printf "%-16s" "IP"}}{{printf "%8s" "PORT"}}{{printf "%10s" "STATE"}}{{printf "%12s" "OFFSET"}}{{printf "%10s" "LAG"}}
{{ range .Pod.Master.Info.Replication.Slaves }}{{ printf "%-16s" .IP}}{{printf "%8d" .Port}}{{printf "%10s" .State}}{{printf "%12d" .Offset}}{{printf "%10d" .Lag}}
{{ end}}

# Stats
EvictedKeys:                   {{printf "%12d" .Pod.Master.Info.Stats.EvictedKeys}}
ExpiredKeys:                   {{printf "%12d" .Pod.Master.Info.Stats.ExpiredKeys}}
InstanteousInputKbps:          {{printf "%12d" .Pod.Master.Info.Stats.InstanteousInputKbps}}
InstanteousOpsPerSecond:       {{printf "%12d" .Pod.Master.Info.Stats.InstanteousOpsPerSecond}}
InstanteousOutputKbps:         {{printf "%12d" .Pod.Master.Info.Stats.InstanteousOutputKbps}}
KeyspaceHits:                  {{printf "%12d" .Pod.Master.Info.Stats.KeyspaceHits}}
KeyspaceMisses:                {{printf "%12d" .Pod.Master.Info.Stats.KeyspaceMisses}}
LatestForkUsec:                {{printf "%12d" .Pod.Master.Info.Stats.LatestForkUsec}}
PubSubChannels:                {{printf "%12d" .Pod.Master.Info.Stats.PubSubChannels}}
PubSubPatterns:                {{printf "%12d" .Pod.Master.Info.Stats.PubSubPatterns}}
RejectedConnections:           {{printf "%12d" .Pod.Master.Info.Stats.RejectedConnections}}
SyncFill:                      {{printf "%12d" .Pod.Master.Info.Stats.SyncFill}}
SyncPartialErr:                {{printf "%12d" .Pod.Master.Info.Stats.SyncPartialErr}}
SyncPartialOk:                 {{printf "%12d" .Pod.Master.Info.Stats.SyncPartialOk}}
TotalCommandsProcessed:        {{printf "%12d" .Pod.Master.Info.Stats.TotalCommandsProcessed}}
TotalConnectionsRecevied:      {{printf "%12d" .Pod.Master.Info.Stats.TotalConnectionsRecevied}}
TotalNetInputBytes:            {{printf "%12d" .Pod.Master.Info.Stats.TotalNetInputBytes}}
TotalNetOutputBytes:           {{printf "%12d" .Pod.Master.Info.Stats.TotalNetOutputBytes}}


`

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}
