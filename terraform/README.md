# Example Terraform Configs

You will want to build the binaries first before running [terraform](https://www.terraform.io/)

The IAM policies are rather minimal and should only be providing the minimal
required access to the resources.

The instance profile is probably the most artificial and will need tweaking
as it only includes permissions to invoke `snedd-initiator`.

 * create IAM roles and policies for lambda functions and sfn state machines
 * attach policies to the correct roles
 * create snedd sfn state machine with it's own role
 * create snedd-initiator lambda function with it's own role
 * create snedd-expirer lambda function with it's own role
 * create IAM instance profile allowed to invoke snedd-initiator

No instances will be created.

## Applying the Config

I use [`aws-vault`](https://github.com/99designs/aws-vault) to manage my
credentials. Applying the terraform config will look something like:

```
$ aws-vault exec $role -- terraform plan
$ aws-vault exec $role -- terraform apply
```
