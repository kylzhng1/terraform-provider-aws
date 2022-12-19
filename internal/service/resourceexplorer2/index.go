package resourceexplorer2

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
// 	awstypes "github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"

// 	"github.com/hashicorp/terraform-plugin-framework-timeouts/timeouts"
// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-provider-aws/internal/enum"
// 	"github.com/hashicorp/terraform-provider-aws/internal/errs"
// 	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
// 	"github.com/hashicorp/terraform-provider-aws/internal/flex"
// 	"github.com/hashicorp/terraform-provider-aws/internal/framework"
// 	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
// 	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
// )

// func init() {
// 	registerFrameworkResourceFactory(newResourceIndex)
// }

// func newResourceIndex(context.Context) (resource.ResourceWithConfigure, error) {
// 	return &resourceIndex{
// 		defaultCreateTimeout: 2 * time.Hour,
// 		defaultUpdateTimeout: 2 * time.Hour,
// 		defaultDeleteTimeout: 10 * time.Minute,
// 	}, nil
// }

// type resourceIndex struct {
// 	framework.ResourceWithConfigure

// 	defaultCreateTimeout time.Duration
// 	defaultUpdateTimeout time.Duration
// 	defaultDeleteTimeout time.Duration
// }

// func (r *resourceIndex) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
// 	response.TypeName = "aws_resourceexplorer2_index"
// }

// func (r *resourceIndex) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
// 	response.Schema = schema.Schema{
// 		Attributes: map[string]schema.Attribute{
// 			"arn": schema.StringAttribute{
// 				Computed: true,
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"id":       framework.IDAttribute(),
// 			"tags":     tftags.TagsAttribute(),
// 			"tags_all": tftags.TagsAttributeComputedOnly(),
// 			"type": schema.StringAttribute{
// 				Required: true,
// 				Validators: []validator.String{
// 					enum.FrameworkValidate[awstypes.IndexType](),
// 				},
// 			},
// 		},
// 		// Blocks: map[string]schema.Block{
// 		// 	"timeouts": timeouts.Block(ctx, timeouts.Opts{
// 		// 		Create: true,
// 		// 		Update: true,
// 		// 		Delete: true,
// 		// 	}),
// 		// },
// 	}
// }

// func (r *resourceIndex) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
// 	var data resourceIndexData

// 	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

// 	if response.Diagnostics.HasError() {
// 		return
// 	}

// 	conn := r.Meta().ResourceExplorer2Client
// 	createTimeout := timeouts.Create(ctx, data.Timeouts, r.defaultCreateTimeout)
// 	defaultTagsConfig := r.Meta().DefaultTagsConfig
// 	ignoreTagsConfig := r.Meta().IgnoreTagsConfig
// 	tags := defaultTagsConfig.MergeTags(tftags.New(data.Tags))

// 	input := &resourceexplorer2.CreateIndexInput{
// 		ClientToken: aws.String(sdkresource.UniqueId()),
// 	}

// 	if len(tags) > 0 {
// 		input.Tags = Tags(tags.IgnoreAWS())
// 	}

// 	output, err := conn.CreateIndex(ctx, input)

// 	if err != nil {
// 		response.Diagnostics.AddError("creating Resource Explorer Index", err.Error())

// 		return
// 	}

// 	arn := aws.ToString(output.Arn)
// 	data.ID = types.StringValue(arn)

// 	if _, err := waitIndexCreated(ctx, conn, createTimeout); err != nil {
// 		response.Diagnostics.AddError(fmt.Sprintf("waiting for Resource Explorer Index (%s) create", data.ID.ValueString()), err.Error())

// 		return
// 	}

// 	if data.Type.ValueString() == string(awstypes.IndexTypeAggregator) {
// 		input := &resourceexplorer2.UpdateIndexTypeInput{
// 			Arn:  flex.StringFromFramework(ctx, data.ID),
// 			Type: awstypes.IndexTypeAggregator,
// 		}

// 		_, err := conn.UpdateIndexType(ctx, input)

// 		if err != nil {
// 			response.Diagnostics.AddError(fmt.Sprintf("updating Resource Explorer Index (%s)", data.ID.ValueString()), err.Error())

// 			return
// 		}

// 		if _, err := waitIndexUpdated(ctx, conn, createTimeout); err != nil {
// 			response.Diagnostics.AddError(fmt.Sprintf("waiting for Resource Explorer Index (%s) update", data.ID.ValueString()), err.Error())

