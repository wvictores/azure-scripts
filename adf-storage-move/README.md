# ADF Storage Move

Terraform script to create Data Factory to move files from one blob container to another blob container.

Example Commands:
```
terraform init
terraform plan -var="resourceGroupName=move-test" -var="location=westus2" -var="factoryName=move-test-df" -var="sourceStorageAccountName=move0test0source0stor" -var="sourceContainerName=stor" -var="destStorageAccountName=move0test0source0stor" -var="destContainerName=dest" -out plan
terraform apply "plan"
```