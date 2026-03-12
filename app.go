package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
	"unsafe"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

// Windows API 函数和常量
var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	dwmapi                    = syscall.NewLazyDLL("dwmapi.dll")
	procFindWindow            = user32.NewProc("FindWindowW")
	procDwmSetWindowAttribute = dwmapi.NewProc("DwmSetWindowAttribute")
)

const (
	DWMWA_WINDOW_CORNER_PREFERENCE = 33
)

// Config 配置项结构体
type Config struct {
	VideoDir    string   `json:"videoDir"`
	Port        int      `json:"port"`
	AllowedExts []string `json:"allowedExts"`
	IsAllowDel  bool     `json:"isAllowDel"`
	IsDeep      bool     `json:"isDeep"`

	// 图片相关配置
	ImageExts     []string `json:"imageExts"`
	IsAllowDelImg bool     `json:"isAllowDelImg"`
}

// 全局列表常量
var (
	videoListCache []string    // 视频列表缓存
	imageListCache []ImageInfo // 图片列表缓存
)

// VideoInfo 视频信息结构体
type VideoInfo struct {
	VideoList    []string `json:"videoList"`
	CurrentVideo string   `json:"currentVideo"`
	HasVideo     bool     `json:"hasVideo"`
	Port         int      `json:"port"`
}

// DeleteResult 删除视频结果结构体
type DeleteResult struct {
	Success      bool   `json:"success"`
	Msg          string `json:"msg"`
	CurrentVideo string `json:"currentVideo"`
	HasVideo     bool   `json:"hasVideo"`
}

// ImageInfo 图片信息结构体
type ImageInfo struct {
	Pic  string `json:"pic"`  // 可访问的URL路径
	PicW int    `json:"picW"` // 图片宽度
	PicH int    `json:"picH"` // 图片高度
	Path string `json:"path"` // 完整文件路径
}

// ImageListResult 图片列表结果结构体
type ImageListResult struct {
	Success   bool        `json:"success"`
	ImageList []ImageInfo `json:"imageList"`
	Msg       string      `json:"msg"`
}

// DeleteImageResult 删除图片结果结构体
type DeleteImageResult struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// App struct
type App struct {
	ctx    context.Context
	config *Config
	server *http.Server
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		config: &Config{
			Port:          3000,
			AllowedExts:   []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv"},
			IsAllowDel:    false,
			IsDeep:        true,
			ImageExts:     []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".JPG", ".JPEG", ".PNG", ".GIF", ".WEBP", ".BMP"},
			IsAllowDelImg: false,
		},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 设置窗口大小和位置
	wailsruntime.WindowSetSize(ctx, 400, 320)
	// wailsruntime.WindowSetPosition(ctx, 100, 100)

	// 延迟设置窗口样式
	go func() {
		time.Sleep(500 * time.Millisecond)
		a.removeWindowCorners() // 移除窗口圆角
	}()
}

// removeWindowCorners 移除窗口圆角
func (a *App) removeWindowCorners() {
	// 将窗口标题转换为UTF-16
	titlePtr, err := syscall.UTF16PtrFromString("视频预览")
	if err != nil {
		return
	}

	// 查找窗口句柄
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return
	}

	// 设置窗口圆角偏好为直角（0 = 直角，1 = 小圆角，2 = 默认圆角）
	procDwmSetWindowAttribute.Call(hwnd, DWMWA_WINDOW_CORNER_PREFERENCE, uintptr(unsafe.Pointer(&[]uint32{0}[0])), 4)
}

// SetConfig 设置配置
//
//export SetConfig
func (a *App) SetConfig(videoDir string, port int, allowedExts []string, imageExts []string, isAllowDel bool, isDeep bool) error {
	a.config.VideoDir = videoDir
	a.config.Port = port
	if len(allowedExts) > 0 {
		a.config.AllowedExts = allowedExts
	}
	if len(imageExts) > 0 {
		a.config.ImageExts = imageExts
	}
	a.config.IsAllowDel = isAllowDel
	a.config.IsDeep = isDeep

	// 确保视频目录存在
	return os.MkdirAll(videoDir, 0755)
}

// GetConfig 获取配置
func (a *App) GetConfig() *Config {
	return a.config
}

// SelectDirectory 选择目录
func (a *App) SelectDirectory() (string, error) {
	selection, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "选择视频目录",
	})
	if err != nil {
		return "", err
	}
	return selection, nil
}

