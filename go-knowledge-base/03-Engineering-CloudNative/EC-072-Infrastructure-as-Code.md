# Infrastructure as Code

> **分类**: 工程与云原生
> **标签**: #iac #terraform #pulumi #cloudformation #automation
> **参考**: Terraform Best Practices, AWS Well-Architected, Azure CAF

---

## 1. Formal Definition

### 1.1 What is Infrastructure as Code?

Infrastructure as Code (IaC) is the practice of managing and provisioning computing infrastructure through machine-readable definition files, rather than physical hardware configuration or interactive configuration tools. IaC enables infrastructure to be versioned, tested, and deployed using the same workflows as application code.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Infrastructure as Code Lifecycle                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   WRITE ─────────────────────────────────────────────────────────────────►  │
│     │                                                                       │
│     │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │
│     └─►│   Define     │───►│   Validate   │───►│   Format     │            │
│        │  Resources   │    │   Syntax     │    │   Code       │            │
│        └──────────────┘    └──────────────┘    └──────────────┘            │
│                                                                             │
│   PLAN ──────────────────────────────────────────────────────────────────►  │
│     │                                                                       │
│     │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │
│     └─►│  Initialize  │───►│  Refresh     │───►│  Plan        │            │
│        │   Backend    │    │   State      │    │  Changes     │            │
│        └──────────────┘    └──────────────┘    └──────────────┘            │
│                                                                             │
│   REVIEW ────────────────────────────────────────────────────────────────►  │
│     │                                                                       │
│     │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │
│     └─►│  Security    │───►│  Cost        │───►│  Peer        │            │
│        │  Scan        │    │  Estimate    │    │  Review      │            │
│        └──────────────┘    └──────────────┘    └──────────────┘            │
│                                                                             │
│   APPLY ─────────────────────────────────────────────────────────────────►  │
│     │                                                                       │
│     │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │
│     └─►│  Apply       │───►│  Verify      │───►│  Document    │            │
│        │  Changes     │    │  Resources   │    │  State       │            │
│        └──────────────┘    └──────────────┘    └──────────────┘            │
│                                                                             │
│   MONITOR ───────────────────────────────────────────────────────────────►  │
│     │                                                                       │
│     │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │
│     └─►│  Drift       │───►│  Compliance  │───►│  Optimize    │            │
│        │  Detection   │    │  Check       │    │  Resources   │            │
│        └──────────────┘    └──────────────┘    └──────────────┘            │
│                                                                             │
│   KEY PRINCIPLES:                                                           │
│   • Version Control: All code in Git with branch protection                 │
│   • Idempotency: Same result on multiple runs                               │
│   • Immutability: Replace rather than modify                                │
│   • Declarative: Define desired state, not steps                            │
│   • Modularity: Reusable components                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 IaC Approaches

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        IaC Approaches                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DECLARATIVE (What)                    IMPERATIVE (How)                     │
│  ━━━━━━━━━━━━━━━━━━━━                  ━━━━━━━━━━━━━━━━━━                   │
│                                                                             │
│  Example: Terraform                    Example: AWS CLI, Ansible            │
│                                                                             │
│  resource "aws_instance" "web" {       aws ec2 run-instances               │
│    ami           = "ami-12345"           --image-id ami-12345               │
│    instance_type = "t3.micro"            --instance-type t3.micro            │
│    tags = {                              --tag-specifications ...           │
│      Name = "WebServer"                                                    │
│    }                                                                     │
│  }                                                                       │
│                                                                             │
│  Pros:                                 Pros:                                │
│  • Idempotent by design                • Fine-grained control               │
│  • Self-documenting                    • Better for one-off tasks           │
│  • State tracking                      • Procedural logic support           │
│  • Drift detection                     • Better for complex workflows       │
│                                                                             │
│  Cons:                                 Cons:                                │
│  • Learning curve                      • Not inherently idempotent          │
│  • State management complexity         • Harder to maintain                 │
│  • Less flexible for complex logic     • No drift detection                 │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  MUTABLE                             IMMUTABLE                              │
│  ━━━━━━━━                            ━━━━━━━━━━━━                           │
│                                                                             │
│  Update existing resources           Replace with new resources             │
│  Configuration management            Infrastructure provisioning            │
│  Ansible, Chef, Puppet               Terraform, CloudFormation, Pulumi      │
│                                                                             │
│  • Gradual changes                   • Clean slate deployments              │
│  • Risk of configuration drift       • Predictable outcomes                 │
│  • Faster for small changes          • Easier rollback                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns

