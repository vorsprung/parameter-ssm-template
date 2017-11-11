SSM Parameter Store legacy
--------------------------

This go program and scripts allow the management of variables in config files, programs etc from the AWS SSM Parameter Store

## Why use the Parameter store?

For convenience reasons: AWS provide a console for editing values, a
security framework for restricting access, encryption keys and a rest
api to the key/values

## Why go?

I manage servers with legacy environments: many of them don't have the
latest aws SDK or cli installed.  The great thing about go is that it
is easy to make a static binary version.  Provided that the host that
is used for the compile and the target hosts are binary compatible,
this one file is all that is need

In the context of preparing configuration this is a great advantage!

See below for static compiler details

## What's here

In Scripts

   - `deletewithcli.bash`  shell script to delete from store
   - `loader.go`   loads key/values from flat files
   - `template.go`  fills in templates
   - `kvtest.templ`  example template

In src
   - `sfill.go`    main program
   - `store_test.go`    unit tests for main program
   note that the tests require AWS keys to be available to connect to
   a AWS parameter store
   - `_kvtextdata.txt`

## How to compile to a static binary

To compile a static binary, I used this command line for the template program
`CGO_ENABLED=0 GOOS=linux go build -o template -a -ldflags '-extldflags "-static"' template.go`

## How to install the aws dependencies
change into the src/sfill directory and issue this command
    `go get ./...`

## How to restrict access to a particular key set 
The policy below can be attached to a user or role to allow access all
the keys under a /test tree.  Replace the account id with your account id

There are details on how to do this and use other methods such as tags here
http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-access.html



    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": [
                    "ssm:DescribeParameters"
                ],
                "Resource": "*"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "ssm:GetParameters",
                    "ssm:GetParameterHistory",
                    "ssm:GetParameter",
                    "ssm:GetParameters",
                    "ssm:GetParametersByPath"
                ],
                "Resource": [
                    "arn:aws:ssm:eu-west-1:ACCOUNTID:parameter/test",
                    "arn:aws:ssm:eu-west-1:ACCOUNTID:parameter/test/*"
                ]
            }
        ]
    }

## Gotchas and restrictions
The keys must match the AWS specification of alphanumeric, -/._ see the
regexp in the Parameterstorestringvalidate function

The template name always ends with a .templ extension.  This extension is
not required by the template script