// GetVideoList 获取视频文件列表
func (a *App) GetVideoList() (*VideoInfo, error) {
	// 使用公共扫描方法获取视频文件列表
	videoList, err := a.scanFiles(a.config.VideoDir, a.config.AllowedExts, a.config.IsDeep)
	if err != nil {
		return nil, err
	}

	// 更新全局视频列表缓存
	videoListCache = videoList

	// 确定当前视频
	currentVideo := ""
	if len(videoList) > 0 {
		currentVideo = videoList[0]
	}

	return &VideoInfo{
		VideoList:    videoList,
		CurrentVideo: currentVideo,
		HasVideo:     len(videoList) > 0,
		Port:         a.config.Port,
	}, nil
}

// SwitchVideo 切换视频
func (a *App) SwitchVideo(current string, direction string) (*VideoInfo, error) {
	videoInfo, err := a.GetVideoList()
	if err != nil {
		return nil, err
	}

	if len(videoInfo.VideoList) == 0 {
		return videoInfo, nil
	}

	// 查找当前视频的索引
	currentIndex := -1
	for i, video := range videoInfo.VideoList {
		if video == current {
			currentIndex = i
			break
		}
	}

	// 如果没找到当前视频，使用第一个视频
	if currentIndex == -1 {
		videoInfo.CurrentVideo = videoInfo.VideoList[0]
		return videoInfo, nil
	}

	// 根据方向切换视频
	if direction == "prev" {
		currentIndex = (currentIndex - 1 + len(videoInfo.VideoList)) % len(videoInfo.VideoList)
	} else if direction == "next" {
		currentIndex = (currentIndex + 1) % len(videoInfo.VideoList)
	}

	videoInfo.CurrentVideo = videoInfo.VideoList[currentIndex]
	return videoInfo, nil
}

// DeleteVideo 删除视频文件
func (a *App) DeleteVideo(fileName string) (*DeleteResult, error) {
	if !a.config.IsAllowDel {
		return &DeleteResult{
			Success:      false,
			Msg:          "删除功能已禁用",
			CurrentVideo: "",
			HasVideo:     false,
		}, nil
	}

	if fileName == "" {
		return &DeleteResult{
			Success:      false,
			Msg:          "请指定要删除的文件",
			CurrentVideo: "",
			HasVideo:     false,
		}, nil
	}

	// 构建视频文件路径（支持子文件夹）
	videoPath := filepath.Join(a.config.VideoDir, filepath.FromSlash(fileName))

	// 检查文件是否存在
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return &DeleteResult{
			Success:      false,
			Msg:          "文件不存在",
			CurrentVideo: "",
			HasVideo:     false,
		}, nil
	}

	// 获取删除前的视频列表和当前视频的索引
	videoInfo, err := a.GetVideoList()
	if err != nil {
		return nil, err
	}

	currentIndex := -1
	for i, video := range videoInfo.VideoList {
		if video == fileName {
			currentIndex = i
			break
		}
	}

	// 删除文件
	if err := os.Remove(videoPath); err != nil {
		return &DeleteResult{
			Success:      false,
			Msg:          "删除失败: " + err.Error(),
			CurrentVideo: "",
			HasVideo:     false,
		}, nil
	}

	// 获取删除后的视频列表
	newVideoInfo, err := a.GetVideoList()
	if err != nil {
		return nil, err
	}

	// 确定新的当前视频
	newVideo := ""
	if newVideoInfo.HasVideo {
		// 如果删除的是最后一个，则跳到第一个
		nextIndex := currentIndex
		if nextIndex >= len(newVideoInfo.VideoList) {
			nextIndex = 0
		}
		newVideo = newVideoInfo.VideoList[nextIndex]
	}

	return &DeleteResult{
		Success:      true,
		Msg:          "删除成功",
		CurrentVideo: newVideo,
		HasVideo:     newVideoInfo.HasVideo,
	}, nil
}

