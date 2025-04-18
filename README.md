<h1 align="center" style="border-bottom: none">
    <a href="https://nimtechnology.com/2023/07/02/zeus-retention-project/" target="_blank"><img alt="Zeus Rotation" width="120px" src="https://nimtechnology.com/wp-content/uploads/2023/07/2185568.png"></a><br>Zeus Rotations
</h1>

<p align="center">Auto Rotations AWS key on Kubernetes</p>

## Instroduction

Zeus is the platform to help you to rotate the AWS Access and Secret Key, which are saved in secret on EKS


## How to Use Zeus Rotations.

**First Step:** Using [EKS IRSA Terraform Module](https://registry.terraform.io/modules/mrnim94/eks-irsa/aws/latest) to  
provide the permission for Zues access AWS.

```plaintext
variable "aws_region" {
  description = "Please enter the region used to deploy this infrastructure"
  type        = string
  default = "us-west-2"  
}

variable "cluster_id" {
  description = "Enter full name of EKS Cluster"
  type        = string
  default = "<Full Name EKS Cluster>" 
}

#Load informations of your EKS cluster
data "aws_eks_cluster" "eks_k8s" {
  name = var.cluster_id
}


module "eks-irsa" {
  source  = "mrnim94/eks-irsa/aws"
  version = "0.0.4"

  aws_region = var.aws_region
  environment = "dev"
  business_divsion = "irsa-zeus-rotations"

  k8s_namespace = "kube-system"
  k8s_service_account = "zeus-rotations"
  json_policy = <<EOT
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iam:CreateAccessKey",
        "iam:ListAccessKeys",
        "iam:DeleteAccessKey"
      ],
      "Resource": "*"
    }
  ]
}
EOT

  aws_iam_openid_connect_provider_arn = "arn:aws:iam::${element(split(":", "${data.aws_eks_cluster.eks_k8s.arn}"), 4)}:oidc-provider/${element(split("//", "${data.aws_eks_cluster.eks_k8s.identity[0].oidc[0].issuer}"), 1)}"
}

output "irsa_iam_role_arn" {
  description = "aws_iam_openid_connect_provider_arn"
  value = module.eks-irsa.irsa_iam_role_arn
}
```

**Second step:** Install Zeus via Helm chart with IRSA that is created in the previous step.

```python
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: zeus-rotations-mdaas-dev
  namespace: argocd
spec:
  destination:
    namespace: kube-system
    name: 'arn:aws:eks:us-west-2:043XXXXX1869:cluster/<Full Name EKS Cluster>'
  project: meta-structure
  source:
    repoURL: https://mrnim94.github.io/zeus
    targetRevision: "0.1.5"
    chart: zeus
    helm:
      values: |
        image:
          repository: "mrnim94/zeus-rotations"
          pullPolicy: IfNotPresent
          tag: "master"
        serviceAccount:
          annotations:
            eks.amazonaws.com/role-arn: "arn:aws:iam::043XXXXX1869:role/irsa-zeus-rotations-dev-irsa-iam-role"
        envVars:
          AWS_REGION: us-east-1
        config:
          schedules:
          - name: change-credetial-aws
            cron: "*/1 * * * *"
            usernameOnAws: nimtechnology
            locations:
              - secretName: credentials-aws
                namespaceOnK8s: default
                style: CredentialOnK8s
                credentialOnK8s: credentials
                profile: dev
              - secretName: secret-aws
                namespaceOnK8s: custom-namespace
                style: AccessKeyOnK8s
                accessKeyOnK8s: accesskey
                secretKeyOnK8s: secretkey
            restartWorkloads:
              - kind: deployment
                name: "argo-workflow-argo-workflows-server"
                namespaceOnK8s: workload-namespace
          - name: remove-access-key
            cron: "*/1 * * * *"
            usernameOnAws: nimtechnology
            action: OnlyDeleteAccessKey
```

### Explain the schedule configuration.

| Level 1 | Level 2 | Level 3 | Value | Type | Description |
| --- | --- | --- | --- | --- | --- |
| **schedules** |   |   | \- (list of schedules) | List | Top-level list for all schedule configurations |
|   | name |   | change-credetial-aws | String | Name of the schedule |
|   | cron |   | \*/1 \* \* \* \* | Cron String | Cron schedule, runs every minute |
|   | usernameOnAws |   | nimtechnology | String | AWS username |
|   | locations |   | \- (list of locations) | List | List of location configurations for the schedule |
|   |   | secretName | credentials-aws | String | Name of the secret in Kubernetes |
|   |   | namespaceOnK8s | default | String | Kubernetes namespace |
|   |   | style | `CredentialOnK8s` or `AccessKeyOnK8s` | String | Style/type of the credential |
|   |   | **credentialOnK8s** (require when style is `CredentialOnK8s`) | credentials | String | Key Name of Secret is holding AWS credential |
|   |   | **profile** (require when style is `CredentialOnK8s`) | dev | String | AWS profile in credential that you want to change |
|   |   | **accessKeyOnK8s** (require when style is `AccessKeyOnK8s`) | accesskey | String | Key Name of Secret is holding AWS access key |
|   |   | **secretKeyOnK8s** (require when style is `AccessKeyOnK8s`) | secretkey | String | Key Name of Secret is holding AWS secret key |
|   | restartWorkloads |   | \- (list of workloads) | List | List of workloads to restart on schedule change |
|   |   | kind | deployment | String | Type of the Kubernetes workload |
|   |   | name | argo-workflow-argo-workflows-server | String | Name of the Kubernetes workload |
|   |   | namespaceOnK8s | default | String | Kubernetes namespace |
|   | **action** |   | nil/null **(optional)** | String | declaring an extra **action** field with the value **OnlyDeleteAccessKey**. Zeus rotation only do a work is remove Access Key of AWS's Acount |

## zeus-rotations on quay.io

If your eks encounter the pull rate limits with the images on Docker Hub.  
YOu can use image on Quay.io

```plaintext
docker pull quay.io/nimtechnology/zeus-rotations
```

## Publish Helm Chart

```plaintext
helm package ./helm-chart/zeus --destination ./helm-chart/
helm repo index . --url https://mrnim94.github.io/zeus
```