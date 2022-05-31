Feature: greeting end-point

Background:
* url demoBaseUrl

Scenario: version

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    { 
        "commands": [
            {
                "command": "terraform version"
            }
        ]
    }
    """
    When method post
    Then status 200
    And match $ == 
    """
    {
    "terraform": [
    {
      "result": "#notnull",
      "success": true
    }
    ]
    }
    """


Scenario: logging

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "scope": "workflow",
        "loglevel": "trace",
        "commands": [
            {
                "command": "env"
            }
        ]
    }
    """
    When method post
    Then status 200
    * def result = $response.terraform[0].result
    * print result
    * match result contains("TF_LOG=trace")
    * match result contains("TF_LOG_PATH=/tmp/out/workflow/tf.log")

Scenario: vars

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "variables": [
            {
                "name": "hello",
                "value": "world"
            },
            {
                "name": "hello1",
                "value": "world1"
            }
        ],
        "commands": [
            {
                "command": "env"
            }
        ]
    }
    """
    When method post
    Then status 200
    * def result = $response.terraform[0].result
    * print result
    * match result contains("TF_VAR_hello=world")
    * match result contains("TF_VAR_hello1=world1")


Scenario: simple run

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "variables": [
            {
                "name": "name",
                "value": "MyName"
            }
        ],
        "commands": [
            {
                "command": "mkdir -p /tmp/out/instance"
            },
            {
                "command": "terraform apply -auto-approve"
            },
            {
                "command": "terraform output -json"
            }
        ]
    }
    """
    When method post
    Then status 200
    * def result = $response.terraform[2].result.output_hello.value
    * match result == "MyName"


Scenario: envs, comma test

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "envs": [
            {
                "name": "AWS_ENVIRONMENT",
                "value": "123"
            }
        ],
        "variables": [
            {
                "name": "name",
                "value": "MyName"
            }
        ],
        "commands": [
            {
                "command": "env"
            }
        ]
    }
    """
    When method post
    Then status 200

Scenario: envs, comma test2

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "variables": [
            {
                "name": "name",
                "value": "MyName"
            }
        ],
        "commands": [
            {
                "command": "env"
            }
        ]
    }
    """
    When method post
    Then status 200

Scenario: envs, comma test3

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "envs": [
            {
                "name": "AWS_ENVIRONMENT",
                "value": "123"
            }
        ],
        "commands": [
            {
                "command": "env"
            }
        ]
    }
    """
    When method post
    Then status 200

Scenario: envs, comma test4

    Given path '/'
    Given header Direktiv-ActionID = 'development'
    Given header Direktiv-Tempdir = '/tmp'
    And request 
    """
    {   
        "commands": [
            {
                "command": "env"
            }
        ]
    }
    """
    When method post
    Then status 200