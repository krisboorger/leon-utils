#!/bin/sh
set -euo pipefail

terraform -chdir=terraform destroy -auto-approve