// StartServer 启动HTTP服务器
//
//export StartServer
func (a *App) StartServer() int {
	// 如果服务器已经在运行，先停止
	if a.server != nil {
		a.StopServer()
	}

	// 创建自定义的ServeMux
	mux := http.NewServeMux()

	// 创建HTTP路由
	mux.HandleFunc("/", a.handleIndex)
	mux.HandleFunc("/mv/list.json", a.handleVideos)
	mux.HandleFunc("/mv/switch.json", a.handleSwitch)
	mux.HandleFunc("/mv/del.json", a.handleDelete)
	mux.HandleFunc("/img/getImgsList.json", a.handleImages)
	mux.HandleFunc("/img/del.json", a.handleDeleteImage)

	// 配置静态文件服务
	mux.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(a.config.VideoDir))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(a.config.VideoDir))))

	// 从嵌入的文件系统中提供assets
	assetsFS, err := fs.Sub(frontendFS, "frontend/dist/assets")
	if err == nil {
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))
	}

	// 创建服务器
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Port),
		Handler: mux,
	}

	// 在goroutine中启动服务器
	go func() {
		fmt.Printf("服务已启动：http://localhost:%d\n", a.config.Port)
		fmt.Printf("文件目录：%s\n", a.config.VideoDir)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("服务器启动失败: %v\n", err)
		}
	}()

	return a.config.Port
}

// StopServer 停止HTTP服务器
func (a *App) StopServer() {
	if a.server != nil {
		a.server.Close()
		a.server = nil
	}

	// 清空全局列表缓存
	videoListCache = []string{}
	imageListCache = []ImageInfo{}
}

// handleIndex 处理首页请求
func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	// 从嵌入的文件系统中读取index.html
	htmlContent, err := frontendFS.ReadFile("frontend/dist/index.html")
	if err != nil {
		http.Error(w, "无法读取前端文件: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 写入响应
	w.Write(htmlContent)
}

// handleVideos 处理获取视频列表请求
func (a *App) handleVideos(w http.ResponseWriter, r *http.Request) {
	// 调用GetVideoList方法
	videoInfo, err := a.GetVideoList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用 json.Marshal 生成正确的 JSON 格式
	videoListJSON, _ := json.Marshal(videoInfo.VideoList)

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应
	w.Write([]byte(fmt.Sprintf(`{
		"success": true,
		"videoList": %s,
		"currentVideo": "%s",
		"hasVideo": %t,
		"port": %d
	}`, string(videoListJSON), videoInfo.CurrentVideo, videoInfo.HasVideo, videoInfo.Port)))
}

// handleSwitch 处理切换视频请求
func (a *App) handleSwitch(w http.ResponseWriter, r *http.Request) {
	current := r.URL.Query().Get("current")
	direction := r.URL.Query().Get("direction")

	// 调用SwitchVideo方法
	videoInfo, err := a.SwitchVideo(current, direction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用 json.Marshal 生成正确的 JSON 格式
	videoListJSON, _ := json.Marshal(videoInfo.VideoList)

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应
	w.Write([]byte(fmt.Sprintf(`{
		"success": true,
		"videoList": %s,
		"currentVideo": "%s",
		"hasVideo": %t
	}`, string(videoListJSON), videoInfo.CurrentVideo, videoInfo.HasVideo)))
}

// handleDelete 处理删除视频请求
func (a *App) handleDelete(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")

	// 调用DeleteVideo方法
	result, err := a.DeleteVideo(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应
	w.Write([]byte(fmt.Sprintf(`{
		"success": %t,
		"msg": "%s",
		"currentVideo": "%s",
		"hasVideo": %t
	}`, result.Success, result.Msg, result.CurrentVideo, result.HasVideo)))
}

// handleImages 处理获取图片列表请求
func (a *App) handleImages(w http.ResponseWriter, r *http.Request) {
	// 调用GetImageList方法
	result, err := a.GetImageList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应
	imageListJSON, _ := json.Marshal(result.ImageList)
	w.Write([]byte(fmt.Sprintf(`{
		"success": %t,
		"imageList": %s,
		"msg": "%s"
	}`, result.Success, string(imageListJSON), result.Msg)))
}

// handleDeleteImage 处理删除图片请求
func (a *App) handleDeleteImage(w http.ResponseWriter, r *http.Request) {
	imagePath := r.URL.Query().Get("path")

	// 调用DeleteImage方法
	result, err := a.DeleteImage(imagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应
	w.Write([]byte(fmt.Sprintf(`{
		"success": %t,
		"msg": "%s"
	}`, result.Success, result.Msg)))
}

// OpenBrowser 打开浏览器
func (a *App) OpenBrowser() error {
	url := fmt.Sprintf("http://localhost:%d/mv", a.config.Port)
	exec.Command("cmd", "/c", "start", url).Run()
	urlTest := fmt.Sprintf("http://localhost:%d/img", a.config.Port)
	return exec.Command("cmd", "/c", "start", urlTest).Run()
}

// SelectFileOrDir 打开文件或目录选择对话框
func (a *App) SelectFileOrDir(title string, isDir bool) (string, error) {
	// 这里使用简单的方法，实际项目中可能需要使用更复杂的对话框
	// 由于Wails框架的限制，我们直接返回一个空字符串，让前端使用原生的文件选择
	return "", nil
}

// SelectImageDirectory 选择图片目录
//
//export SelectImageDirectory
func (a *App) SelectImageDirectory() (string, error) {
	selection, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "选择图片目录",
	})
	if err != nil {
		return "", err
	}
	return selection, nil
}

// scanFiles 公共文件扫描方法
func (a *App) scanFiles(dir string, allowedExts []string, isDeep bool) ([]string, error) {
	if dir == "" {
		return []string{}, nil
	}

	var fileList []string
	var err error

	if isDeep {
		// 深度扫描：递归扫描所有子文件夹
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // 跳过无法访问的文件
			}

			// 跳过目录
			if info.IsDir() {
				return nil
			}

			// 检查文件扩展名
			ext := strings.ToLower(filepath.Ext(path))
			for _, allowedExt := range allowedExts {
				if ext == strings.ToLower(allowedExt) {
					// 计算相对路径
					relPath, err := filepath.Rel(dir, path)
					if err != nil {
						return nil // 跳过无法计算相对路径的文件
					}
					// 将路径分隔符统一为 /
					relPath = filepath.ToSlash(relPath)
					fileList = append(fileList, relPath)
					break
				}
			}
			return nil
		})
	} else {
		// 浅层扫描：只扫描当前目录下的文件
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			// 跳过目录
			if entry.IsDir() {
				continue
			}

			// 检查文件扩展名
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			for _, allowedExt := range allowedExts {
				if ext == strings.ToLower(allowedExt) {
					fileList = append(fileList, entry.Name())
					break
				}
			}
		}
	}

	if err != nil {
		return nil, err
	}

	// 对文件列表进行排序（按文件名字母顺序）
	sort.Strings(fileList)

	return fileList, nil
}