### 2.1 Terraform Module Structure

```hcl
# modules/vpc/main.tf
# Virtual Private Cloud module

locals {
  common_tags = merge(
    var.additional_tags,
    {
      Environment = var.environment
      ManagedBy   = "terraform"
      Project     = var.project_name
    }
  )
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-vpc"
    }
  )

  lifecycle {
    prevent_destroy = var.prevent_destroy
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-igw"
    }
  )
}

# Public Subnets
resource "aws_subnet" "public" {
  count = length(var.availability_zones)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-public-${count.index + 1}"
      Type = "public"
    }
  )
}

# Private Subnets
resource "aws_subnet" "private" {
  count = length(var.availability_zones)

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 100)
  availability_zone = var.availability_zones[count.index]

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-private-${count.index + 1}"
      Type = "private"
    }
  )
}

# NAT Gateways (one per AZ for high availability)
resource "aws_eip" "nat" {
  count = var.enable_nat_gateway ? length(var.availability_zones) : 0

  domain = "vpc"

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-nat-eip-${count.index + 1}"
    }
  )

  depends_on = [aws_internet_gateway.main]
}

resource "aws_nat_gateway" "main" {
  count = var.enable_nat_gateway ? length(var.availability_zones) : 0

  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-nat-${count.index + 1}"
    }
  )
}

# Route Tables
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-public-rt"
    }
  )
}

resource "aws_route_table" "private" {
  count = length(var.availability_zones)

  vpc_id = aws_vpc.main.id

  dynamic "route" {
    for_each = var.enable_nat_gateway ? [1] : []
    content {
      cidr_block     = "0.0.0.0/0"
      nat_gateway_id = aws_nat_gateway.main[count.index].id
    }
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${var.project_name}-${var.environment}-private-rt-${count.index + 1}"
    }
  )
}

# Route Table Associations
resource "aws_route_table_association" "public" {
  count = length(var.availability_zones)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count = length(var.availability_zones)

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# VPC Flow Logs
resource "aws_flow_log" "main" {
  count = var.enable_flow_logs ? 1 : 0

  vpc_id                   = aws_vpc.main.id
  traffic_type             = "ALL"
  log_destination_type     = "cloud-watch-logs"
  log_destination          = aws_cloudwatch_log_group.flow_logs[0].arn
  iam_role_arn             = aws_iam_role.flow_logs[0].arn
  max_aggregation_interval = 60

  tags = local.common_tags
}

resource "aws_cloudwatch_log_group" "flow_logs" {
  count = var.enable_flow_logs ? 1 : 0

  name              = "/aws/vpc/${var.project_name}-${var.environment}-flow-logs"
  retention_in_days = var.flow_logs_retention_days

  tags = local.common_tags
}
```

```hcl
# modules/vpc/variables.tf
variable "environment" {
  description = "Environment name"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be dev, staging, or prod."
  }
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
}

variable "enable_nat_gateway" {
  description = "Enable NAT gateway for private subnets"
  type        = bool
  default     = true
}

variable "enable_flow_logs" {
  description = "Enable VPC flow logs"
  type        = bool
  default     = true
}

variable "flow_logs_retention_days" {
  description = "Flow logs retention in days"
  type        = number
  default     = 30
}

variable "prevent_destroy" {
  description = "Prevent destruction of VPC"
  type        = bool
  default     = true
}

variable "additional_tags" {
  description = "Additional tags to apply"
  type        = map(string)
  default     = {}
}
```

```hcl
# modules/vpc/outputs.tf
output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = aws_subnet.private[*].id
}

output "nat_gateway_ids" {
  description = "List of NAT gateway IDs"
  value       = aws_nat_gateway.main[*].id
}

output "availability_zones" {
  description = "List of availability zones used"
  value       = var.availability_zones
}
```

### 2.2 Pulumi Go Implementation

