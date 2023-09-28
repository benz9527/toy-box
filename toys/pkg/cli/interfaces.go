package cli

import (
	"errors"
)

type kubeOpVerb uint

const (
	kubeOpVerbNil kubeOpVerb = iota
	kubeOpVerbGet
	kubeOpVerbDel
)

type kubeResultOutFormatType uint

const (
	kubeOutTable kubeResultOutFormatType = iota
	kubeOutJson
	kubeOutYaml
)

var (
	errForbidCli          = errors.New("forbid cli to execute")
	errUnknownOpVerb      = errors.New("unknown kube op")
	errCliTerminated      = errors.New("cli parse terminated")
	errKubeCmdPrintHelp   = errors.New("cli print help")
	errKubeOpPrintHelp    = errors.New("cli op print help")
	errKubeCliForbid      = errors.New("cli forbid to execute")
	errCliExecTimeout     = errors.New("execution timeout")
	errCliExecEmptyResult = errors.New("empty result")
)

type Commander interface {
	SetArguments(args []string, maps ...*CliKubeResourceMapsDTO)
	Exec() (string, error)
}

var _ visitor = (*commandContext)(nil)

type commandContext struct {
	resource        string
	result          string // After formatted.
	subResource     string // Maybe join by comma?
	verb            kubeOpVerb
	outFormat       kubeResultOutFormatType
	resourceMaps    *CliKubeResourceMapsDTO // Release forbid
	genericObj      []byte
	restCommands    []string
	requestHeaders  map[string]string
	flags           map[string][]string // Temporarily store the flags
	queryParameters map[string]string
}

func newCommandInfo(commands ...string) *commandContext {
	list := make([]string, 0, len(commands))
	list = append(list, commands...)
	return &commandContext{
		restCommands:    list,
		requestHeaders:  make(map[string]string, 8),
		flags:           make(map[string][]string, 8),
		queryParameters: make(map[string]string, 8),
	}
}

func (i *commandContext) visit(fn visitorFn) error {
	return fn(i, nil)
}

func (i *commandContext) release() {
	// For GC.
	i.flags = nil
	i.restCommands = nil
	i.requestHeaders = nil
	i.queryParameters = nil
	i.genericObj = nil
}

type visitorFn func(*commandContext, error) error

type visitor interface {
	visit(fn visitorFn) error
}

type decoratedVisitor struct {
	visitor    visitor
	decorators []visitorFn // Loop and execute in order.
}

func newDecoratedVisitor(base visitor, vfns ...visitorFn) visitor {
	if len(vfns) == 0 {
		return base
	}
	return decoratedVisitor{base, vfns}
}

var _ visitor = &decoratedVisitor{}

func (d decoratedVisitor) visit(fn visitorFn) error {
	return d.visitor.visit(func(info *commandContext, err error) error {
		if err != nil {
			return err
		}
		if err = fn(info, nil); err != nil {
			return err
		}
		for i := range d.decorators {
			if err = d.decorators[i](info, nil); err != nil {
				return err
			}
		}
		return nil
	})
}
