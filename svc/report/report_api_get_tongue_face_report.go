package report

import (
	"bytes"
	"context"
	"encoding/json"
	gerr "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取舌面报告
func (s *ReportAPIHandler) GetTongueFaceReport(ctx context.Context, req *pb.GetTongueFaceReportRequest, rsp *pb.GetTongueFaceReportResponse) error {
	err := validateGetTongueFaceReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 1. 获取报告
	report, err := s.reportStore.GetReportByReportID(ctx, req.GetReportId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if report == nil {
		return errors.Errorf(ErrReportNotFound, "report not found by report_id %s", req.GetReportId())
	}

	// 只有未完成的可以获取
	if report.GetTongueFaceReportStatus() != domain.TongueFaceReportStatusNotGet {
		return errors.Errorf(codes.InvalidRequest, "tongue face report status can not get by report_id %s", req.GetReportId())
	}

	// 获取性别
	gender := report.GetGender()

	// 2. 调用舌面分析接口
	respBody, err := s.getTongueFaceReport(req.GetFaceImageUrl(), req.GetTongueImageUrl(), gender)
	if err != nil {
		// 更新report 获取舌面失败
		err = s.reportStore.UpdateReportTongueFaceReportStatus(ctx, req.GetFaceImageUrl(), req.GetTongueImageUrl(), req.GetReportId(), err.Error(), domain.TongueFaceReportStatusFailed)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 返回数据
		rsp.TongueFaceReport = &pb.TongueFaceReport{
			Success: false,
		}
		return nil
	}

	// 更新report 获取舌面成功
	err = s.reportStore.UpdateReportTongueFaceReportStatus(ctx, req.GetFaceImageUrl(), req.GetTongueImageUrl(), req.GetReportId(), respBody, domain.TongueFaceReportStatusSuccess)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回数据
	rsp.TongueFaceReport = &pb.TongueFaceReport{
		Success:        true,
		TongueImageUrl: req.GetTongueImageUrl(),
		FaceImageUrl:   req.GetFaceImageUrl(),
		Data:           respBody,
	}

	return nil
}

// 验证request
func validateGetTongueFaceReportRequest(req *pb.GetTongueFaceReportRequest) error {
	if req.GetReportId() == "" {
		return gerr.New("report id should not be empty")
	}
	if req.GetTongueImageUrl() == "" {
		return gerr.New("tongue image url should not be empty")
	}
	if req.GetFaceImageUrl() == "" {
		return gerr.New("face image url should not be empty")
	}
	return nil
}

func (s *ReportAPIHandler) getTongueFaceReport(faceImg, tongueImg string, gender int32) (string, error) {
	// 定义请求体结构体（与 JSON 结构对应，方便序列化）
	type FaceTongueRequest struct {
		Scene   int    `json:"scene"`    // 场景值，固定为 2
		FfImage string `json:"ff_image"` // 人脸图片 URL
		TfImage string `json:"tf_image"` // 舌头正面图片 URL
		TbImage string `json:"tb_image"` // 舌头背面图片 URL
		Gender  string `json:"gender"`   // 性别：男/女
	}

	// 构造请求体数据
	requestBody := FaceTongueRequest{
		Scene:   2,                        // 固定值
		FfImage: s.getS3ImgUrl(faceImg),   // 替换为实际人脸图片 URL
		TfImage: s.getS3ImgUrl(tongueImg), // 替换为实际舌头正面图片 URL
		Gender:  toAliTfGender(gender),
	}

	// 2. 将请求体序列化为 JSON 字符串
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 3. 创建 HTTP 请求
	req, err := http.NewRequest("POST", s.faceTongueAPI, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// 4. 设置请求头
	req.Header.Set("Authorization", "APPCODE "+s.faceTongueAppCode) // 拼接 APPCODE
	req.Header.Set("Content-Type", "application/json")              // JSON 格式请求体

	// 5. 配置 HTTP 客户端（设置超时时间，避免无限等待）
	client := &http.Client{}

	// 6. 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // 确保响应体关闭，避免资源泄露

	// 7. 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 8. 检查响应状态码（200 表示成功，根据接口文档调整）
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码：%d，响应内容：%s", resp.StatusCode, string(respBody))
	}
	return string(respBody), nil
}

func (s *ReportAPIHandler) getS3ImgUrl(imgUrl string) string {
	if imgUrl == "" {
		return ""
	}
	link, _ := url.Parse(s.s3Domain)
	link.Path = path.Join(link.Path, imgUrl)
	return link.String()
}
