// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Run Hashicorp Terrafrom from Direktiv",
    "title": "terraform",
    "version": "1.0.0",
    "x-direktiv-meta": {
      "categories": [
        "Cloud",
        "Tools"
      ],
      "container": "direktiv/terraform",
      "issues": "https://github.com/direktiv-apps/terraform/issues",
      "license": "[Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)",
      "long-description": null,
      "maintainer": "[direktiv.io](https://www.direktiv.io) ",
      "url": "https://github.com/direktiv-apps/terraform"
    }
  },
  "paths": {
    "/": {
      "post": {
        "parameters": [
          {
            "type": "string",
            "description": "direktiv action id is an UUID. \nFor development it can be set to 'development'\n",
            "name": "Direktiv-ActionID",
            "in": "header"
          },
          {
            "type": "string",
            "description": "direktiv temp dir is the working directory for that request\nFor development it can be set to e.g. '/tmp'\n",
            "name": "Direktiv-TempDir",
            "in": "header"
          },
          {
            "name": "body",
            "in": "body",
            "schema": {
              "type": "object",
              "properties": {
                "commands": {
                  "description": "Commands to execute in order.",
                  "type": "array",
                  "items": {
                    "type": "string"
                  },
                  "example": [
                    "terraform -chdir=out/workflow/tfbase.tar.gz plan"
                  ]
                },
                "continue": {
                  "description": "If set to true all commands are getting executed and errors ignored.",
                  "type": "boolean",
                  "default": false,
                  "example": true
                },
                "variables": {
                  "description": "Variables set for all commands. This translatyes into TF_VAR_* environment variables.",
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "name": {
                        "description": "Name of the variable.",
                        "type": "string"
                      },
                      "value": {
                        "description": "Value of the variable.",
                        "type": "string"
                      }
                    }
                  },
                  "example": [
                    {
                      "name": "instance_name",
                      "value": "myinstance"
                    }
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "nice greeting",
            "schema": {
              "type": "object",
              "additionalProperties": false,
              "example": {
                "greeting": "Hello YourName"
              }
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            },
            "headers": {
              "Direktiv-ErrorCode": {
                "type": "string"
              },
              "Direktiv-ErrorMessage": {
                "type": "string"
              }
            }
          }
        },
        "x-direktiv": {
          "cmds": [
            {
              "action": "foreach",
              "continue": "{{ .Body.Continue }}",
              "env": [
                "TF_IN_AUTOMATION=y",
                "TF_INPUT=0"
              ],
              "exec": "/bin/bash -c \"{{ .Item }}\"",
              "loop": ".Commands",
              "runtime-envs": "[\n{{- range $index, $element := .Body.Variables }}\n{{- if $index}},{{- end}}\n\"TF_VAR_{{ $element.Name }}={{ $element.Value }}\"\n{{- end }}\n]\n"
            }
          ],
          "debug": true,
          "output": "{\n  \"terraform\": {{ index . 0 | toJson }}\n}\n"
        },
        "x-direktiv-errors": {
          "io.direktiv.command.error": "Command execution failed",
          "io.direktiv.output.error": "Template error for output generation of the service",
          "io.direktiv.ri.error": "Can not create information object from request"
        },
        "x-direktiv-examples": [
          {
            "content": "- id: tf\n     type: action\n     action:\n      function: get\n      files:\n      # Contains all required .tf files. Can point to a plain text .tf file as well.\n      - scope: workflow\n        key: tfbase.tar.gz\n        as: tf\n        type: tar.gz\n      input: \n        commands:\n        # the execution dir (chdir) is \"tf\" which we create in the \"files\" section\n        # Storing the state in \"../out/workflow/terraform.tfstate\" will store the state in workflow scope. \n        - terraform -chdir=tf apply -state=../out/workflow/terraform.tfstate -no-color -auto-approve",
            "title": "Basic"
          },
          {
            "content": "- id: tf\n     type: action\n       action:\n        function: get\n        secrets: [\"password\"]\n        files:\n        - scope: workflow\n          key: main.tf\n        input: \n          commands:\n          # Uses tfstate with a jq component. Can run same .tf file for different instances. \n          - terraform apply -state=out/workflow/terraform-jq(.instance).tfstate -no-color -auto-approve\n          # returns state of the change and can be used in a switch later\n          - terraform plan -detailed-exitcode | echo $?\n          variables:\n          - name: instance_name\n            value: jq(.instance)\n          # Use of Direktiv secrets or fetch secrets earlier in the flow.\n          - password:\n            value: jq(.secrets.password)",
            "title": "Example with Variables and Secrets"
          },
          {
            "content": "- id: tf\n     type: action\n       action:\n        function: get\n        files:\n        - scope: workflow\n          key: main.tf\n        input: \n          commands:\n          # return graph as base64\n          - terraform graph | dot -Tpng | base64 -w0\n          # store graph as Direktiv variable\n          - terraform graph | dot -Tpng \u003e out/workflow/graph.png",
            "title": "Visualize"
          }
        ],
        "x-direktiv-function": "functions:\n  - id: terraform\n    image: direktiv/terraform\n    type: knative-workflow"
      },
      "delete": {
        "parameters": [
          {
            "type": "string",
            "description": "On cancel Direktiv sends a DELETE request to\nthe action with id in the header\n",
            "name": "Direktiv-ActionID",
            "in": "header"
          }
        ],
        "responses": {
          "200": {
            "description": ""
          }
        },
        "x-direktiv": {
          "cancel": "echo 'cancel {{ .DirektivActionID }}'"
        }
      }
    }
  },
  "definitions": {
    "direktivFile": {
      "type": "object",
      "x-go-type": {
        "import": {
          "package": "github.com/direktiv/apps/go/pkg/apps"
        },
        "type": "DirektivFile"
      }
    },
    "error": {
      "type": "object",
      "required": [
        "errorCode",
        "errorMessage"
      ],
      "properties": {
        "errorCode": {
          "type": "string"
        },
        "errorMessage": {
          "type": "string"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Run Hashicorp Terrafrom from Direktiv",
    "title": "terraform",
    "version": "1.0.0",
    "x-direktiv-meta": {
      "categories": [
        "Cloud",
        "Tools"
      ],
      "container": "direktiv/terraform",
      "issues": "https://github.com/direktiv-apps/terraform/issues",
      "license": "[Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)",
      "long-description": null,
      "maintainer": "[direktiv.io](https://www.direktiv.io) ",
      "url": "https://github.com/direktiv-apps/terraform"
    }
  },
  "paths": {
    "/": {
      "post": {
        "parameters": [
          {
            "type": "string",
            "description": "direktiv action id is an UUID. \nFor development it can be set to 'development'\n",
            "name": "Direktiv-ActionID",
            "in": "header"
          },
          {
            "type": "string",
            "description": "direktiv temp dir is the working directory for that request\nFor development it can be set to e.g. '/tmp'\n",
            "name": "Direktiv-TempDir",
            "in": "header"
          },
          {
            "name": "body",
            "in": "body",
            "schema": {
              "type": "object",
              "properties": {
                "commands": {
                  "description": "Commands to execute in order.",
                  "type": "array",
                  "items": {
                    "type": "string"
                  },
                  "example": [
                    "terraform -chdir=out/workflow/tfbase.tar.gz plan"
                  ]
                },
                "continue": {
                  "description": "If set to true all commands are getting executed and errors ignored.",
                  "type": "boolean",
                  "default": false,
                  "example": true
                },
                "variables": {
                  "description": "Variables set for all commands. This translatyes into TF_VAR_* environment variables.",
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/VariablesItems0"
                  },
                  "example": [
                    {
                      "name": "instance_name",
                      "value": "myinstance"
                    }
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "nice greeting",
            "schema": {
              "type": "object",
              "additionalProperties": false,
              "example": {
                "greeting": "Hello YourName"
              }
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            },
            "headers": {
              "Direktiv-ErrorCode": {
                "type": "string"
              },
              "Direktiv-ErrorMessage": {
                "type": "string"
              }
            }
          }
        },
        "x-direktiv": {
          "cmds": [
            {
              "action": "foreach",
              "continue": "{{ .Body.Continue }}",
              "env": [
                "TF_IN_AUTOMATION=y",
                "TF_INPUT=0"
              ],
              "exec": "/bin/bash -c \"{{ .Item }}\"",
              "loop": ".Commands",
              "runtime-envs": "[\n{{- range $index, $element := .Body.Variables }}\n{{- if $index}},{{- end}}\n\"TF_VAR_{{ $element.Name }}={{ $element.Value }}\"\n{{- end }}\n]\n"
            }
          ],
          "debug": true,
          "output": "{\n  \"terraform\": {{ index . 0 | toJson }}\n}\n"
        },
        "x-direktiv-errors": {
          "io.direktiv.command.error": "Command execution failed",
          "io.direktiv.output.error": "Template error for output generation of the service",
          "io.direktiv.ri.error": "Can not create information object from request"
        },
        "x-direktiv-examples": [
          {
            "content": "- id: tf\n     type: action\n     action:\n      function: get\n      files:\n      # Contains all required .tf files. Can point to a plain text .tf file as well.\n      - scope: workflow\n        key: tfbase.tar.gz\n        as: tf\n        type: tar.gz\n      input: \n        commands:\n        # the execution dir (chdir) is \"tf\" which we create in the \"files\" section\n        # Storing the state in \"../out/workflow/terraform.tfstate\" will store the state in workflow scope. \n        - terraform -chdir=tf apply -state=../out/workflow/terraform.tfstate -no-color -auto-approve",
            "title": "Basic"
          },
          {
            "content": "- id: tf\n     type: action\n       action:\n        function: get\n        secrets: [\"password\"]\n        files:\n        - scope: workflow\n          key: main.tf\n        input: \n          commands:\n          # Uses tfstate with a jq component. Can run same .tf file for different instances. \n          - terraform apply -state=out/workflow/terraform-jq(.instance).tfstate -no-color -auto-approve\n          # returns state of the change and can be used in a switch later\n          - terraform plan -detailed-exitcode | echo $?\n          variables:\n          - name: instance_name\n            value: jq(.instance)\n          # Use of Direktiv secrets or fetch secrets earlier in the flow.\n          - password:\n            value: jq(.secrets.password)",
            "title": "Example with Variables and Secrets"
          },
          {
            "content": "- id: tf\n     type: action\n       action:\n        function: get\n        files:\n        - scope: workflow\n          key: main.tf\n        input: \n          commands:\n          # return graph as base64\n          - terraform graph | dot -Tpng | base64 -w0\n          # store graph as Direktiv variable\n          - terraform graph | dot -Tpng \u003e out/workflow/graph.png",
            "title": "Visualize"
          }
        ],
        "x-direktiv-function": "functions:\n  - id: terraform\n    image: direktiv/terraform\n    type: knative-workflow"
      },
      "delete": {
        "parameters": [
          {
            "type": "string",
            "description": "On cancel Direktiv sends a DELETE request to\nthe action with id in the header\n",
            "name": "Direktiv-ActionID",
            "in": "header"
          }
        ],
        "responses": {
          "200": {
            "description": ""
          }
        },
        "x-direktiv": {
          "cancel": "echo 'cancel {{ .DirektivActionID }}'"
        }
      }
    }
  },
  "definitions": {
    "VariablesItems0": {
      "type": "object",
      "properties": {
        "name": {
          "description": "Name of the variable.",
          "type": "string"
        },
        "value": {
          "description": "Value of the variable.",
          "type": "string"
        }
      }
    },
    "direktivFile": {
      "type": "object",
      "x-go-type": {
        "import": {
          "package": "github.com/direktiv/apps/go/pkg/apps"
        },
        "type": "DirektivFile"
      }
    },
    "error": {
      "type": "object",
      "required": [
        "errorCode",
        "errorMessage"
      ],
      "properties": {
        "errorCode": {
          "type": "string"
        },
        "errorMessage": {
          "type": "string"
        }
      }
    }
  }
}`))
}
