package cli

import (
	"errors"
	"fmt"
	"strings"
)

var _ visitor = (*kubeResourceVisitor)(nil)

type kubeResourceVisitor struct {
	visitor decoratedVisitor
}

func newKubeResourceVisitor(v visitor) visitor {
	return &kubeResourceVisitor{visitor: decoratedVisitor{visitor: v}}
}

func (k *kubeResourceVisitor) visit(fn visitorFn) error {
	return k.visitor.visit(func(info *commandContext, err error) error {
		res := info.restCommands[0]
		if len(res) == 0 {
			return errors.New("missing resource after op in command")
		}
		if strings.HasPrefix(res, "-") {
			return errors.New("unknown resource start with '-'")
		}

		// resource plural or resource alias (i.e. short name) mapped as resource plural
		var (
			ok           bool
			completedRes any
		)
		switch info.verb {
		case kubeOpVerbGet:
			if info.resourceMaps.GetResourcesMap == nil {
				return fmt.Errorf("internal error %w", errCliTerminated)
			}
			completedRes, ok = info.resourceMaps.GetResourcesMap.Load(res)
		case kubeOpVerbDel:
			if info.resourceMaps.DelResourcesMap == nil {
				return fmt.Errorf("internal error %w", errCliTerminated)
			}
			completedRes, ok = info.resourceMaps.DelResourcesMap.Load(res)
		default:
			return errCliTerminated
		}

		if !ok {
			return fmt.Errorf("unable to find resource '%s' to execute, error: %w", res, errCliTerminated)
		}
		info.restCommands = info.restCommands[1:]
		info.resource = completedRes.(string)

		// Sub-resource fetch.
		var sub string
		if len(info.restCommands) > 0 {
			sub = info.restCommands[0]
		}
		if len(sub) > 0 && !strings.HasPrefix(sub, "-") {
			info.restCommands = info.restCommands[1:]
			info.subResource = sub
		}

		switch info.verb {
		case kubeOpVerbGet:
			k.visitor.decorators = append(k.visitor.decorators, getFlagVisitorFns...)
		case kubeOpVerbDel:
			k.visitor.decorators = append(k.visitor.decorators, delFlagVisitorFns...)
		}

		for _, dfn := range k.visitor.decorators {
			if err = dfn(info, nil); err != nil {
				return err
			}
		}

		if err = fn(info, nil); err != nil { // Do request.
			return err
		}

		return outputFormatFactory(info.outFormat)(info, nil)
	})
}
