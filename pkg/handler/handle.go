package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogo/status"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"

	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1alpha"
)

var forward_PipelinePublicService_TriggerUserPipeline_0 = runtime.ForwardResponseMessage
var forward_PipelinePublicService_TriggerUserPipelineRelease_0 = runtime.ForwardResponseMessage

func convertFormData(ctx context.Context, mux *runtime.ServeMux, req *http.Request) ([]*structpb.Struct, error) {

	err := req.ParseMultipartForm(4 << 20)
	if err != nil {
		return nil, err
	}

	inputsMap := map[int]map[string]interface{}{}

	maxInputIdx := 0

	for k, v := range req.MultipartForm.Value {
		if strings.HasPrefix(k, "inputs[") {
			k = k[7:]

			inputIdx, err := strconv.Atoi(k[:strings.Index(k, "]")])
			if err != nil {
				return nil, err
			}

			if inputIdx > maxInputIdx {
				maxInputIdx = inputIdx
			}

			k = k[strings.Index(k, "]")+2:]

			var key string
			isArray := false
			keyIdx := 0
			if strings.Contains(k, "[") {
				key = k[:strings.Index(k, "[")]
				keyIdx, err = strconv.Atoi(k[len(key)+1 : strings.Index(k, "]")])
				if err != nil {
					return nil, err
				}
				isArray = true
			} else {
				key = k
			}

			if _, ok := inputsMap[inputIdx]; !ok {
				inputsMap[inputIdx] = map[string]interface{}{}
			}

			if isArray {
				if _, ok := inputsMap[inputIdx][key]; !ok {
					inputsMap[inputIdx][key] = map[int]interface{}{}
				}
				var b interface{}
				unmarshalErr := json.Unmarshal([]byte(v[0]), &b)
				if unmarshalErr != nil {
					return nil, unmarshalErr
				}
				inputsMap[inputIdx][key].(map[int]interface{})[keyIdx] = b
			} else {
				var b interface{}
				unmarshalErr := json.Unmarshal([]byte(v[0]), &b)
				if unmarshalErr != nil {
					return nil, unmarshalErr
				}
				inputsMap[inputIdx][key] = b
			}

		}
	}

	for k, v := range req.MultipartForm.File {
		if strings.HasPrefix(k, "inputs[") {
			k = k[7:]

			inputIdx, err := strconv.Atoi(k[:strings.Index(k, "]")])
			if err != nil {
				return nil, err
			}

			if inputIdx > maxInputIdx {
				maxInputIdx = inputIdx
			}

			k = k[strings.Index(k, "]")+2:]

			var key string
			isArray := false
			keyIdx := 0
			if strings.Contains(k, "[") {
				key = k[:strings.Index(k, "[")]
				keyIdx, err = strconv.Atoi(k[len(key)+1 : strings.Index(k, "]")])
				if err != nil {
					return nil, err
				}
				isArray = true
			} else {
				key = k
			}

			if _, ok := inputsMap[inputIdx]; !ok {
				inputsMap[inputIdx] = map[string]interface{}{}
			}

			file, err := v[0].Open()
			if err != nil {
				return nil, err
			}

			byteContainer, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			v := base64.StdEncoding.EncodeToString(byteContainer)
			if isArray {
				if _, ok := inputsMap[inputIdx][key]; !ok {
					inputsMap[inputIdx][key] = map[int]interface{}{}
				}

				inputsMap[inputIdx][key].(map[int]interface{})[keyIdx] = v
			} else {
				inputsMap[inputIdx][key] = v
			}

		}
	}

	inputs := make([]*structpb.Struct, maxInputIdx+1)
	for inputIdx, inputValue := range inputsMap {
		inputs[inputIdx] = &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		}
		for key, value := range inputValue {

			switch value := value.(type) {
			case map[int]interface{}:
				maxItemIdx := 0
				for itemIdx := range value {
					if itemIdx > maxItemIdx {
						maxItemIdx = itemIdx
					}
				}
				vals := make([]interface{}, maxItemIdx+1)
				for itemIdx, itemValue := range value {
					vals[itemIdx] = itemValue
				}

				structVal, err := structpb.NewList(vals)
				if err != nil {
					return nil, err
				}

				inputs[inputIdx].GetFields()[key] = structpb.NewListValue(structVal)

			default:
				structVal, err := structpb.NewValue(value)
				if err != nil {
					return nil, err
				}
				inputs[inputIdx].GetFields()[key] = structVal
			}

		}
	}
	return inputs, nil
}

