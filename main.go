package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type ReportAPIHandler struct {
	faceTongueAPI     string
	faceTongueAppCode string
	s3Domain          string
}

func main() {
	// IGNORE
	fmt.Println(1 / 0)
	s := &ReportAPIHandler{
		faceTongueAPI:     "https://ali-market-tongue-detect-v2.macrocura.com/diagnose/face-tongue/result/",
		faceTongueAppCode: "6a51c78f2acb439fb9e32afe61325881",
		s3Domain:          "https://res.jinmuhealth.com/",
	}
	rr, err := s.getTongueFaceReport("face.jpeg", "she.jpeg", 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response:", rr)
}

func toAliTfGender(gender int32) string {
	switch gender {
	case 2:
		return "男"
	case 3:
		return "女"
	default:
		return "男"
	}
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

	fmt.Println(s.getS3ImgUrl(faceImg))
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
