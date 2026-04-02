# Terraform Providers

> **分类**: 成熟应用领域

---

## Provider 开发

```go
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Provider() *schema.Provider {
    return &schema.Provider{
        Schema: map[string]*schema.Schema{
            "api_key": {
                Type:        schema.TypeString,
                Required:    true,
                Sensitive:   true,
                DefaultFunc: schema.EnvDefaultFunc("API_KEY", nil),
            },
        },
        ResourcesMap: map[string]*schema.Resource{
            "mycloud_server": resourceServer(),
        },
    }
}
```

---

## Resource 实现

```go
func resourceServer() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceServerCreate,
        ReadContext:   resourceServerRead,
        UpdateContext: resourceServerUpdate,
        DeleteContext: resourceServerDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "size": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    name := d.Get("name").(string)
    size := d.Get("size").(string)

    server, err := client.CreateServer(name, size)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(server.ID)
    return nil
}
```
