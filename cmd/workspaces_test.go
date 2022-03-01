package cmd

var fake = `{
	"data": [
		{
			"id": "ws-z55zWHsNCbuXUDRb",
			"type": "workspaces",
			"attributes": {
				"allow-destroy-plan": true,
				"auto-apply": false,
				"auto-destroy-at": null,
				"created-at": "2019-03-04T15:29:40.193Z",
				"environment": "default",
				"locked": false,
				"name": "jira-insite-stg",
				"queue-all-runs": false,
				"speculative-enabled": true,
				"structured-run-output-enabled": false,
				"terraform-version": "1.1.4",
				"working-directory": "jira",
				"global-remote-state": true,
				"updated-at": "2022-03-01T15:22:47.343Z",
				"resource-count": 66,
				"apply-duration-average": 22000,
				"plan-duration-average": 46000,
				"policy-check-failures": null,
				"run-failures": 8,
				"workspace-kpis-runs-count": 30,
				"latest-change-at": "2022-03-01T15:22:44.794Z",
				"operations": true,
				"execution-mode": "remote",
				"vcs-repo": {
					"branch": "develop",
					"ingress-submodules": false,
					"identifier": "silintl/jira-terraform",
					"display-identifier": "silintl/jira-terraform",
					"oauth-token-id": "ot-nGSDjSQYkGYwBs28",
					"webhook-url": "https://app.terraform.io/webhooks/vcs/83e198f0-dc6e-4f84-9420-388a29eb6eef",
					"repository-http-url": "https://bitbucket.org/silintl/jira-terraform",
					"service-provider": "bitbucket_hosted"
				},
				"vcs-repo-identifier": "silintl/jira-terraform",
				"permissions": {
					"can-update": true,
					"can-destroy": true,
					"can-queue-destroy": true,
					"can-queue-run": true,
					"can-queue-apply": true,
					"can-read-state-versions": true,
					"can-create-state-versions": true,
					"can-read-variable": true,
					"can-update-variable": true,
					"can-lock": true,
					"can-unlock": true,
					"can-force-unlock": true,
					"can-read-settings": true,
					"can-manage-tags": true
				},
				"actions": {
					"is-destroyable": true
				},
				"description": null,
				"file-triggers-enabled": false,
				"trigger-prefixes": [],
				"source": null,
				"source-name": null,
				"source-url": null,
				"tag-names": []
			},
			"relationships": {
				"organization": {
					"data": {
						"id": "gtis",
						"type": "organizations"
					}
				},
				"current-run": {
					"data": {
						"id": "run-p1c5Ca2edVjw3yvm",
						"type": "runs"
					},
					"links": {
						"related": "/api/v2/runs/run-p1c5Ca2edVjw3yvm"
					}
				},
				"latest-run": {
					"data": {
						"id": "run-p1c5Ca2edVjw3yvm",
						"type": "runs"
					},
					"links": {
						"related": "/api/v2/runs/run-p1c5Ca2edVjw3yvm"
					}
				},
				"outputs": {
					"data": [
						{
							"id": "wsout-2upQMg7qyRbtNiue",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-QTjLvj3KuUBJ4Y5W",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-de9ySBYKw3FHeX23",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-iCDQpo2gdamho8cb",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-83oJK66W3zgc3kay",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-PU5jmpLvhE9Xman8",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-64wQrPFFyCLN5CTj",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-fozdtcJPEeRaTh2b",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-XbYTfkUpYPrBNFhH",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-9KRsk7zeEkssdaEv",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-5hMJtENzzfNJ9wbo",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-83hJnkY9Rsn7XWgJ",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-ozYcM4QmidtDA5KS",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-7zVKronjVg9STN6Y",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-wVLNjuGawkvP1H61",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-EGk7DmYEmNvXyunx",
							"type": "workspace-outputs"
						},
						{
							"id": "wsout-E3pQkFsmYucKxfVg",
							"type": "workspace-outputs"
						}
					]
				},
				"remote-state-consumers": {
					"links": {
						"related": "/api/v2/workspaces/ws-z55zWHsNCbuXUDRb/relationships/remote-state-consumers"
					}
				},
				"current-state-version": {
					"data": {
						"id": "sv-PPYSHYdqzrpUUNNr",
						"type": "state-versions"
					},
					"links": {
						"related": "/api/v2/workspaces/ws-z55zWHsNCbuXUDRb/current-state-version"
					}
				},
				"current-configuration-version": {
					"data": {
						"id": "cv-Atanc2UPVoHqHdGB",
						"type": "configuration-versions"
					},
					"links": {
						"related": "/api/v2/configuration-versions/cv-Atanc2UPVoHqHdGB"
					}
				},
				"agent-pool": {
					"data": null
				},
				"readme": {
					"data": {
						"id": "413731",
						"type": "workspace-readme"
					}
				}
			},
			"links": {
				"self": "/api/v2/organizations/gtis/workspaces/jira-insite-stg"
			}
		}
	],
	"links": {
		"self": "https://app.terraform.io/api/v2/organizations/gtis/workspaces?page%5Bnumber%5D=1&page%5Bsize%5D=1",
		"first": "https://app.terraform.io/api/v2/organizations/gtis/workspaces?page%5Bnumber%5D=1&page%5Bsize%5D=1",
		"prev": null,
		"next": "https://app.terraform.io/api/v2/organizations/gtis/workspaces?page%5Bnumber%5D=2&page%5Bsize%5D=1",
		"last": "https://app.terraform.io/api/v2/organizations/gtis/workspaces?page%5Bnumber%5D=156&page%5Bsize%5D=1"
	},
	"meta": {
		"status-counts": {
			"pending": 1,
			"plan-queued": 0,
			"planning": 0,
			"planned": 0,
			"confirmed": 0,
			"apply-queued": 0,
			"applying": 0,
			"applied": 63,
			"discarded": 5,
			"errored": 4,
			"canceled": 0,
			"cost-estimating": 0,
			"cost-estimated": 0,
			"policy-checking": 0,
			"policy-override": 0,
			"policy-checked": 0,
			"policy-soft-failed": 0,
			"planned-and-finished": 83,
			"post-plan-running": 0,
			"post-plan-completed": 0,
			"pre-apply-running": 0,
			"pre-apply-completed": 0,
			"fetching": 0,
			"none": 0,
			"total": 156
		},
		"pagination": {
			"current-page": 1,
			"page-size": 1,
			"prev-page": null,
			"next-page": 2,
			"total-pages": 156,
			"total-count": 156
		}
	}
}`