// GetImageList 获取图片文件列表
//
//export GetImageList
func (a *App) GetImageList() (*ImageListResult, error) {
	if a.config.VideoDir == "" {
		return &ImageListResult{
			Success:   true,
			ImageList: []ImageInfo{},
			Msg:       "请先选择目录",
		}, nil
	}

	// 使用公共扫描方法获取图片文件列表
	imageFiles, err := a.scanFiles(a.config.VideoDir, a.config.ImageExts, a.config.IsDeep)
	if err != nil {
		return &ImageListResult{
			Success:   false,
			ImageList: []ImageInfo{},
			Msg:       "扫描图片失败: " + err.Error(),
		}, nil
	}

	// 构建图片信息列表
	imageList := []ImageInfo{}

	for _, imageFile := range imageFiles {
		fullPath := filepath.Join(a.config.VideoDir, filepath.FromSlash(imageFile))

		// 获取图片尺寸
		width, height, err := a.getImageDimensions(fullPath)
		if err != nil {
			// 如果获取尺寸失败，使用默认值
			width = 0
			height = 0
		}

		// 计算相对路径（URL路径），前端会拼接API_PREFIX
		urlPath := "/images/" + imageFile

		imageList = append(imageList, ImageInfo{
			Pic:  urlPath,
			PicW: width,
			PicH: height,
			Path: fullPath,
		})
	}

	// 更新全局图片列表缓存
	imageListCache = imageList

	return &ImageListResult{
		Success:   true,
		ImageList: imageList,
		Msg:       fmt.Sprintf("共找到 %d 张图片", len(imageList)),
	}, nil
}

// getImageDimensions 获取图片尺寸
func (a *App) getImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// 使用image包解码图片配置，不需要解码整个图片
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}

// DeleteImage 删除图片文件
//
//export DeleteImage
func (a *App) DeleteImage(imagePath string) (*DeleteImageResult, error) {
	if !a.config.IsAllowDelImg {
		return &DeleteImageResult{
			Success: false,
			Msg:     "删除功能已禁用",
		}, nil
	}

	if imagePath == "" {
		return &DeleteImageResult{
			Success: false,
			Msg:     "请指定要删除的文件",
		}, nil
	}

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return &DeleteImageResult{
			Success: false,
			Msg:     "文件不存在",
		}, nil
	}

	// 删除文件
	if err := os.Remove(imagePath); err != nil {
		return &DeleteImageResult{
			Success: false,
			Msg:     "删除失败: " + err.Error(),
		}, nil
	}

	return &DeleteImageResult{
		Success: true,
		Msg:     "删除成功",
	}, nil
}
