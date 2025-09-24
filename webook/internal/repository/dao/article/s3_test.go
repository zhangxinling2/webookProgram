package article

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestS3(t *testing.T) {
	// 腾讯云中对标 s3 和 OSS 的产品叫做 COS
	cosId, ok := os.LookupEnv("COS_APP_ID")
	if !ok {
		panic("没有找到环境变量 COS_APP_ID ")
	}
	cosKey, ok := os.LookupEnv("COS_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 COS_APP_SECRET")
	}
	// 1. 自定义 Endpoint Resolver 来处理腾讯云 COS 的特定端点
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://cos.ap-nanjing.myqcloud.com", // 你的腾讯云 COS 端点
			SigningRegion: "ap-nanjing",                          // 区域
			// 注意：部分 S3 兼容服务可能需要调整 SigningName
			// SourcesRegion: "ap-nanjing",
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-nanjing"), // 设置区域
		config.WithCredentialsProvider( // 设置静态凭证
			credentials.NewStaticCredentialsProvider(cosId, cosKey, ""),
		),
		config.WithEndpointResolverWithOptions(customResolver), // 设置自定义端点解析器
	)
	// 3. 创建 S3 客户端，并设置 UsePathStyle
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // 强制使用路径样式，即 /bucket/key 的形态
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 4. 上传对象 (PutObject)
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("webook-1314583317"), // 存储桶名称
		Key:         aws.String("126"),               // 对象键
		Body:        bytes.NewReader([]byte("测试内容 abc")),
		ContentType: aws.String("text/plain;charset=utf-8"),
	})
	if err != nil {
		panic("上传对象失败: " + err.Error())
	}

	// 5. 获取对象 (GetObject) - 注意 Key 需与上传时一致
	res, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("webook-1314583317"),
		Key:    aws.String("126"), // 注意：这里应使用上传时使用的 Key "126"，而非 "测试文件"
	})
	if err != nil {
		panic("获取对象失败: " + err.Error())
	}
	defer res.Body.Close() // V2 中需要记得关闭 Body

	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic("读取对象内容失败: " + err.Error())
	}
	t.Log(string(data))
}