```go
package main

import (
    "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// VPCComponent encapsulates VPC resources
type VPCComponent struct {
    pulumi.ResourceState

    VPC            *ec2.Vpc
    PublicSubnets  []*ec2.Subnet
    PrivateSubnets []*ec2.Subnet
    IGW            *ec2.InternetGateway
    NATGateways    []*ec2.NatGateway
}

// VPCArgs defines VPC configuration
type VPCArgs struct {
    CIDRBlock          string
    AvailabilityZones  []string
    EnableNATGateway   bool
    EnableFlowLogs     bool
    Tags               pulumi.StringMap
}

// NewVPCComponent creates a new VPC
func NewVPCComponent(ctx *pulumi.Context, name string, args *VPCArgs, opts ...pulumi.ResourceOption) (*VPCComponent, error) {
    component := &VPCComponent{}
    err := ctx.RegisterComponentResource("custom:resource:VPC", name, component, opts...)
    if err != nil {
        return nil, err
    }

    // Create VPC
    vpc, err := ec2.NewVpc(ctx, name, &ec2.VpcArgs{
        CidrBlock:          pulumi.String(args.CIDRBlock),
        EnableDnsHostnames: pulumi.Bool(true),
        EnableDnsSupport:   pulumi.Bool(true),
        Tags: pulumi.StringMap{
            "Name": pulumi.String(name),
        },
    }, pulumi.Parent(component))
    if err != nil {
        return nil, err
    }
    component.VPC = vpc

    // Create Internet Gateway
    igw, err := ec2.NewInternetGateway(ctx, name+"-igw", &ec2.InternetGatewayArgs{
        VpcId: vpc.ID(),
        Tags: pulumi.StringMap{
            "Name": pulumi.String(name + "-igw"),
        },
    }, pulumi.Parent(component))
    if err != nil {
        return nil, err
    }
    component.IGW = igw

    // Create Subnets
    for i, az := range args.AvailabilityZones {
        // Public subnet
        publicSubnet, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-public-%d", name, i), &ec2.SubnetArgs{
            VpcId:               vpc.ID(),
            CidrBlock:           pulumi.String(cidrSubnet(args.CIDRBlock, 8, i)),
            AvailabilityZone:    pulumi.String(az),
            MapPublicIpOnLaunch: pulumi.Bool(true),
            Tags: pulumi.StringMap{
                "Name": pulumi.String(fmt.Sprintf("%s-public-%d", name, i)),
                "Type": pulumi.String("public"),
            },
        }, pulumi.Parent(component))
        if err != nil {
            return nil, err
        }
        component.PublicSubnets = append(component.PublicSubnets, publicSubnet)

        // Private subnet
        privateSubnet, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-private-%d", name, i), &ec2.SubnetArgs{
            VpcId:            vpc.ID(),
            CidrBlock:        pulumi.String(cidrSubnet(args.CIDRBlock, 8, i+100)),
            AvailabilityZone: pulumi.String(az),
            Tags: pulumi.StringMap{
                "Name": pulumi.String(fmt.Sprintf("%s-private-%d", name, i)),
                "Type": pulumi.String("private"),
            },
        }, pulumi.Parent(component))
        if err != nil {
            return nil, err
        }
        component.PrivateSubnets = append(component.PrivateSubnets, privateSubnet)
    }

    return component, nil
}

// cidrSubnet calculates subnet CIDR
func cidrSubnet(cidr string, newbits, netnum int) string {
    // Simplified implementation
    return fmt.Sprintf("10.0.%d.0/24", netnum)
}

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "")

        environment := cfg.Require("environment")

        // Create VPC
        vpc, err := NewVPCComponent(ctx, "main-vpc", &VPCArgs{
            CIDRBlock:         "10.0.0.0/16",
            AvailabilityZones: []string{"us-east-1a", "us-east-1b", "us-east-1c"},
            EnableNATGateway:  true,
            EnableFlowLogs:    true,
        })
        if err != nil {
            return err
        }

        // Export outputs
        ctx.Export("vpcId", vpc.VPC.ID())
        ctx.Export("publicSubnetIds", pulumi.ToStringArrayOutput(
            func() []pulumi.StringOutput {
                outputs := make([]pulumi.StringOutput, len(vpc.PublicSubnets))
                for i, subnet := range vpc.PublicSubnets {
                    outputs[i] = subnet.ID()
                }
                return outputs
            }(),
        ))

        return nil
    })
}
```

---

## 3. Production-Ready Configurations

### 3.1 Terraform Workspace Structure

```
terraform/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── backend.tf
│   │   └── terraform.tfvars
│   ├── staging/
│   │   └── ...
│   └── prod/
│       └── ...
├── modules/
│   ├── vpc/
│   ├── eks/
│   ├── rds/
│   └── iam/
├── policies/
│   └── sentinel/
└── scripts/
    └── setup.sh
```

