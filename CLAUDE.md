When adding new tools that make API calls:
- Use built in functionality in the cloud-sdk
- Refer to the Terraform provider - often very close to what we're implementing https://github.com/temporalio/terraform-provider-temporalcloud
- Update the readme to describe the new tool
- Update test_mcp.sh with the new tool