// 			return
// 		}
// 	}

// 	// Set values for unknowns.
// 	data.ARN = types.StringValue(arn)
// 	data.TagsAll = flex.FlattenFrameworkStringValueMap(ctx, tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map())

// 	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
// }

// func (r *resourceIndex) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
// 	var data resourceIndexData

// 	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

// 	if response.Diagnostics.HasError() {
// 		return
// 	}

// 	conn := r.Meta().ResourceExplorer2Client
// 	defaultTagsConfig := r.Meta().DefaultTagsConfig
// 	ignoreTagsConfig := r.Meta().IgnoreTagsConfig

// 	output, err := findIndex(ctx, conn)

// 	if tfresource.NotFound(err) {
// 		response.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
// 		response.State.RemoveResource(ctx)

// 		return
// 	}

// 	if err != nil {
// 		response.Diagnostics.AddError(fmt.Sprintf("reading Resource Explorer Index (%s)", data.ID.ValueString()), err.Error())

// 		return
// 	}

// 	data.ARN = flex.StringToFramework(ctx, output.Arn)
// 	data.Type = types.StringValue(string(output.Type))

// 	tags := KeyValueTags(output.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)
// 	// AWS APIs often return empty lists of tags when none have been configured.
// 	if tags := tags.RemoveDefaultConfig(defaultTagsConfig).Map(); len(tags) == 0 {
// 		data.Tags = tftags.Null
// 	} else {
// 		data.Tags = flex.FlattenFrameworkStringValueMap(ctx, tags)
// 	}
// 	data.TagsAll = flex.FlattenFrameworkStringValueMap(ctx, tags.Map())

// 	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
// }

// func (r *resourceIndex) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
// 	var old, new resourceIndexData

// 	response.Diagnostics.Append(request.State.Get(ctx, &old)...)

// 	if response.Diagnostics.HasError() {
// 		return
// 	}

// 	response.Diagnostics.Append(request.Plan.Get(ctx, &new)...)

// 	if response.Diagnostics.HasError() {
// 		return
// 	}

// 	conn := r.Meta().ResourceExplorer2Client
// 	updateTimeout := timeouts.Update(ctx, new.Timeouts, r.defaultUpdateTimeout)

// 	if !new.Type.Equal(old.Type) {
// 		input := &resourceexplorer2.UpdateIndexTypeInput{
// 			Arn:  flex.StringFromFramework(ctx, new.ID),
// 			Type: awstypes.IndexType(new.Type.ValueString()),
// 		}

// 		_, err := conn.UpdateIndexType(ctx, input)

// 		if err != nil {
// 			response.Diagnostics.AddError(fmt.Sprintf("updating Resource Explorer Index (%s)", new.ID.ValueString()), err.Error())

// 			return
// 		}

// 		if _, err := waitIndexUpdated(ctx, conn, updateTimeout); err != nil {
// 			response.Diagnostics.AddError(fmt.Sprintf("waiting for Resource Explorer Index (%s) update", new.ID.ValueString()), err.Error())

// 			return
// 		}
// 	}

// 	if !new.TagsAll.Equal(old.TagsAll) {
// 		if err := UpdateTags(ctx, conn, new.ID.ValueString(), old.TagsAll, new.TagsAll); err != nil {
// 			response.Diagnostics.AddError(fmt.Sprintf("updating Resource Explorer Index (%s) tags", new.ID.ValueString()), err.Error())

// 			return
// 		}
// 	}

// 	response.Diagnostics.Append(response.State.Set(ctx, &new)...)
// }

// func (r *resourceIndex) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
// 	var data resourceIndexData

// 	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

// 	if response.Diagnostics.HasError() {
// 		return
// 	}

// 	conn := r.Meta().ResourceExplorer2Client
// 	deleteTimeout := timeouts.Delete(ctx, data.Timeouts, r.defaultDeleteTimeout)

// 	tflog.Debug(ctx, "deleting Resource Explorer Index", map[string]interface{}{
// 		"id": data.ID.ValueString(),
// 	})
// 	_, err := conn.DeleteIndex(ctx, &resourceexplorer2.DeleteIndexInput{
// 		Arn: flex.StringFromFramework(ctx, data.ID),
// 	})

