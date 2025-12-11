package app

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

func (s *AppAPIHandler) BatchGetDevice(ctx context.Context, req *pb.BatchGetDeviceRequest, rsp *pb.BatchGetDeviceResponse) error {

	// 2.发送批量获取设备的请求
	// 数据转化为Svc的数据
	svcDevices := make([]string, len(req.GetDevices()))
	for k, v := range req.GetDevices() {
		svcDevices[k] = v.GetMac()
	}
	// 发送请求
	getDevicesRsp, err := s.deviceAPI.ListDevices(ctx, &devicepb.ListDevicesRequest{
		Mac: svcDevices,
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据给app
	appDevice := make([]*pb.Device, len(getDevicesRsp.GetDevices()))
	for k, v := range getDevicesRsp.GetDevices() {
		appDevice[k] = toAppDevice(v)
	}

	rsp.Devices = appDevice
	return nil
}
