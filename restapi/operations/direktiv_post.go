package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/direktiv/apps/go/pkg/apps"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

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
	Commands    []interface{}
	DirektivDir string
}

type accParamsTemplate struct {
	PostBody
	Commands    []interface{}
	DirektivDir string
}

type ctxInfo struct {
	cf        context.CancelFunc
	cancelled bool
}

func PostDirektivHandle(params PostParams) middleware.Responder {
	fmt.Printf("params in: %+v", params)
	resp := &PostOKBody{}

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

	sm.Store(*params.DirektivActionID, &ctxInfo{
		cancel,
		false,
	})

	defer sm.Delete(*params.DirektivActionID)

	var responses []interface{}

	var paramsCollector []interface{}
	accParams := accParams{
		params,
		nil,
		ri.Dir(),
	}

	ret, err = runCommand0(ctx, accParams, ri)
	responses = append(responses, ret)

	// if foreach returns an error there is no continue
	cont = false

	if err != nil && !cont {

		errName := cmdErr

		// if the delete function added the cancel tag
		ci, ok := sm.Load(*params.DirektivActionID)
		if ok {
			cinfo, ok := ci.(*ctxInfo)
			if ok && cinfo.cancelled {
				errName = "direktiv.actionCancelled"
				err = fmt.Errorf("action got cancel request")
			}
		}

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

	// validate
	resp.UnmarshalBinary(responseBytes)
	err = resp.Validate(strfmt.Default)

	if err != nil {
		fmt.Printf("error parsing output object: %+v\n", err)
		return generateError(outErr, err)
	}

	return NewPostOK().WithPayload(resp)
}

// foreach command
type LoopStruct0 struct {
	accParams
	Item        interface{}
	DirektivDir string
}

func runCommand0(ctx context.Context,
	params accParams, ri *apps.RequestInfo) ([]map[string]interface{}, error) {

	ri.Logger().Infof("foreach command over .Commands")

	var cmds []map[string]interface{}

	for a := range params.Body.Commands {

		ls := &LoopStruct0{
			params,
			params.Body.Commands[a],
			params.DirektivDir,
		}
		fmt.Printf("object going in command template: %+v\n", ls)

		cmd, err := templateString(`{{ .Item.Command }}`, ls)
		if err != nil {
			ir := make(map[string]interface{})
			ir[successKey] = false
			ir[resultKey] = err.Error()
			cmds = append(cmds, ir)
			continue
		}

		silent := convertTemplateToBool("{{ .Item.Silent }}", ls, false)
		print := convertTemplateToBool("{{ .Item.Print }}", ls, true)
		cont := convertTemplateToBool("{{ .Item.Continue }}", ls, false)
		output := ""

		envs := []string{}
		env0, _ := templateString(`TF_IN_AUTOMATION=y`, ls)
		envs = append(envs, env0)
		env1, _ := templateString(`TF_INPUT=0`, ls)
		envs = append(envs, env1)
		env2, _ := templateString(`TF_LOG={{ default "off" .Body.Loglevel }}`, ls)
		envs = append(envs, env2)
		env3, _ := templateString(`TF_LOG_PATH={{ .DirektivDir }}/out/{{ default "instance" .Body.Scope }}/tf.log`, ls)
		envs = append(envs, env3)

		envTempl, err := templateString(`[
{{- range $index, $element := .Body.Variables }}
{{- if $index}},{{- end}}
"TF_VAR_{{ $element.Name }}={{ $element.Value }}"
{{- end }}
{{- $lenVar := len .Body.Variables }}
{{- $lenEnvs := len .Body.Envs }}
{{- if and (gt $lenVar 0) (gt $lenEnvs 0) }},{{- end}}
{{- range $index, $element := .Body.Envs }}
{{- if $index}},{{- end}}
"{{ $element.Name }}={{ $element.Value }}"
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
