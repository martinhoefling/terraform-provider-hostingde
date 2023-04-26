# Configuration-based authentication
# This is not recommended, use environment variables to configure the provider:
# HOSTINGDE_AUTH_TOKEN, (Optional: HOSTINGDE_ACCOUNT_ID)
# Can be omitted when using the environment variables

provider "hostingde" {
  auth_token = "YOUR_API_TOKEN"
  account_id = "YOUR_ACCOUNT_ID"
}