// 	if err != nil {
// 		response.Diagnostics.AddError(fmt.Sprintf("deleting Resource Explorer Index (%s)", data.ID.ValueString()), err.Error())

// 		return
// 	}

// 	if _, err := waitIndexDeleted(ctx, conn, deleteTimeout); err != nil {
// 		response.Diagnostics.AddError(fmt.Sprintf("waiting for Resource Explorer Index (%s) delete", data.ID.ValueString()), err.Error())

// 		return
// 	}
// }

// func (r *resourceIndex) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
// }

// func (r *resourceIndex) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
// 	r.SetTagsAll(ctx, request, response)
// }

// type resourceIndexData struct {
// 	ARN      types.String `tfsdk:"arn"`
// 	ID       types.String `tfsdk:"id"`
// 	Tags     types.Map    `tfsdk:"tags"`
// 	TagsAll  types.Map    `tfsdk:"tags_all"`
// 	Timeouts types.Object `tfsdk:"timeouts"`
// 	Type     types.String `tfsdk:"type"`
// }

// func findIndex(ctx context.Context, conn *resourceexplorer2.Client) (*resourceexplorer2.GetIndexOutput, error) {
// 	input := &resourceexplorer2.GetIndexInput{}

// 	output, err := conn.GetIndex(ctx, input)

// 	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
// 		return nil, &sdkresource.NotFoundError{
// 			LastError:   err,
// 			LastRequest: input,
// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	if output == nil {
// 		return nil, tfresource.NewEmptyResultError(input)
// 	}

// 	if state := output.State; state == awstypes.IndexStateDeleted {
// 		return nil, &sdkresource.NotFoundError{
// 			Message:     string(state),
// 			LastRequest: input,
// 		}
// 	}

// 	return output, nil
// }

// func statusIndex(ctx context.Context, conn *resourceexplorer2.Client) sdkresource.StateRefreshFunc {
// 	return func() (interface{}, string, error) {
// 		output, err := findIndex(ctx, conn)

// 		if tfresource.NotFound(err) {
// 			return nil, "", nil
// 		}

// 		if err != nil {
// 			return nil, "", err
// 		}

// 		return output, string(output.State), nil
// 	}
// }

// func waitIndexCreated(ctx context.Context, conn *resourceexplorer2.Client, timeout time.Duration) (*resourceexplorer2.GetIndexOutput, error) {
// 	stateConf := &sdkresource.StateChangeConf{
// 		Pending: enum.Slice(awstypes.IndexStateCreating),
// 		Target:  enum.Slice(awstypes.IndexStateActive),
// 		Refresh: statusIndex(ctx, conn),
// 		Timeout: timeout,
// 	}

// 	outputRaw, err := stateConf.WaitForStateContext(ctx)

// 	if output, ok := outputRaw.(*resourceexplorer2.GetIndexOutput); ok {
// 		return output, err
// 	}

// 	return nil, err
// }

// func waitIndexUpdated(ctx context.Context, conn *resourceexplorer2.Client, timeout time.Duration) (*resourceexplorer2.GetIndexOutput, error) { //nolint:unparam
// 	stateConf := &sdkresource.StateChangeConf{
// 		Pending: enum.Slice(awstypes.IndexStateUpdating),
// 		Target:  enum.Slice(awstypes.IndexStateActive),
// 		Refresh: statusIndex(ctx, conn),
// 		Timeout: timeout,
// 	}

// 	outputRaw, err := stateConf.WaitForStateContext(ctx)

// 	if output, ok := outputRaw.(*resourceexplorer2.GetIndexOutput); ok {
// 		return output, err
// 	}

// 	return nil, err
// }

// func waitIndexDeleted(ctx context.Context, conn *resourceexplorer2.Client, timeout time.Duration) (*resourceexplorer2.GetIndexOutput, error) {
// 	stateConf := &sdkresource.StateChangeConf{
// 		Pending: enum.Slice(awstypes.IndexStateDeleting),
// 		Target:  []string{},
// 		Refresh: statusIndex(ctx, conn),
// 		Timeout: timeout,
// 	}

// 	outputRaw, err := stateConf.WaitForStateContext(ctx)

// 	if output, ok := outputRaw.(*resourceexplorer2.GetIndexOutput); ok {
// 		return output, err
// 	}

// 	return nil, err
// }
