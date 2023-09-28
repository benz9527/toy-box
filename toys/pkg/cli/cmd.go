package cli

import (
	"github.com/samber/lo"
	KRest "k8s.io/client-go/rest"
)

var (
	cmdWhiteList = []string{
		"kubectl", "oc",
	}
	commanderMapping = make(map[string]func(request *KRest.Request) Commander, 8)
)

// Format:
// - kubectl <op> <resource> [(<flag> <option>)...]
// Others format (need implemented):
// - cat <path>

func NewDecoratedCommander(cmdPrefix string, requests ...*KRest.Request) Commander {
	if !isExecutable(cmdPrefix) {
		return &defaultCommander{cmd: cmdPrefix, k8sRequestNum: -1}
	}
	c, ok := commanderMapping[cmdPrefix]
	if !ok {
		return &defaultCommander{cmd: cmdPrefix, k8sRequestNum: -1}
	}
	if len(requests) < 1 {
		return &defaultCommander{cmd: cmdPrefix, k8sRequestNum: len(requests)}
	}
	commander := c(requests[0])
	return commander
}

func isExecutable(cmd string) bool {
	return lo.Contains[string](cmdWhiteList, cmd)
}

var _ Commander = (*defaultCommander)(nil)

type defaultCommander struct {
	cmd           string
	k8sRequestNum int // -1 means not the kubectl cli; it >= 0 but it < 1 means kubectl cli without args after 'kubectl'
}

func (d *defaultCommander) SetArguments(args []string, maps ...*CliKubeResourceMapsDTO) {
	// do nothing
}

func (d *defaultCommander) Exec() (string, error) {
	if d.k8sRequestNum <= -1 {
		return "default unable to execute '" + d.cmd + "'", errForbidCli
	}
	return "unable to execute '" + d.cmd + "' with empty condition", errCliTerminated
}

func init() {
	commanderMapping["kubectl"] = newKubeCommander
	commanderMapping["oc"] = newOcCommander
}
