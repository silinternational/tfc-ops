# Terraform Cloud Ops Tool
This application can be helpful in making copies/clones of a workspace and bringing its variables over
to the new one. It can also be used for listing workspaces and listing or modifying variables in workspaces.

## Required ENV vars
- `ATLAS_TOKEN` - Must be set as an environment variable. Get this by going to 
https://app.terraform.io/app/settings/tokens and generating a new token.
- `ATLAS_TOKEN_DESTINATION` - Only necessary if cloning to a new organization in TF Cloud.


## Installation
There are three ways to download/install this script:

1. Download a pre-built binary for your operating system by going to [/dist](https://github.com/silinternational/tfc-ops/tree/master/dist)
2. If you're a Go developer you can install it by running `go get -u https://github.com/silinternational/tfc-ops.git`
3. If you're a Go developer and want to modify the source before running, clone this repo and run with `go run main.go ...`

## Cloning a TF Cloud Workspace
Examples.

Get help about the command.

```$ go run main.go workspaces clone -h```

Clone just the workspace (no variables) to the same organization.

```$ go run main.go workspaces clone -o=my-org -s=source-workspace -n=new-workspace```

Clone a workspace, its variables and its state to a different organization in TF Cloud.

Note: Sensitive variables will get a placeholder in the new workspace, the value of
which will need to be corrected manually.  Also, the environment variables from the source
workspace will need to be moved in the destination workspace from the normal variables section
down to the environment variables section (e.g. `CONFIRM_DESTROY`).

```
$ go run main.go workspaces clone -c=true -t=true -o=org1 -p=org2 -d=true \
$   -s=source-workspace -n=destination-workspace -v=org2-vcs-token
```


## Getting a list of all TF Cloud Workspaces with some of their attributes 
Examples.

Get help about the command.

```$ go run main.go workspaces list -h```

List the workspaces with at least one of their attributes/pieces of data.

```$ go run main.go workspaces list -o=gtis -a=id,name,createdat,environment,workingdirectory,terraformversion,vcsrepo```

## Usage

### General Help
```text
$ tfc-ops -h
Perform TF Cloud operations, e.g. clone a workspace or manage variables within a workspace

Usage:
  tfc-ops [command]

Available Commands:
  help        Help about any command
  variables   Operations and reports related to on variables
  workspaces  Operations and reports related to on workspaces

Flags:
  -h, --help   help for tfc-ops

Use "tfc-ops [command] --help" for more information about a command.
```

### Workspaces Help
```text
Top Level Command for the following sub-commands

Usage:
  tfc-ops workspaces [command]

Available Commands:
  clone       Clone a V2 Workspace
  list        List Workspaces

Flags:
  -h, --help   help for workspaces
```


### Clone Help
```text
$ tfc-ops workspaces clone -h
Clone a TF Enterprise Version 2 Workspace

Usage:
  tfc-ops clone [flags]

Flags:
  -t, --copyState                     [optional] (e.g. "-t=true") whether to copy the state of the Source Workspace (only possible if copying to a new account).
  -c, --copyVariables                 [optional] (e.g. "-c=true") whether to copy the values of the Source Workspace variables.
  -d, --differentDestinationAccount   [optional] (e.g. "-d=true") whether to clone to a different TF account.
  -h, --help                          help for clone
  -p, --new-organization string       Name of the Destination Organization in TF Enterprise (version 2)
  -v, --new-vcs-token-id string       The new Organization's VCS repo's oauth-token-id
  -n, --new-workspace string          Name of the new Workspace in TF Enterprise (version 2)
  -o, --organization string           Name of the Organization in TF Enterprise (version 2)
  -s, --source-workspace string       Name of the Source Workspace in TF Enterprise (version 2)
```

### List Help
```text
$ tfc-ops workspaces list -h
Lists the TF workspaces with (some of) their attributes

Usage:
  tfc-ops workspaces list [flags]

Flags:
  -a, --attributes string     [required] - Workspace attributes to list: id,name,createdat,environment,workingdirectory,terraformversion,vcsrepo
  -h, --help                  help for list
  -o, --organization string   [required] - Name of Terraform Enterprise Organization
```

### Variables Help
```text
Top Level Command for the following sub-commands

Usage:
  tfc-ops variables [command]

Available Commands:
  list        Report on variables
  update      Update/add a variable in a V2 Workspace

Flags:
  -h, --help   help for variables
```

### List Variables
```text
$ tfc-ops variables -h

Show the values of variables with a key or value containing a certain string

Usage:
  tfc-ops variables list [flags]

Flags:
  -h, --help                    help for variables
  -k, --key_contains string     [required] if value_contains is blank - string contained in the Terraform variable keys to report on
  -o, --organization string     [required] - Name of Terraform Enterprise Organization
  -v, --value_contains string   [required] if key_contains is blank - string contained in the Terraform variable values to report on
  -w, --workspace string        [optional] Name of the Workspace in TF Enterprise
```

### Update Help
```text
$ tfc-ops variables update -h
Update or add a variable in a TF Enterprise Version 2 Workspace based on a complete case-insensitive match

Usage:
  tfc-ops variables update [flags]

Flags:
  -a, --add-key-if-not-found            [optional] (e.g. "-a=true") whether to add a new variable if a matching key is not found.
  -d, --dry-run-mode                    [optional] (e.g. "-d=true") dry run mode only.
  -h, --help                            help for update
  -n, --new-variable-value string       The desired new value of the variable
  -o, --organization string             Name of the Organization in TF Enterprise (version 2)
  -v, --search-on-variable-value        [optional] (e.g. "-v=true") whether to do the search based on the value of the variables. (Must be false if add-key-if-not-found is true
  -s, --variable-search-string string   The string to match in the current variables (either in the Key or Value - see other flags)
  -w, --workspace string                [optional] Name of the Workspace in TF Enterprise (version 2)
```

## License
terraform-enterprise-monitor is released under the Apache 2.0 license. See 
[LICENSE](https://github.com/silinternational/terraform-enterprise-monitor/blob/master/LICENSE)
