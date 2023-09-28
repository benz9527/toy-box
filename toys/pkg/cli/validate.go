package cli

var _ visitor = (*kubeValidateVisitor)(nil)

type kubeValidateVisitor struct {
	visitor visitor
}

func newKubeValidateVisitor(v visitor) *kubeValidateVisitor {
	return &kubeValidateVisitor{visitor: v}
}

func (k *kubeValidateVisitor) visit(fn visitorFn) error {
	return k.visitor.visit(func(info *commandContext, err error) error {
		if len(info.restCommands) <= 0 {
			return errCliTerminated
		}

		if len(info.restCommands) == 1 {
			if info.restCommands[0] == "-h" || info.restCommands[0] == "--help" {
				return errKubeCmdPrintHelp
			}
		}

		if err = fn(info, nil); err != nil {
			return err
		}

		return nil
	})
}
