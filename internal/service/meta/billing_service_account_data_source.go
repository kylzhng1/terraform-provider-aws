package meta

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
)

func init() {
	registerFrameworkDataSourceFactory(newDataSourceBillingServiceAccount)
}

// newDataSourceBillingServiceAccount instantiates a new DataSource for the aws_billing_service_account data source.
func newDataSourceBillingServiceAccount(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceBillingServiceAccount{}, nil
}

type dataSourceBillingServiceAccount struct {
	framework.DataSourceWithConfigure
}

// Metadata should return the full name of the data source, such as
// examplecloud_thing.
func (d *dataSourceBillingServiceAccount) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) { // nosemgrep:ci.meta-in-func-name
	response.TypeName = "aws_billing_service_account"
}

// Schema returns the schema for this data source.
func (d *dataSourceBillingServiceAccount) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
// Config values should be read from the ReadRequest and new state values set on the ReadResponse.
func (d *dataSourceBillingServiceAccount) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data dataSourceBillingServiceAccountData

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// See http://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/billing-getting-started.html#step-2
	const billingAccountID = "386209384616"

	arn := arn.ARN{
		Partition: d.Meta().Partition,
		Service:   "iam",
		AccountID: billingAccountID,
		Resource:  "root",
	}

	data.ARN = types.StringValue(arn.String())
	data.ID = types.StringValue(billingAccountID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

type dataSourceBillingServiceAccountData struct {
	ARN types.String `tfsdk:"arn"`
	ID  types.String `tfsdk:"id"`
}
