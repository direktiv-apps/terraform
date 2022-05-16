package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/direktiv/apps/go/pkg/apps"
	"github.com/go-openapi/runtime/middleware"

	"terraform/models"
)

const (
	successKey = "success"
	resultKey  = "result"

	// http related
	statusKey  = "status"
	codeKey    = "code"
	headersKey = "headers"
)

var sm sync.Map

const (
	cmdErr = "io.direktiv.command.error"
	outErr = "io.direktiv.output.error"
	riErr  = "io.direktiv.ri.error"
)

type accParams struct {
	PostParams
	Commands []interface{}
}

type accParamsTemplate struct {
	PostBody
	Commands []interface{}
}

func PostDirektivHandle(params PostParams) middleware.Responder {
	fmt.Printf("params in: %+v", params)
	var resp interface{}

	var (
		err  error
		ret  interface{}
		cont bool
	)

	ri, err := apps.RequestinfoFromRequest(params.HTTPRequest)
	if err != nil {
		return generateError(riErr, err)
	}

	ctx, cancel := context.WithCancel(params.HTTPRequest.Context())
	sm.Store(*params.DirektivActionID, cancel)
	defer sm.Delete(params.DirektivActionID)

	var responses []interface{}

	var paramsCollector []interface{}
	accParams := accParams{
		params,
		nil,
	}

	ret, err = runCommand0(ctx, accParams, ri)
	responses = append(responses, ret)

	cont = convertTemplateToBool("{{ .Body.Continue }}", accParams, true)

	if err != nil && !cont {
		errName := cmdErr
		return generateError(errName, err)
	}

	paramsCollector = append(paramsCollector, ret)
	accParams.Commands = paramsCollector

	fmt.Printf("object going in output template: %+v\n", responses)

	s, err := templateString(`{
  "terraform": {{ index . 0 | toJson }}
}
`, responses)
	if err != nil {
		return generateError(outErr, err)
	}
	fmt.Printf("object from output template: %+v\n", s)

	responseBytes := []byte(s)

	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		fmt.Printf("error parsing output template: %+v\n", err)
		return generateError(outErr, err)
	}

	return NewPostOK().WithPayload(resp)
}

// foreach command
type LoopStruct0 struct {
	accParams
	Item interface{}
}

func runCommand0(ctx context.Context,
	params accParams, ri *apps.RequestInfo) ([]map[string]interface{}, error) {

	ri.Logger().Infof("foreach command over .Commands")

	var cmds []map[string]interface{}

	for a := range params.Body.Commands {

		ls := &LoopStruct0{
			params,
			params.Body.Commands[a],
		}
		fmt.Printf("object going in command template: %+v\n", ls)

		cmd, err := templateString(`/bin/bash -c "{{ .Item }}"`, ls)
		if err != nil {
			ir := make(map[string]interface{})
			ir[successKey] = false
			ir[resultKey] = err.Error()
			cmds = append(cmds, ir)
			continue
		}

		silent := convertTemplateToBool("<no value>", ls, false)
		print := convertTemplateToBool("<no value>", ls, true)
		cont := convertTemplateToBool("{{ .Body.Continue }}", ls, false)
		output := ""

		envs := []string{}
		env0, _ := templateString(`TF_IN_AUTOMATION=y`, ls)
		envs = append(envs, env0)
		env1, _ := templateString(`TF_INPUT=0`, ls)
		envs = append(envs, env1)

		envTempl, err := templateString(`[
{{- range $index, $element := .Body.Variables }}
{{- if $index}},{{- end}}
"TF_VAR_{{ $element.Name }}={{ $element.Value }}"
{{- end }}
]
`, ls)
		if err != nil {
			ir := make(map[string]interface{})
			ir[successKey] = false
			ir[resultKey] = err.Error()
			cmds = append(cmds, ir)
			continue
		}
		var addEnvs []string
		err = json.Unmarshal([]byte(envTempl), &addEnvs)
		if err != nil {
			ir := make(map[string]interface{})
			ir[successKey] = false
			ir[resultKey] = err.Error()
			cmds = append(cmds, ir)
			continue
		}
		envs = append(envs, addEnvs...)

		r, err := runCmd(ctx, cmd, envs, output, silent, print, ri)
		if err != nil {
			ir := make(map[string]interface{})
			ir[successKey] = false
			ir[resultKey] = err.Error()
			cmds = append(cmds, ir)

			if cont {
				continue
			}

			return cmds, err

		}
		cmds = append(cmds, r)

	}

	return cmds, nil

}

// end commands

func generateError(code string, err error) *PostDefault {

	d := NewPostDefault(0).WithDirektivErrorCode(code).
		WithDirektivErrorMessage(err.Error())

	errString := err.Error()

	errResp := models.Error{
		ErrorCode:    &code,
		ErrorMessage: &errString,
	}

	d.SetPayload(&errResp)

	return d
}

func HandleShutdown() {
	// nothing for generated functions
}
