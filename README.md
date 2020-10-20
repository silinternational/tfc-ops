# Terraform Enterprise Migration Tool
This application was initially used to migrate Terraform Enterprise (TFE) Legacy environments to new 
Terraform Enterprise workspaces.  However, that is not needed anymore. Instead, this application 
can be helpful in making copies/clones of a workspace and bringing its variables over to the new one.

The original migration process is a 3-4 step process using this application, but it was much faster than 
manually copying everything over using the Terraform Enterprise web interface. 

## Disclaimer
While we were able to use this application and process to migrate ~75 environments successfully, we obviously cannot
guarantee it will work for you or that it'll be flawless. So with that said, start with a small migration of a test
environment and then work up before you batch migrate everything so that you can be sure it is working properly.

## Required ENV vars
- `ATLAS_TOKEN` - Must be set as an environment variable. Get this by going to 
https://app.terraform.io/app/settings/tokens and generating a new token.
- `ATLAS_TOKEN_DESTINATION` - Only necessary if cloning to a new organization in TF Cloud.


## Installation
There are three ways to download/install this script:

1. Download a pre-built binary for your operating system by going to [/dist](https://github.com/silinternational/terraform-enterprise-migrator/tree/master/dist)
2. If you're a Go developer you can install it by running `go get -u https://github.com/silinternationa/terraform-enterprise-migrator`
3. If you're a Go developer and want to modify the source before running, clone this repo and run with `go run main.go ...`

## Cloning a TF Cloud Workspace
Examples.

Get help about the command.

```$ go run main.go clone -h```

Clone just the workspace (no variables) to the same organization.

```$ go run main.go clone -o=my-org -s=source-workspace -n=new-workspace```

Clone a workspace, its variables and its state to a different organization in TF Cloud.

Note: Sensitive variables will get a placeholder in the new workspace, the value of
which will need to be corrected manually.  Also, the environment variables from the source
workspace will need to be moved in the destination workspace from the normal variables section
down to the environment variables section (e.g. `CONFIRM_DESTROY`).

```
$ go run main.go clone -c=true -t=true -o=org1 -p=org2 -d=true \
$   -s=source-workspace -n=destination-workspace -v=org2-vcs-token
```


## Getting a list of all TF Cloud Workspaces with some of their attributes 
Examples.

Get help about the command.

```$ go run main.go list -h```

List the workspaces with at least one of their attributes/pieces of data.

```$ go run main.go list -o=gtis -a=id,name,createdat,environment,workingdirectory,terraformversion,vcsrepo```


## Original Migration Process
As mentioned above the migration process is at least 3 steps but may be 4 if you mark any variables as `sensitive` 
in your configurations. The first step is to generate a plan file. This file is a CSV file that you'll need to fill 
in some missing columns of data for. Unfortunately the v1 TFE APIs do not provide all information required to create 
new workspaces from them. Specifically it does not include VCS repository information including the repo url, 
branch, and working directory. It also does not provide the version of Terraform used. So this script generates a 
CSV file with the information it does know, you update the file with the information it needs, and then run a 
migration from the plan. Finally if you have any variables marked as sensitive they cannot be read so they need to be 
updated manually in the TFE web interface. After the migration completes a list of environments that were migrated 
will be printed out along with any sensitive variables for each so you know which will need to be updated.

1. Generate `plan.csv`: `terraform-enterprise-migrator plan -l legacyOrg -n newOrg`
2. Edit `plan.csv` however you want and fill in all empty columns (`TerraformVersion`,`RepoID`,`Branch`,`Directory`)
3. Run migration: `terraform-enterprise-migrator migrate -u vcsUsername`. Notice the `-u` flag here, the value should be
the username for the version control system you are connected with, for example your Github username. TFE uses the 
existing VCS connection you have when creating the new workspace and attaching a webhook to the VCS repo.

## Example
Here is a video of us using this script to migrate a dozen environments from legacy to new. You can see the script 
running on the left and the Terraform Enterprise interface on the right to see the new workspaces just appearing as the 
script runs. You can see how much time this script saves when we move 12 environments in 3 minutes. Manually it would 
have taken closer to 3 hours!

![Example video](video.gif)

## Usage

### General Help
```text
$ terraform-enterprise-migrator -h
Migration is a three step process. First run plan, modify the generated file as needed,
then run migrate to process the plan file

Usage:
  terraform-enterprise-migrator [command]

Available Commands:
  clone       Clone a V2 Workspace
  help        Help about any command
  list        List Workspaces
  migrate     Perform migration plan
  plan        Generate migration plan file
  update      Update/add a variable in a V2 Workspace
  variables   Report on variables

Flags:
  -h, --help   help for terraform-enterprise-migrator

Use "terraform-enterprise-migrator [command] --help" for more information about a command.
```

### Plan Help
```text
$ terraform-enterprise-migrator plan -h
Generates a plan.csv file with list of environments from legacy organization
for mapping to new organization.

Usage:
  terraform-enterprise-migrator plan [flags]

Flags:
  -f, --file string     optional - Name of migration plan CSV file (default "plan.csv")
  -h, --help            help for plan
  -l, --legacy string   required - Name of Terraform Enterprise Legacy Organization
  -n, --new string      required - Name of new Terraform Enterprise Organization
```

### Migrate Help
```text
$ terraform-enterprise-migrator migrate -h
Processes plan.csv to validate migration plan and perform the work

Usage:
  terraform-enterprise-migrator migrate [flags]

Flags:
  -f, --file string           optional - Name of migration plan CSV file (default "plan.csv")
  -h, --help                  help for migrate
  -u, --vcs-username string   Name of the VCS User in TF Enterprise (new version) to allow us to get the right VCS Token ID (the GitHub or Bitbucket username used to connect TFE)
```

### Clone Help
```text
$ terraform-enterprise-migrator clone -h
Clone a TF Enterprise Version 2 Workspace

Usage:
  terraform-enterprise-migrator clone [flags]

Flags:
  -t, --copyState                     optional (e.g. "-t=true") whether to copy the state of the Source Workspace (only possible if copying to a new account).
  -c, --copyVariables                 optional (e.g. "-c=true") whether to copy the values of the Source Workspace variables.
  -d, --differentDestinationAccount   optional (e.g. "-d=true") whether to clone to a different TF account.
  -h, --help                          help for clone
  -p, --new-organization string       Name of the Destination Organization in TF Enterprise (version 2)
  -v, --new-vcs-token-id string       The new Organization's VCS repo's oauth-token-id
  -n, --new-workspace string          Name of the new Workspace in TF Enterprise (version 2)
  -o, --organization string           Name of the Organization in TF Enterprise (version 2)
  -s, --source-workspace string       Name of the Source Workspace in TF Enterprise (version 2)
```

### List Help
```text
$ terraform-enterprise-migrator list -h
Lists the TF workspaces with (some of) their attributes

Usage:
  terraform-enterprise-migrator list [flags]

Flags:
  -a, --attributes string     required - Workspace attributes to list: id,name,createdat,environment,workingdirectory,terraformversion,vcsrepo
  -h, --help                  help for list
  -o, --organization string   required - Name of Terraform Enterprise Organization
```


### Update Help
```text
$ terraform-enterprise-migrator update -h
Update or add a variable in a TF Enterprise Version 2 Workspace based on a complete case-insensitive match

Usage:
  terraform-enterprise-migrator update [flags]

Flags:
  -a, --add-key-if-not-found            optional (e.g. "-a=true") whether to add a new variable if a matching key is not found.
  -d, --dry-run-mode                    optional (e.g. "-d=true") dry run mode only.
  -h, --help                            help for update
  -n, --new-variable-value string       The desired new value of the variable
  -o, --organization string             Name of the Organization in TF Enterprise (version 2)
  -v, --search-on-variable-value        optional (e.g. "-v=true") whether to do the search based on the value of the variables. (Must be false if add-key-if-not-found is true
  -s, --variable-search-string string   The string to match in the current variables (either in the Key or Value - see other flags)
  -w, --workspace string                Name of the Workspace in TF Enterprise (version 2)
```

### Variables Help
```text
$ terraform-enterprise-migrator variables -h
Show the values of variables with a key or value containing a certain string

Usage:
  terraform-enterprise-migrator variables [flags]

Flags:
  -h, --help                    help for variables
  -k, --key_contains string     required if value_contains is blank - string contained in the Terraform variable keys to report on
  -o, --organization string     required - Name of Terraform Enterprise Organization
  -v, --value_contains string   required if key_contains is blank - string contained in the Terraform variable values to report on
```

## License
terraform-enterprise-monitor is released under the Apache 2.0 license. See 
[LICENSE](https://github.com/silinternational/terraform-enterprise-monitor/blob/master/LICENSE)