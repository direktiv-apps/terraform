
# terraform 1.0

Run Hashicorp Terrafrom from Direktiv

---
- #### Categories: cloud, tools, infrastructure
- #### Image: gcr.io/direktiv/apps/terraform 
- #### License: [Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)
- #### Issue Tracking: https://github.com/direktiv-apps/terraform/issues
- #### URL: https://github.com/direktiv-apps/terraform
- #### Maintainer: [direktiv.io](https://www.direktiv.io) 
---

## About terraform

<no value>

### Example(s)
  #### Function Configuration
  ```yaml
  functions:
  - id: terraform
    image: gcr.io/direktiv/apps/terraform:1.0
    type: knative-workflow
  ```
   #### Basic
   ```yaml
   - id: tf
     type: action
      action:
        files:
        # Contains all required .tf files after init. Can point to a plain text .tf file as well.
        - scope: workflow
          key: tf.tar.gz
          as: out/workflow/tf
          type: tar.gz
        function: get
        input: 
          variables:
          - name: name
            value: MyName
          commands:
          - command: terraform -chdir=out/workflow/tf apply -no-color -auto-approve
          - command: terraform -chdir=out/workflow/tf output -json
   ```
   #### Example with Variables and Secrets
   ```yaml
   - id: tf
     type: action
       action:
        function: get
        secrets: ["password"]
        files:
        - scope: workflow
          key: main.tf
        input: 
          commands:
          # Uses tfstate with a jq component. Can run same .tf file for different instances. 
          - terraform apply -state=out/workflow/terraform-jq(.instance).tfstate -no-color -auto-approve
          # returns state of the change and can be used in a switch later
          - terraform plan -detailed-exitcode | echo $?
          variables:
          - name: instance_name
            value: jq(.instance)
          # Use of Direktiv secrets or fetch secrets earlier in the flow.
          - password:
            value: jq(.secrets.password)
   ```
   #### Visualize
   ```yaml
   - id: tf
     type: action
       action:
        function: get
        files:
        - scope: workflow
          key: main.tf
        input: 
          commands:
          # return graph as base64
          - terraform graph | dot -Tpng | base64 -w0
          # store graph as Direktiv variable
          - terraform graph | dot -Tpng > out/workflow/graph.png
   ```

### Request



#### Request Attributes
[PostParamsBody](#post-params-body)

### Response
  List of executed commands.
#### Reponse Types
    
  

[PostOKBody](#post-o-k-body)
#### Example Reponses
    
```json
[
  {
    "result": "VTQ3U....c2ZaN0FJaldjVnkra2tKV==",
    "success": true
  },
  {
    "format_version": "1.0",
    "result": null,
    "success": true
  }
]
```

### Errors
| Type | Description
|------|---------|
| io.direktiv.command.error | Command execution failed |
| io.direktiv.output.error | Template error for output generation of the service |
| io.direktiv.ri.error | Can not create information object from request |


### Types
#### <span id="post-o-k-body"></span> postOKBody

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| terraform | [][PostOKBodyTerraformItems](#post-o-k-body-terraform-items)| `[]*PostOKBodyTerraformItems` |  | |  |  |


#### <span id="post-o-k-body-terraform-items"></span> postOKBodyTerraformItems

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| result | [interface{}](#interface)| `interface{}` | ✓ | |  |  |
| success | boolean| `bool` | ✓ | |  |  |


#### <span id="post-params-body"></span> postParamsBody

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| commands | [][PostParamsBodyCommandsItems](#post-params-body-commands-items)| `[]*PostParamsBodyCommandsItems` |  | | Array of commands. |  |
| loglevel | string| `string` |  | `"off"`| Terraform log level, default off |  |
| scope | string| `string` |  | `"instance"`| Scope where the log file is stored, default instance. Filename `tf.log`. |  |
| variables | [][PostParamsBodyVariablesItems](#post-params-body-variables-items)| `[]*PostParamsBodyVariablesItems` |  | | Variables set for all commands. This translatyes into TF_VAR_* environment variables. | `[{"name":"instance_name","value":"myinstance"}]` |


#### <span id="post-params-body-commands-items"></span> postParamsBodyCommandsItems

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| command | string| `string` |  | | Command to run | `terraform version` |
| continue | boolean| `bool` |  | | Stops excecution if command fails, otherwise proceeds with next command |  |
| print | boolean| `bool` |  | `true`| If set to false the command will not print the full command with arguments to logs. |  |
| silent | boolean| `bool` |  | | If set to false the command will not print output to logs. |  |


#### <span id="post-params-body-variables-items"></span> postParamsBodyVariablesItems

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| name | string| `string` |  | | Name of the variable. |  |
| value | string| `string` |  | | Value of the variable. |  |

 
