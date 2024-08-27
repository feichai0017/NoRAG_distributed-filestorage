package config

const (
	// UploadServiceHost : 上传服务监听的地址
	UploadServiceHost = "0.0.0.0:8080"
	// UploadLBHost: 上传服务LB地址
	UploadLBHost = "http://upload.fileserver.com"
	// DownloadLBHost: 下载服务LB地址
	DownloadLBHost = "http://download.fileserver.com"
	// TracerAgentHost: tracing agent地址
	TracerAgentHost = "127.0.0.1:6831"
)
