package cli

import (
	"errors"
	"fmt"
	"strings"
)

func splitFlagsVisitor(info *commandContext, err error) error {
	if len(info.restCommands) <= 0 {
		return nil
	}

	for {
		if len(info.restCommands) <= 0 {
			break
		}

		flag := info.restCommands[0]
		if strings.HasPrefix(flag, "-") {
			info.restCommands = info.restCommands[1:]
			if strings.Contains(flag, "=") {
				// flag with value like '-l=axyn=xxx'
				vals := strings.SplitAfterN(flag, "=", strings.Index(flag, "="))
				flag = strings.TrimSuffix(vals[0], "=")
				if len(info.flags[flag]) == 0 {
					list := make([]string, 0, 4)
					list = append(list, vals[1])
					info.flags[flag] = list
				} else {
					info.flags[flag] = append(info.flags[flag], vals[1])
				}
			} else {
				// flag without value like ['-l', 'axyn=xxx']
				list := make([]string, 0, 4)
				var isBoolFlag bool
				for _, val := range info.restCommands {
					if strings.HasPrefix(val, "-") {
						isBoolFlag = true
						list = append(list, "")
						break
					}
					list = append(list, val)
				}
				if !isBoolFlag && len(list) > 0 {
					info.restCommands = info.restCommands[len(list):]
				}

				if len(info.flags[flag]) == 0 {
					info.flags[flag] = list
				} else {
					info.flags[flag] = append(info.flags[flag], list...)
				}
			}
		}
	}
	return nil
}

func getLabelFlagsVisitor(info *commandContext, err error) error {
	short, sok := info.flags["-l"]
	long, lok := info.flags["--selector"]
	if sok && lok && len(short) > 0 && len(long) > 0 {
		return errors.New("flag '-l' or '--selector' used more than once")
	}

	if sok && !lok {
		if len(short) > 1 {
			return errors.New("flag '-l' or '--selector' used more than once")
		}
		info.queryParameters["labelSelector"] = short[0]
		return nil
	}

	if !sok && lok {
		if len(long) > 1 {
			return errors.New("flag '-l' or '--selector' used more than once")
		}
		info.queryParameters["labelSelector"] = long[0]
		return nil
	}

	var flags = make(map[string][]string, 8)
	for k, v := range info.flags {
		if k == "-l" || k == "--selector" {
			continue
		}
		list := make([]string, 0, len(v))
		list = append(list, v...)
		flags[k] = list
	}
	info.flags = flags

	return nil
}

func getOutputFlagsVisitor(info *commandContext, err error) error {
	short, sok := info.flags["-o"]
	long, lok := info.flags["--output"]

	if !sok && !lok {
		info.outFormat = kubeOutTable
		info.requestHeaders["Accept"] = "application/json;as=Table;v=v1;g=meta.k8s.io,application/json;as=Table;v=v1beta1;g=meta.k8s.io,application/json"
		return nil
	}

	if sok && lok && len(short) > 0 && len(long) > 0 {
		return errors.New("flag '-o' or '--output' used more than once")
	}

	if sok && !lok {
		if len(short) > 1 {
			return errors.New("flag '-o' or '--output' used more than once")
		}
		switch short[0] {
		case "table":
			info.outFormat = kubeOutTable
			info.requestHeaders["Accept"] = "application/json;as=Table;v=v1;g=meta.k8s.io,application/json;as=Table;v=v1beta1;g=meta.k8s.io,application/json"
		case "json":
			info.outFormat = kubeOutJson
			info.requestHeaders["Accept"] = "application/json"
		case "yaml":
			info.outFormat = kubeOutYaml
			info.requestHeaders["Accept"] = "application/json"
		default:
			return errors.New("unknown get output format type '" + short[0] + "'")
		}
		return nil
	}

	if !sok && lok {
		if len(long) > 1 {
			return errors.New("flag '-o' or '--output' used more than once")
		}
		switch long[0] {
		case "table":
			info.outFormat = kubeOutTable
			info.requestHeaders["Accept"] = "application/json;as=Table;v=v1;g=meta.k8s.io,application/json;as=Table;v=v1beta1;g=meta.k8s.io,application/json"
		case "json":
			info.outFormat = kubeOutJson
			info.requestHeaders["Accept"] = "application/json"
		case "yaml":
			info.outFormat = kubeOutYaml
			info.requestHeaders["Accept"] = "application/json"
		default:
			return errors.New("unknown get output format type '" + long[0] + "'")
		}
		return nil
	}

	var flags = make(map[string][]string, 8)
	for k, v := range info.flags {
		if k == "-o" || k == "--output" {
			continue
		}
		list := make([]string, 0, len(v))
		list = append(list, v...)
		flags[k] = list
	}
	info.flags = flags

	return nil
}

func getChunkSizeFlagVisitor(info *commandContext, err error) error {
	long, lok := info.flags["--chunk-size"]
	if lok && len(long) > 1 {
		return errors.New("flag '--chunk-size' used more than once")
	}

	if !lok {
		info.queryParameters["limit"] = "500"
		return nil
	}

	if len(long) > 0 {
		if long[0] == "0" {
			// Not add limit query parameter as disable.
			return nil
		} else if long[0] == "" {
			return errors.New("flag '--chunk-size' missing value")
		}
	}

	info.queryParameters["limit"] = long[0]

	var flags = make(map[string][]string, 8)
	for k, v := range info.flags {
		if k == "--chunk-size" {
			continue
		}
		list := make([]string, 0, len(v))
		list = append(list, v...)
		flags[k] = list
	}
	info.flags = flags

	return nil
}

func finallyVerifyFlagsVisitor(info *commandContext, err error) error {
	if len(info.flags) > 0 {
		var errmsg string
		for k, v := range info.flags {
			if len(errmsg) == 0 {
				errmsg = fmt.Sprintf("[unknown flag] %s and values (%+v)", k, v)
				continue
			}
			errmsg = fmt.Sprintf("%s ; [unknown flag] %s and values (%+v)", errmsg, k, v)
		}
	}
	return nil
}
