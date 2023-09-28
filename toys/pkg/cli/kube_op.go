package cli

import (
	"fmt"
)

var _ visitor = (*kubeOpVisitor)(nil)

type kubeOpVisitor struct {
	visitor visitor
}

func newKubeOpVisitor(v visitor) visitor {
	return &kubeOpVisitor{visitor: v}
}

func (k *kubeOpVisitor) visit(fn visitorFn) error {
	return k.visitor.visit(func(info *commandContext, err error) error {
		op := info.restCommands[0]
		switch op {
		case "get":
			info.verb = kubeOpVerbGet
		case "delete":
			info.verb = kubeOpVerbDel
		default:
			return fmt.Errorf("unable to execute 'kubectl %s', error: %w", op, errUnknownOpVerb)
		}
		info.restCommands = info.restCommands[1:]

		if len(info.restCommands) == 1 {
			if info.restCommands[0] == "-h" || info.restCommands[0] == "--help" {
				return errKubeOpPrintHelp
			}
		}

		if err := fn(info, nil); err != nil {
			return err
		}
		return nil
	})
}

var (
	getFlagVisitorFns = make([]visitorFn, 0, 8)
	delFlagVisitorFns = make([]visitorFn, 0, 8)
)

func flagsRegister() {
	getFlagVisitorFns = append(getFlagVisitorFns,
		splitFlagsVisitor,
		getLabelFlagsVisitor,
		getOutputFlagsVisitor,
		getChunkSizeFlagVisitor,
		finallyVerifyFlagsVisitor, // This is the last one of all decorators.
	)
}
