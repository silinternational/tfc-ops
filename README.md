# Terraform Enterprise Migration Tool
This application is used to migrate Terraform Enterprise (TFE) Legacy environments to new Terraform Enterprise workspaces.
HashiCorp rebuilt the Terraform Enterprise experience and is ending support for the legacy version on March 30, 2018. 
We as an organization have a LOT of Terraform environments ~75 at this time, so moving our environments from the
legacy version to the new version manually wasn't a reasonable task. Thankfully both legacy and new services have 
APIs which were able to provide the majority of information needed to automate the migration.

The migration process is a 3-4 step process using this application, but it is much faster than manually copying 
everything over using the Terraform Enterprise web interface. 

## Disclaimer
While we were able to use this application and process to migrate ~75 environments successfully, we obviously cannot
guarantee it will work for you or that it'll be flawless. So with that said, start with a small migration of a test
environment and then work up before you batch migrate everything so that you can be sure it is working properly.

## Required ENV vars
- `ATLAS_TOKEN` - Must be set as an environment variable. Get this by going to 
https://app.terraform.io/app/settings/tokens and generating a new token.

## Installation
There are three ways to download/install this script:

1. Download a pre-built binary for your operating system by going to [/dist](https://github.com/silinternational/terraform-enterprise-migrator/tree/master/dist)
2. If you're a Go developer you can install it by running `go get -u https://github.com/silinternationa/terraform-enterprise-migrator`
3. If you're a Go developer and want to modify the source before running, clone this repo and run with `go run main.go ...`

## Migration Process
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
  help        Help about any command
  migrate     Perform migration plan
  plan        Generate migration plan file

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
  -c, --copyVariables             optional (e.g. "-c=true") whether to copy the values of the Source Workspace variables.
  -h, --help                      help for clone
  -n, --new-workspace string      Name of the New Workspace in TF Enterprise (version 2)
  -o, --organization string       Name of the Organization in TF Enterprise (version 2)
  -s, --source-workspace string   Name of the Source Workspace in TF Enterprise (version 2)
```

## License
terraform-enterprise-monitor is released under the Apache 2.0 license. See 
[LICENSE](https://github.com/silinternational/terraform-enterprise-monitor/blob/master/LICENSE)