```hcl
# environments/prod/backend.tf
terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "terraform-state-prod"
    key            = "infrastructure/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = "production"
      ManagedBy   = "terraform"
    }
  }
}
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    IaC Security Best Practices                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  STATE MANAGEMENT                                                           │
│  ✓ Remote state with encryption at rest                                     │
│  ✓ State locking to prevent conflicts                                       │
│  ✓ Sensitive data in state - use encryption                                 │
│  ✓ Regular state backups                                                    │
│                                                                             │
│  SECRETS MANAGEMENT                                                         │
│  ✓ No hardcoded secrets in code                                             │
│  ✓ Use secret stores (Vault, AWS Secrets Manager)                           │
│  ✓ Mark sensitive outputs                                                   │
│  ✓ Rotate credentials automatically                                         │
│                                                                             │
│  ACCESS CONTROL                                                             │
│  ✓ Least privilege for deployment credentials                               │
│  ✓ MFA for state access                                                     │
│  ✓ Service accounts per environment                                         │
│  ✓ Audit all changes                                                        │
│                                                                             │
│  VALIDATION                                                                 │
│  ✓ Policy as Code (OPA, Sentinel, Checkov)                                  │
│  ✓ Pre-commit hooks for security scans                                      │
│  ✓ Cost estimation before apply                                             │
│  ✓ Drift detection                                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Decision Matrices

### 5.1 IaC Tool Selection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       IaC Tool Comparison Matrix                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Tool        │  Best For              │  Learning │  State    │  Ecosystem │
│              │                        │  Curve    │  Mgmt     │            │
├──────────────┼────────────────────────┼───────────┼───────────┼────────────│
│  Terraform   │  Multi-cloud,          │  Medium   │  Required │  ★★★★★     │
│              │  modularity            │           │           │            │
│  ────────────┼────────────────────────┼───────────┼───────────┼────────────│
│  Pulumi      │  Developer experience, │  Low      │  Required │  ★★★★☆     │
│              │  complex logic         │           │           │            │
│  ────────────┼────────────────────────┼───────────┼───────────┼────────────│
│  CloudFormation│ AWS-native,          │  Medium   │  Managed  │  ★★★☆☆     │
│              │  simple stacks         │           │           │            │
│  ────────────┼────────────────────────┼───────────┼───────────┼────────────│
│  ARM/Bicep   │  Azure-native          │  Low      │  Managed  │  ★★★☆☆     │
│  ────────────┼────────────────────────┼───────────┼───────────┼────────────│
│  Ansible     │  Configuration         │  Low      │  None     │  ★★★★☆     │
│              │  management            │           │           │            │
│  ────────────┼────────────────────────┼───────────┼───────────┼────────────│
│  Crossplane  │  Kubernetes-native     │  Medium   │  K8s      │  ★★★☆☆     │
│              │  control plane         │           │           │            │
│                                                                             │
│  Recommendation:                                                            │
│  • Multi-cloud: Terraform                                                   │
│  • Developer-centric org: Pulumi                                            │
│  • Single cloud: Native tool (CloudFormation/ARM)                           │
│  • Hybrid: Terraform + Ansible                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         IaC Best Practices Summary                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CODE ORGANIZATION                                                          │
│  ✓ Use modules for reusability                                              │
│  ✓ Environment-based directory structure                                    │
│  ✓ Consistent naming conventions                                            │
│  ✓ Version pin providers and modules                                        │
│                                                                             │
│  STATE MANAGEMENT                                                           │
│  ✓ Remote state with locking                                                │
│  ✓ State encryption at rest                                                 │
│  ✓ Separate state per environment                                           │
│  ✓ Regular state backups                                                    │
│                                                                             │
│  SECURITY                                                                   │
│  ✓ No secrets in code                                                       │
│  ✓ Policy as Code validation                                                │
│  ✓ Least privilege credentials                                              │
│  ✓ Audit logging enabled                                                    │
│                                                                             │
│  WORKFLOW                                                                   │
│  ✓ Git-based workflow with PR reviews                                       │
│  ✓ Automated CI/CD for terraform                                            │
│  ✓ Plan before apply                                                        │
│  ✓ Drift detection and remediation                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Terraform Best Practices
2. AWS Well-Architected IaC Lens
3. Azure Cloud Adoption Framework
4. Pulumi Documentation
5. Infrastructure as Code Book by Kief Morris
