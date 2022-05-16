
# terraform 1.0.0

Run Hashicorp Terrafrom from Direktiv

---
- #### Categories: Cloud, Tools
- #### Image: direktiv/terraform 
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
    image: direktiv/terraform
    type: knative-workflow
  ```
   #### Basic
   ```yaml
   - id: tf
     type: action
     action:
      function: get
      files:
      # Contains all required .tf files. Can point to a plain text .tf file as well.
      - scope: workflow
        key: tfbase.tar.gz
        as: tf
        type: tar.gz
      input: 
        commands:
        # the execution dir (chdir) is "tf" which we create in the "files" section
        # Storing the state in "../out/workflow/terraform.tfstate" will store the state in workflow scope. 
        - terraform -chdir=tf apply -state=../out/workflow/terraform.tfstate -no-color -auto-approve
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
    "format_version": 1,
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
| output | [][PostOKBodyOutputItems](#post-o-k-body-output-items)| `[]*PostOKBodyOutputItems` |  | |  |  |


#### <span id="post-o-k-body-output-items"></span> postOKBodyOutputItems

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| result | [interface{}](#interface)| `interface{}` | ✓ | |  |  |
| success | boolean| `bool` | ✓ | |  |  |


#### <span id="post-params-body"></span> postParamsBody

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| commands | []string| `[]string` |  | | Commands to execute in order. | `["terraform -chdir=out/workflow/tfbase.tar.gz plan"]` |
| continue | boolean| `bool` |  | | If set to true all commands are getting executed and errors ignored. | `true` |
| variables | [][PostParamsBodyVariablesItems](#post-params-body-variables-items)| `[]*PostParamsBodyVariablesItems` |  | | Variables set for all commands. This translatyes into TF_VAR_* environment variables. | `[{"name":"instance_name","value":"myinstance"}]` |


#### <span id="post-params-body-variables-items"></span> postParamsBodyVariablesItems

  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| name | string| `string` |  | | Name of the variable. |  |
| value | string| `string` |  | | Value of the variable. |  |

 
