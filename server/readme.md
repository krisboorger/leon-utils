# Server
## Purpose of this component
The goal of the server is to collect answers received from clients and return them the combined answer counts.

## How to deploy it?
1. Create a cert - more on that in the readme of `cert/` subdirectory.
2. Deploy the server
Configs necessary to deploy this server to a K8s cluster on Azure are already included.
If you are logged in with `az cli` and have terraform, kubectl and helm installed locally, simply run:
```
./deploy/deploy_aks.sh
```
directly from this dir.

3. After the script finishes, add the IP address printed to your DNS A record, as instructed.

4. When the server is no longer needed, in order to destroy all the resources, run:
```
./deploy/destroy_aks.sh
```
This should get rid of any resources left.
