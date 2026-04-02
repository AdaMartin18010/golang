# 基础设施即代码 (IaC)

> **分类**: 成熟应用领域

---

## Pulumi

```go
import (
    "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // 创建 S3 Bucket
        bucket, err := s3.NewBucket(ctx, "my-bucket", &s3.BucketArgs{
            Website: &s3.BucketWebsiteArgs{
                IndexDocument: pulumi.String("index.html"),
            },
        })
        if err != nil {
            return err
        }

        // 导出 bucket 名称
        ctx.Export("bucketName", bucket.ID())
        ctx.Export("bucketEndpoint", bucket.WebsiteEndpoint())

        return nil
    })
}
```

---

## CDK for Terraform

```go
import (
    "github.com/aws/constructs-go/constructs/v10"
    "github.com/aws/jsii-runtime-go"
    "github.com/hashicorp/terraform-cdk-go/cdktf"
    "github.com/cdktf/cdktf-provider-aws-go/aws/v19/instance"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
    stack := cdktf.NewTerraformStack(scope, &id)

    instance.NewInstance(stack, jsii.String("instance"), &instance.InstanceConfig{
        Ami:          jsii.String("ami-12345678"),
        InstanceType: jsii.String("t2.micro"),
    })

    return stack
}

func main() {
    app := cdktf.NewApp(nil)
    NewMyStack(app, "cdktf-go")
    app.Synth()
}
```

---

## 云 SDK

### AWS SDK

```go
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/ec2"

cfg, _ := config.LoadDefaultConfig(context.TODO())
client := ec2.NewFromConfig(cfg)

// 创建 EC2
result, _ := client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
    ImageId:      aws.String("ami-12345678"),
    InstanceType: types.InstanceTypeT2Micro,
    MinCount:     aws.Int32(1),
    MaxCount:     aws.Int32(1),
})
```

### Azure SDK

```go
import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"

cred, _ := azidentity.NewDefaultAzureCredential(nil)
client, _ := armcompute.NewVirtualMachinesClient("subscription-id", cred, nil)
```