// HandleTrigger
func HandleTrigger(mux *runtime.ServeMux, client pipelinePB.PipelinePublicServiceClient, w http.ResponseWriter, req *http.Request, pathParams map[string]string) {

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
	var err error
	var annotatedContext context.Context
	var resp protoreflect.ProtoMessage
	var md runtime.ServerMetadata

	annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/vdp.pipeline.v1alpha.PipelinePublicService/TriggerUserPipeline", runtime.WithHTTPPathPattern("/v1alpha/{name=users/*/pipelines/*}/trigger"))
	if err != nil {
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		inputs, err := convertFormData(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err = request_PipelinePublicService_TriggerUserPipeline_0_form(annotatedContext, inboundMarshaler, client, &pipelinePB.TriggerUserPipelineRequest{
			Inputs: inputs,
		}, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

	} else {
		resp, md, err = request_PipelinePublicService_TriggerUserPipeline_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
	}

	annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)

	forward_PipelinePublicService_TriggerUserPipeline_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

}

// HandleTriggerAsync
func HandleTriggerAsync(mux *runtime.ServeMux, client pipelinePB.PipelinePublicServiceClient, w http.ResponseWriter, req *http.Request, pathParams map[string]string) {

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
	var err error
	var annotatedContext context.Context
	var resp protoreflect.ProtoMessage
	var md runtime.ServerMetadata

	annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/vdp.pipeline.v1alpha.PipelinePublicService/TriggerAsyncUserPipeline", runtime.WithHTTPPathPattern("/v1alpha/{name=users/*/pipelines/*}/triggerAsync"))
	if err != nil {
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		inputs, err := convertFormData(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err = request_PipelinePublicService_TriggerAsyncUserPipeline_0_form(annotatedContext, inboundMarshaler, client, &pipelinePB.TriggerAsyncUserPipelineRequest{
			Inputs: inputs,
		}, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

	} else {
		resp, md, err = request_PipelinePublicService_TriggerAsyncUserPipeline_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
	}

	annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)

	forward_PipelinePublicService_TriggerUserPipeline_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

}

// ref: the generated protogen-go files
func request_PipelinePublicService_TriggerUserPipeline_0(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq pipelinePB.TriggerUserPipelineRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerUserPipeline(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// ref: the generated protogen-go files
func request_PipelinePublicService_TriggerUserPipeline_0_form(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, protoReq *pipelinePB.TriggerUserPipelineRequest, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerUserPipeline(ctx, protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_PipelinePublicService_TriggerAsyncUserPipeline_0(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq pipelinePB.TriggerAsyncUserPipelineRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerAsyncUserPipeline(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// ref: the generated protogen-go files
func request_PipelinePublicService_TriggerAsyncUserPipeline_0_form(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, protoReq *pipelinePB.TriggerAsyncUserPipelineRequest, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerAsyncUserPipeline(ctx, protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// HandleTrigger
func HandleTriggerRelease(mux *runtime.ServeMux, client pipelinePB.PipelinePublicServiceClient, w http.ResponseWriter, req *http.Request, pathParams map[string]string) {

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
	var err error
	var annotatedContext context.Context
	var resp protoreflect.ProtoMessage
	var md runtime.ServerMetadata

	annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/vdp.pipeline.v1alpha.PipelinePublicService/TriggerUserPipelineRelease", runtime.WithHTTPPathPattern("/v1alpha/{name=users/*/pipelines/*/releases/*}/trigger"))
	if err != nil {
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		inputs, err := convertFormData(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err = request_PipelinePublicService_TriggerUserPipelineRelease_0_form(annotatedContext, inboundMarshaler, client, &pipelinePB.TriggerUserPipelineReleaseRequest{
			Inputs: inputs,
		}, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

	} else {
		resp, md, err = request_PipelinePublicService_TriggerUserPipelineRelease_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
	}

	annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)

	forward_PipelinePublicService_TriggerUserPipelineRelease_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

}

// HandleTriggerAsync
func HandleTriggerAsyncRelease(mux *runtime.ServeMux, client pipelinePB.PipelinePublicServiceClient, w http.ResponseWriter, req *http.Request, pathParams map[string]string) {

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
	var err error
	var annotatedContext context.Context
	var resp protoreflect.ProtoMessage
	var md runtime.ServerMetadata

	annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/vdp.pipeline.v1alpha.PipelinePublicService/TriggerAsyncUserPipelineRelease", runtime.WithHTTPPathPattern("/v1alpha/{name=users/*/pipelines/*/releases/*}/triggerAsync"))
	if err != nil {
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		inputs, err := convertFormData(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err = request_PipelinePublicService_TriggerAsyncUserPipelineRelease_0_form(annotatedContext, inboundMarshaler, client, &pipelinePB.TriggerAsyncUserPipelineReleaseRequest{
			Inputs: inputs,
		}, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

	} else {
		resp, md, err = request_PipelinePublicService_TriggerAsyncUserPipelineRelease_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
	}

	annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)

	forward_PipelinePublicService_TriggerUserPipelineRelease_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

}

// ref: the generated protogen-go files
func request_PipelinePublicService_TriggerUserPipelineRelease_0(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq pipelinePB.TriggerUserPipelineReleaseRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerUserPipelineRelease(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// ref: the generated protogen-go files
func request_PipelinePublicService_TriggerUserPipelineRelease_0_form(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, protoReq *pipelinePB.TriggerUserPipelineReleaseRequest, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerUserPipelineRelease(ctx, protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_PipelinePublicService_TriggerAsyncUserPipelineRelease_0(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq pipelinePB.TriggerAsyncUserPipelineReleaseRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerAsyncUserPipelineRelease(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// ref: the generated protogen-go files
func request_PipelinePublicService_TriggerAsyncUserPipelineRelease_0_form(ctx context.Context, marshaler runtime.Marshaler, client pipelinePB.PipelinePublicServiceClient, protoReq *pipelinePB.TriggerAsyncUserPipelineReleaseRequest, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "name")
	}

	protoReq.Name, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "name", err)
	}

	msg, err := client.TriggerAsyncUserPipelineRelease(ctx, protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}
