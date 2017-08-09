[![CircleCI](https://circleci.com/gh/uchimanajet7/ca-cli.svg?style=svg)](https://circleci.com/gh/uchimanajet7/ca-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/uchimanajet7/ca-cli)](https://goreportcard.com/report/github.com/uchimanajet7/ca-cli)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/uchimanajet7/ca-cli/blob/master/LICENSE)
 
# ca-cli
The `ca` is command line tool to invoke `Cloud Automator` API


## Description
The `ca` is the initial letter of the `Cloud Automator`.

The `ca` is command line tool. You can easily use the `Cloud Automator` API. The execution result of the API is displayed on the command line as JSON.

***see also:***

- "Cloud Automator" - Change your way of cloud operation
	- https://cloudautomator.com/en/
- Cloud Automator API
	- https://cloudautomator.com/api_docs/v1/api.html

## Demo
![ca_cli_demo](https://user-images.githubusercontent.com/6448792/29017893-d140d520-7b93-11e7-98d0-6df767fc3643.gif)

## Features
- Works on Mac / Linux / Windows by binary file cross - compiled with Go.
- Since it works by just copying a single binary file, you do not have to worry about environment building and dependency.
- Sends http request to the `Cloud Automator` API using the specified argument. 
- Response from API Output JSON directly to standard output.
- You can easily process the output of the command and pass it to another command or tool.

## Requirement
- Go 1.8+
	- It is not necessary when using the released binary.
	- https://github.com/uchimanajet7/ca-cli/releases
- `API Key` created in your account of `Cloud Automator`.
	- For details on how to create API key, see url below.
	- http://blog.serverworks.co.jp/tech/2017/07/07/releasing-cloudautomator-rest-api/

## Usage
First of all, create a profile using the `configure` command.

```	console
$ ./ca configure
```
`API Key` and `Endpoint` can be registered by executing the command.

```console
$ ./ca configure
Register the information necessary for execution as a profile of "default".

API Key:
API Key: 123**************************abc
Endpoint:
```
`API Key` is **required**, `Endpoint` is set to change to default.

The `ca` command invoke the API using the registered `API Key`.

### Multiple profiles

You can use it when you have multiple `Cloud Automator` accounts, or want to access multiple `Endpoint`.

Add the `--profile` flag to the command and specify the profile name.

You can create profiles named user1 and user2 by doing as follows.

```console
$ ./ca configure --profile user1
  :
  (Enter information about user1)

$ ./ca configure --profile user2
  :
  (Enter information about user2)
```
To use the registered profile, please use the `--profile` flag.

### More information

Please see the help.

```	console
$ ./ca help
A command line tool to invoke Cloud Automator(CA) API

Usage:
  ca [command]

Available Commands:
  aws-account Manage CA AWS accounts
  configure   Manage ca-cli profiles
  help        Help about any command
  job         Manage CA jobs
  log         Manage CA job logs
  version     Print the version number of ca-cli

Flags:
  -h, --help             help for ca
  -p, --profile string   Specify profile name (default "default")

Use "ca [command] --help" for more information about a command.
```

## Installation

Please select the package file for your own environment from the releases page, download and unpack it, and put the executable file in a place where included in PATH.

- Releases · uchimanajet7/ca-cli
	- https://github.com/uchimanajet7/ca-cli/releases

If you build from source yourself.

```	console
$ go get github.com/uchimanajet7/ca-cli
$ cd $GOPATH/src/github.com/uchimanajet7/ca-cli
$ make
```

## Object type parameters

It is necessary to specify the `object type` as parameter when ca `job create` command.

An execution example is shown below.

```	console
$ ./ca job create \
--name "ca job create exsample" \
--aws-account-id 1 \
--rule-type cron \
--rule-value hour=2,minutes=0,schedule_type=weekly,weekly_schedule=monday,friday \
--action-type create_image \
--action-value region=ap-northeast-1\
,specify_image_instance=identifier\
,instance_id=i-xxxxxxxxxxxxxxxxx\
,generation=1\
,image_name=exsample-ami\
,description="Job Create Exsample Cloud Automator CLI"\
,reboot_instance=true\
,additional_tag_key=name\
,additional_tag_value=exsample\
,add_same_tag_to_snapshot=true\
,trace_status=true\
,recreate_image_if_ami_status_failed=true
```

`--rule-value` and `--action-value` flag are `object type`.
**API parameter name** `=` **value** format is connected with `,` and passed to the flag. If you want to express **array** by **value**, you can connect values with `,`.

***see also:***

- Cloud Automator API
	- https://cloudautomator.com/api_docs/v1/api.html#ジョブ-post

## Useful JSON tools
Besides the famous JSON tool `jq`, there are useful ones implemented by golang.

- jq
	- https://stedolan.github.io/jq/

Especially `jid` is convenient because I can drill down into JSON interactively.

- simeji/jid: json incremental digger
	- https://github.com/simeji/jid
- JSONをインタラクティブに掘り下げるコマンド jid - Qiita
	- http://qiita.com/simeji/items/dd0464b7ed91c51ee618

## Author
[uchimanajet7](https://github.com/uchimanajet7)


## Licence
[MIT](https://github.com/uchimanajet7/ca-cli/blob/master/LICENSE)
