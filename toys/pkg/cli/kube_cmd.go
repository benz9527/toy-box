package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/rest"
)

var _ Commander = (*kubeCommander)(nil)

// TODO(Ben) Pooling this resource.
type kubeCommander struct {
	cmd        string
	restArgs   []string
	postSetErr error
	info       *commandContext
	request    *rest.Request
}

func newKubeCommander(request *rest.Request) Commander {
	return &kubeCommander{
		cmd:     "kubectl",
		request: request,
	}
}

func newOcCommander(request *rest.Request) Commander {
	return &kubeCommander{
		cmd:     "oc",
		request: request,
	}
}

func (k *kubeCommander) SetArguments(args []string, maps ...*CliKubeResourceMapsDTO) {
	k.restArgs = args
	k.info = newCommandInfo(k.restArgs...)
	k.info.resourceMaps = &CliKubeResourceMapsDTO{
		AddResourcesMap:  &sync.Map{},
		DelResourcesMap:  &sync.Map{},
		GetResourcesMap:  &sync.Map{},
		ModResourcesMap:  &sync.Map{},
		ListResourcesMap: &sync.Map{},
	}
	if len(maps) == 0 {
		return
	}
	// Avoid the nil panic.
	if maps[0].AddResourcesMap != nil {
		k.info.resourceMaps.AddResourcesMap = maps[0].AddResourcesMap
	}
	if maps[0].DelResourcesMap != nil {
		k.info.resourceMaps.DelResourcesMap = maps[0].DelResourcesMap
	}
	if maps[0].GetResourcesMap != nil {
		k.info.resourceMaps.GetResourcesMap = maps[0].GetResourcesMap
	}
	if maps[0].ModResourcesMap != nil {
		k.info.resourceMaps.ModResourcesMap = maps[0].ModResourcesMap
	}
	if maps[0].ListResourcesMap != nil {
		k.info.resourceMaps.ListResourcesMap = maps[0].ListResourcesMap
	}
}

func (k *kubeCommander) Exec() (string, error) {
	var v visitor = newKubeValidateVisitor(k.info)
	v = newKubeOpVisitor(v)
	v = newKubeResourceVisitor(v)
	defer k.info.release()
	err := v.visit(func(info *commandContext, err error) error {
		var reqVerb string
		switch info.verb {
		case kubeOpVerbGet:
			reqVerb = "GET"
		case kubeOpVerbDel:
			fallthrough
		default:
			return errKubeCliForbid
		}
		req := k.request.Verb(reqVerb).Resource(info.resource)

		// Add request headers if present.
		if len(info.requestHeaders) > 0 {
			for k, v := range info.requestHeaders {
				req.SetHeader(k, v)
			}
		}

		// Add query parameters if present.
		if len(info.queryParameters) > 0 {
			for k, v := range info.queryParameters {
				req.Param(k, v)
			}
		}

		// Sub-resource if present.
		if len(info.subResource) > 0 {
			req.Name(info.subResource)
		}

		todo := context.TODO()
		todoCancel, cancelFn := context.WithTimeout(todo, 5*time.Second)
		defer cancelFn()
		info.genericObj, err = req.Do(todoCancel).Raw()
		if errors.Is(err, context.DeadlineExceeded) {
			return errCliExecTimeout
		}
		return err
	})
	if errors.Is(err, errKubeCmdPrintHelp) {
		return k.printHelp()
	}
	if errors.Is(err, errKubeOpPrintHelp) {
		return k.printOpHelp(k.info.verb, k.info.resourceMaps)
	}

	if err != nil {
		return err.Error(), err
	}

	return k.info.result, nil
}

func (k *kubeCommander) printHelp() (string, error) {
	return fmt.Sprintf(`This is not an actual 'kubectl' of Kubernetes, only works for current system
to mock session management service operation in Kubernetes Cluster.

Command format:
%s <op> <resource> [subresource] [<option> ...]

Sub-command as '<op>', sub-command list:
- get
- delete
# delete sub-command is in the works.
`, k.cmd), nil
}

func (k *kubeCommander) printOpHelp(opVerb kubeOpVerb, maps *CliKubeResourceMapsDTO) (string, error) {
	switch opVerb {
	case kubeOpVerbGet:
		return opGetHelp(maps.GetResourcesMap, k.cmd), nil
	case kubeOpVerbDel:
		return opDelHelp(maps.DelResourcesMap, k.cmd), nil
	default:

	}
	err := errors.New("not help to print")
	return err.Error(), err
}

func opGetHelp(m *sync.Map, indicators ...string) string {
	if len(indicators) <= 0 {
		indicators = []string{"kubectl"}
	}
	if len(strings.TrimSpace(indicators[0])) <= 0 {
		indicators[0] = "kubectl"
	}
	tips := fmt.Sprintf(`Command format:
%s get <resource> [subresource] [<option> ...]

Resources:
`, indicators[0])
	if m != nil {
		m.Range(func(key, value any) bool {
			tips += key.(string) + "\n"
			return true
		})
	}
	tips += `
Options:
-o,--output    Options: table,json,yaml. Default is 'table' if missing.
-l,--selector  Example: '-l axyn=xxx', '-l=axyn=xxx', '-l axyn=xxx,key2=val2'.
               More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#api
--chunk-size   Return large lists in chunks rather than all at once. Pass 0 to disable. Default is 500.
               This flag is beta and may change in the future.

Note:
Namespace is fixed and forbid to set on current system namespace scope cli.
`
	return tips
}

func opDelHelp(m *sync.Map, indicators ...string) string {
	if len(indicators) <= 0 {
		indicators = []string{"kubectl"}
	}
	if len(strings.TrimSpace(indicators[0])) <= 0 {
		indicators[0] = "kubectl"
	}
	tips := fmt.Sprintf(`Command format:
%s delete <resource> [subresource] [<option> ...]

Resources:
`, indicators[0])
	if m != nil {
		m.Range(func(key, value any) bool {
			tips += key.(string) + "\n"
			return true
		})
	}
	tips += `
Note:
Namespace is fixed and forbid to set on current system namespace scope cli.
Delete operation is temporarily forbidden.
`
	return tips
}

func init() {
	flagsRegister()
}
