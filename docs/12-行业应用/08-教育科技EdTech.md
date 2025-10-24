# Go在教育科技（EdTech）中的应用

> **简介**: 系统介绍Go语言在在线教育、学习管理系统、教学互动平台等教育科技领域的架构设计、技术实践与工程落地

---

## 📚 目录

- [Go在教育科技（EdTech）中的应用](#go在教育科技edtech中的应用)
  - [📚 目录](#-目录)
  - [1. 教育科技概览](#1-教育科技概览)
    - [1.1 行业特点](#11-行业特点)
    - [1.2 Go的优势](#12-go的优势)
  - [2. 在线学习平台架构](#2-在线学习平台架构)
    - [2.1 整体架构](#21-整体架构)
    - [2.2 微服务划分](#22-微服务划分)
  - [3. 学习管理系统（LMS）](#3-学习管理系统lms)
    - [3.1 课程管理](#31-课程管理)
    - [3.2 学习进度跟踪](#32-学习进度跟踪)
  - [4. 实时互动教学](#4-实时互动教学)
    - [4.1 在线问答系统](#41-在线问答系统)
    - [4.2 实时白板](#42-实时白板)
  - [5. 课程内容管理](#5-课程内容管理)
    - [5.1 视频处理](#51-视频处理)
  - [6. 学习数据分析](#6-学习数据分析)
    - [6.1 学习行为分析](#61-学习行为分析)
  - [7. 考试评测系统](#7-考试评测系统)
    - [7.1 在线考试](#71-在线考试)
  - [8. 视频直播与点播](#8-视频直播与点播)
    - [8.1 直播系统](#81-直播系统)
  - [9. 完整项目：在线学习平台](#9-完整项目在线学习平台)
    - [9.1 项目结构](#91-项目结构)
    - [9.2 核心API实现](#92-核心api实现)
  - [💡 总结](#-总结)
    - [核心要点](#核心要点)
    - [进阶方向](#进阶方向)
  - [🔗 相关资源](#-相关资源)

---

## 1. 教育科技概览

### 1.1 行业特点

**核心需求**:

- 高并发访问（课程直播、作业提交）
- 低延迟互动（在线问答、实时批注）
- 大规模存储（视频、文档、作业）
- 数据安全（学生信息、成绩记录）
- 个性化学习（推荐算法、学习路径）

**技术挑战**:

- 视频流媒体处理
- 大量用户并发学习
- 实时互动体验
- 学习数据分析
- 考试防作弊

### 1.2 Go的优势

```go
// Go在EdTech中的优势
优势特性:
✅ 高并发处理 - 支持大量学生同时在线
✅ 高性能 - 快速响应用户请求
✅ 易于部署 - 单一二进制文件，便于分发
✅ 丰富生态 - 完善的Web框架和工具链
✅ 云原生 - 天然适合微服务和容器化
```

---

## 2. 在线学习平台架构

### 2.1 整体架构

```go
package architecture

/*
┌─────────────────────────────────────────────────────────────┐
│                      前端应用层                               │
│  Web端      移动端      小程序      TV端                       │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      API 网关层                               │
│  认证授权    路由分发    限流降级    日志监控                   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      业务服务层                               │
│  用户服务   课程服务   学习服务   考试服务   互动服务            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      基础设施层                               │
│  数据库     缓存      消息队列    对象存储    CDN              │
└─────────────────────────────────────────────────────────────┘
*/
```

### 2.2 微服务划分

```go
package services

// UserService 用户服务
type UserService struct {
    repo UserRepository
}

// 用户注册
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
    // 验证手机号/邮箱
    if err := s.validateContact(req.Contact); err != nil {
        return nil, err
    }

    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(req.Password), 
        bcrypt.DefaultCost,
    )
    if err != nil {
        return nil, err
    }

    user := &User{
        ID:           generateID(),
        Username:     req.Username,
        Contact:      req.Contact,
        PasswordHash: string(hashedPassword),
        Role:         RoleStudent,
        Status:       StatusActive,
        CreatedAt:    time.Now(),
    }

    return s.repo.Create(ctx, user)
}

// CourseService 课程服务
type CourseService struct {
    repo    CourseRepository
    storage StorageClient
}

// 创建课程
func (s *CourseService) CreateCourse(ctx context.Context, req CreateCourseRequest) (*Course, error) {
    course := &Course{
        ID:          generateID(),
        Title:       req.Title,
        Description: req.Description,
        TeacherID:   req.TeacherID,
        Category:    req.Category,
        Price:       req.Price,
        Status:      StatusDraft,
        CreatedAt:   time.Now(),
    }

    // 保存课程
    if err := s.repo.Create(ctx, course); err != nil {
        return nil, err
    }

    // 上传封面图
    if req.CoverImage != nil {
        coverURL, err := s.storage.Upload(ctx, req.CoverImage)
        if err != nil {
            return nil, err
        }
        course.CoverURL = coverURL
        s.repo.Update(ctx, course)
    }

    return course, nil
}

// LearningService 学习服务
type LearningService struct {
    repo     LearningRepository
    progress ProgressTracker
}

// 开始学习
func (s *LearningService) StartLearning(ctx context.Context, userID, courseID string) error {
    // 检查用户是否已购买课程
    enrolled, err := s.repo.IsEnrolled(ctx, userID, courseID)
    if err != nil {
        return err
    }
    if !enrolled {
        return ErrNotEnrolled
    }

    // 创建学习记录
    record := &LearningRecord{
        UserID:    userID,
        CourseID:  courseID,
        StartTime: time.Now(),
        Progress:  0,
        Status:    StatusLearning,
    }

    return s.repo.CreateRecord(ctx, record)
}

// 更新学习进度
func (s *LearningService) UpdateProgress(
    ctx context.Context,
    userID, courseID, lessonID string,
    progress int,
) error {
    return s.progress.Update(ctx, &ProgressUpdate{
        UserID:   userID,
        CourseID: courseID,
        LessonID: lessonID,
        Progress: progress,
        UpdateAt: time.Now(),
    })
}
```

---

## 3. 学习管理系统（LMS）

### 3.1 课程管理

```go
package lms

import (
    "context"
    "time"
)

// Course 课程模型
type Course struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    TeacherID   string    `json:"teacher_id"`
    Category    string    `json:"category"`
    Level       string    `json:"level"` // 初级/中级/高级
    Price       int64     `json:"price"` // 价格（分）
    Status      string    `json:"status"`
    CoverURL    string    `json:"cover_url"`
    Duration    int       `json:"duration"` // 总时长（分钟）
    Chapters    []Chapter `json:"chapters"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Chapter 章节
type Chapter struct {
    ID       string   `json:"id"`
    Title    string   `json:"title"`
    Order    int      `json:"order"`
    Lessons  []Lesson `json:"lessons"`
}

// Lesson 课时
type Lesson struct {
    ID        string        `json:"id"`
    Title     string        `json:"title"`
    Type      string        `json:"type"` // video/doc/quiz
    Content   LessonContent `json:"content"`
    Duration  int           `json:"duration"`
    Order     int           `json:"order"`
    IsFree    bool          `json:"is_free"` // 是否试听
}

// LessonContent 课时内容
type LessonContent struct {
    VideoURL     string   `json:"video_url,omitempty"`
    DocumentURL  string   `json:"document_url,omitempty"`
    QuizID       string   `json:"quiz_id,omitempty"`
    Attachments  []string `json:"attachments,omitempty"`
}

// CourseManager 课程管理器
type CourseManager struct {
    repo    CourseRepository
    storage StorageClient
    cache   CacheClient
}

// GetCourseDetail 获取课程详情
func (m *CourseManager) GetCourseDetail(ctx context.Context, courseID string) (*Course, error) {
    // 先从缓存获取
    if course, err := m.cache.GetCourse(ctx, courseID); err == nil {
        return course, nil
    }

    // 从数据库获取
    course, err := m.repo.GetByID(ctx, courseID)
    if err != nil {
        return nil, err
    }

    // 加载章节和课时
    chapters, err := m.repo.GetChapters(ctx, courseID)
    if err != nil {
        return nil, err
    }

    for i := range chapters {
        lessons, err := m.repo.GetLessons(ctx, chapters[i].ID)
        if err != nil {
            return nil, err
        }
        chapters[i].Lessons = lessons
    }

    course.Chapters = chapters

    // 缓存课程信息
    m.cache.SetCourse(ctx, course, 1*time.Hour)

    return course, nil
}

// PublishCourse 发布课程
func (m *CourseManager) PublishCourse(ctx context.Context, courseID string) error {
    course, err := m.repo.GetByID(ctx, courseID)
    if err != nil {
        return err
    }

    // 验证课程内容完整性
    if err := m.validateCourse(course); err != nil {
        return err
    }

    // 更新状态为已发布
    course.Status = StatusPublished
    course.UpdatedAt = time.Now()

    if err := m.repo.Update(ctx, course); err != nil {
        return err
    }

    // 清除缓存
    m.cache.DeleteCourse(ctx, courseID)

    return nil
}

func (m *CourseManager) validateCourse(course *Course) error {
    if course.Title == "" {
        return ErrMissingTitle
    }
    if course.CoverURL == "" {
        return ErrMissingCover
    }
    if len(course.Chapters) == 0 {
        return ErrNoChapters
    }
    return nil
}
```

### 3.2 学习进度跟踪

```go
package tracking

import (
    "context"
    "sync"
    "time"
)

// ProgressTracker 进度跟踪器
type ProgressTracker struct {
    repo  ProgressRepository
    mu    sync.RWMutex
    cache map[string]*Progress // userID:courseID -> Progress
}

// Progress 学习进度
type Progress struct {
    UserID          string            `json:"user_id"`
    CourseID        string            `json:"course_id"`
    TotalLessons    int               `json:"total_lessons"`
    CompletedLessons int              `json:"completed_lessons"`
    Progress        int               `json:"progress"` // 百分比
    LastLesson      string            `json:"last_lesson"`
    LessonProgress  map[string]int    `json:"lesson_progress"` // lessonID -> 进度
    TotalDuration   int               `json:"total_duration"` // 总时长（秒）
    LearnedDuration int               `json:"learned_duration"` // 已学时长（秒）
    UpdatedAt       time.Time         `json:"updated_at"`
}

func NewProgressTracker(repo ProgressRepository) *ProgressTracker {
    return &ProgressTracker{
        repo:  repo,
        cache: make(map[string]*Progress),
    }
}

// GetProgress 获取学习进度
func (t *ProgressTracker) GetProgress(ctx context.Context, userID, courseID string) (*Progress, error) {
    key := userID + ":" + courseID

    // 先从内存缓存获取
    t.mu.RLock()
    if progress, exists := t.cache[key]; exists {
        t.mu.RUnlock()
        return progress, nil
    }
    t.mu.RUnlock()

    // 从数据库获取
    progress, err := t.repo.GetProgress(ctx, userID, courseID)
    if err != nil {
        return nil, err
    }

    // 缓存到内存
    t.mu.Lock()
    t.cache[key] = progress
    t.mu.Unlock()

    return progress, nil
}

// UpdateLessonProgress 更新课时进度
func (t *ProgressTracker) UpdateLessonProgress(
    ctx context.Context,
    userID, courseID, lessonID string,
    duration int, // 当前观看时长
) error {
    progress, err := t.GetProgress(ctx, userID, courseID)
    if err != nil {
        return err
    }

    // 更新课时进度
    if progress.LessonProgress == nil {
        progress.LessonProgress = make(map[string]int)
    }
    progress.LessonProgress[lessonID] = duration
    progress.LastLesson = lessonID
    progress.LearnedDuration += duration

    // 检查课时是否完成（观看超过90%）
    if t.isLessonCompleted(lessonID, duration) {
        progress.CompletedLessons++
    }

    // 计算整体进度
    progress.Progress = (progress.CompletedLessons * 100) / progress.TotalLessons
    progress.UpdatedAt = time.Now()

    // 保存到数据库
    if err := t.repo.UpdateProgress(ctx, progress); err != nil {
        return err
    }

    // 更新缓存
    key := userID + ":" + courseID
    t.mu.Lock()
    t.cache[key] = progress
    t.mu.Unlock()

    return nil
}

func (t *ProgressTracker) isLessonCompleted(lessonID string, watchedDuration int) bool {
    // 获取课时总时长
    lessonDuration := t.getLessonDuration(lessonID)
    if lessonDuration == 0 {
        return false
    }

    // 判断是否观看超过90%
    return float64(watchedDuration)/float64(lessonDuration) >= 0.9
}

func (t *ProgressTracker) getLessonDuration(lessonID string) int {
    // 从数据库或缓存获取课时时长
    // 这里简化处理
    return 600 // 假设10分钟
}
```

---

## 4. 实时互动教学

### 4.1 在线问答系统

```go
package interactive

import (
    "context"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

// QuestionAnswer 问答系统
type QuestionAnswer struct {
    clients   map[string]*Client // userID -> Client
    questions chan *Question
    answers   chan *Answer
    mu        sync.RWMutex
}

// Client WebSocket客户端
type Client struct {
    ID       string
    UserID   string
    CourseID string
    Role     string // student/teacher
    Conn     *websocket.Conn
    Send     chan []byte
}

// Question 问题
type Question struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Username  string    `json:"username"`
    CourseID  string    `json:"course_id"`
    Content   string    `json:"content"`
    Images    []string  `json:"images,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

// Answer 答案
type Answer struct {
    ID         string    `json:"id"`
    QuestionID string    `json:"question_id"`
    UserID     string    `json:"user_id"`
    Username   string    `json:"username"`
    Content    string    `json:"content"`
    IsTeacher  bool      `json:"is_teacher"`
    CreatedAt  time.Time `json:"created_at"`
}

func NewQuestionAnswer() *QuestionAnswer {
    return &QuestionAnswer{
        clients:   make(map[string]*Client),
        questions: make(chan *Question, 100),
        answers:   make(chan *Answer, 100),
    }
}

// Run 运行问答系统
func (qa *QuestionAnswer) Run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case question := <-qa.questions:
            qa.broadcastQuestion(question)
        case answer := <-qa.answers:
            qa.broadcastAnswer(answer)
        }
    }
}

// HandleClient 处理客户端连接
func (qa *QuestionAnswer) HandleClient(client *Client) {
    // 注册客户端
    qa.mu.Lock()
    qa.clients[client.ID] = client
    qa.mu.Unlock()

    defer func() {
        // 移除客户端
        qa.mu.Lock()
        delete(qa.clients, client.ID)
        qa.mu.Unlock()
        client.Conn.Close()
    }()

    // 读取消息
    for {
        var msg Message
        err := client.Conn.ReadJSON(&msg)
        if err != nil {
            break
        }

        switch msg.Type {
        case "question":
            qa.handleQuestion(client, msg.Data)
        case "answer":
            qa.handleAnswer(client, msg.Data)
        }
    }
}

func (qa *QuestionAnswer) handleQuestion(client *Client, data interface{}) {
    question := &Question{
        ID:        generateID(),
        UserID:    client.UserID,
        CourseID:  client.CourseID,
        CreatedAt: time.Now(),
    }

    // 解析问题内容
    // ...

    qa.questions <- question
}

func (qa *QuestionAnswer) handleAnswer(client *Client, data interface{}) {
    answer := &Answer{
        ID:        generateID(),
        UserID:    client.UserID,
        IsTeacher: client.Role == "teacher",
        CreatedAt: time.Now(),
    }

    // 解析答案内容
    // ...

    qa.answers <- answer
}

func (qa *QuestionAnswer) broadcastQuestion(question *Question) {
    qa.mu.RLock()
    defer qa.mu.RUnlock()

    for _, client := range qa.clients {
        if client.CourseID == question.CourseID {
            client.Send <- encodeMessage("question", question)
        }
    }
}

func (qa *QuestionAnswer) broadcastAnswer(answer *Answer) {
    qa.mu.RLock()
    defer qa.mu.RUnlock()

    for _, client := range qa.clients {
        client.Send <- encodeMessage("answer", answer)
    }
}
```

### 4.2 实时白板

```go
package whiteboard

import (
    "sync"
    "time"
)

// Whiteboard 白板
type Whiteboard struct {
    ID        string
    CourseID  string
    Elements  []*Element
    History   []*Operation
    mu        sync.RWMutex
    clients   map[string]*Client
}

// Element 白板元素
type Element struct {
    ID      string      `json:"id"`
    Type    string      `json:"type"` // line/rect/circle/text/image
    Data    interface{} `json:"data"`
    Style   Style       `json:"style"`
    ZIndex  int         `json:"z_index"`
}

// Style 样式
type Style struct {
    Color       string  `json:"color"`
    StrokeWidth float64 `json:"stroke_width"`
    Fill        string  `json:"fill"`
    FontSize    int     `json:"font_size,omitempty"`
}

// Operation 操作记录
type Operation struct {
    Type      string    `json:"type"` // add/update/delete
    ElementID string    `json:"element_id"`
    Element   *Element  `json:"element,omitempty"`
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
}

// AddElement 添加元素
func (w *Whiteboard) AddElement(elem *Element) {
    w.mu.Lock()
    defer w.mu.Unlock()

    w.Elements = append(w.Elements, elem)

    // 记录操作
    op := &Operation{
        Type:      "add",
        ElementID: elem.ID,
        Element:   elem,
        Timestamp: time.Now(),
    }
    w.History = append(w.History, op)

    // 广播给所有客户端
    w.broadcast(op)
}

// UpdateElement 更新元素
func (w *Whiteboard) UpdateElement(elemID string, elem *Element) error {
    w.mu.Lock()
    defer w.mu.Unlock()

    for i, e := range w.Elements {
        if e.ID == elemID {
            w.Elements[i] = elem

            op := &Operation{
                Type:      "update",
                ElementID: elemID,
                Element:   elem,
                Timestamp: time.Now(),
            }
            w.History = append(w.History, op)
            w.broadcast(op)

            return nil
        }
    }

    return ErrElementNotFound
}

// DeleteElement 删除元素
func (w *Whiteboard) DeleteElement(elemID string) error {
    w.mu.Lock()
    defer w.mu.Unlock()

    for i, e := range w.Elements {
        if e.ID == elemID {
            w.Elements = append(w.Elements[:i], w.Elements[i+1:]...)

            op := &Operation{
                Type:      "delete",
                ElementID: elemID,
                Timestamp: time.Now(),
            }
            w.History = append(w.History, op)
            w.broadcast(op)

            return nil
        }
    }

    return ErrElementNotFound
}

func (w *Whiteboard) broadcast(op *Operation) {
    for _, client := range w.clients {
        client.Send <- encodeMessage("operation", op)
    }
}

// Clear 清空白板
func (w *Whiteboard) Clear() {
    w.mu.Lock()
    defer w.mu.Unlock()

    w.Elements = []*Element{}

    op := &Operation{
        Type:      "clear",
        Timestamp: time.Now(),
    }
    w.History = append(w.History, op)
    w.broadcast(op)
}
```

---

## 5. 课程内容管理

### 5.1 视频处理

```go
package content

import (
    "context"
    "fmt"
    "os/exec"
)

// VideoProcessor 视频处理器
type VideoProcessor struct {
    storage    StorageClient
    transcoder TranscoderClient
}

// ProcessVideo 处理视频
func (p *VideoProcessor) ProcessVideo(ctx context.Context, videoPath string) (*VideoResult, error) {
    // 1. 获取视频信息
    info, err := p.getVideoInfo(videoPath)
    if err != nil {
        return nil, err
    }

    // 2. 转码为多种清晰度
    qualities := []Quality{Quality360p, Quality480p, Quality720p, Quality1080p}
    var outputs []VideoOutput

    for _, quality := range qualities {
        output, err := p.transcode(ctx, videoPath, quality)
        if err != nil {
            continue // 某个清晰度失败不影响其他
        }
        outputs = append(outputs, output)
    }

    // 3. 生成缩略图
    thumbnail, err := p.generateThumbnail(videoPath)
    if err != nil {
        return nil, err
    }

    // 4. 上传到云存储
    result := &VideoResult{
        VideoID:   generateID(),
        Duration:  info.Duration,
        Size:      info.Size,
        Outputs:   outputs,
        Thumbnail: thumbnail,
    }

    for i := range outputs {
        url, err := p.storage.Upload(ctx, outputs[i].Path)
        if err != nil {
            return nil, err
        }
        result.Outputs[i].URL = url
    }

    thumbnailURL, err := p.storage.Upload(ctx, thumbnail)
    if err != nil {
        return nil, err
    }
    result.Thumbnail = thumbnailURL

    return result, nil
}

// VideoInfo 视频信息
type VideoInfo struct {
    Duration int64  // 时长（秒）
    Width    int    // 宽度
    Height   int    // 高度
    Size     int64  // 文件大小（字节）
    Bitrate  int    // 码率
    Format   string // 格式
}

func (p *VideoProcessor) getVideoInfo(videoPath string) (*VideoInfo, error) {
    // 使用 ffprobe 获取视频信息
    cmd := exec.Command("ffprobe",
        "-v", "error",
        "-show_entries", "format=duration,size:stream=width,height,bit_rate",
        "-of", "json",
        videoPath,
    )

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    // 解析 JSON 输出
    var info VideoInfo
    // ... 解析逻辑

    return &info, nil
}

// Quality 视频清晰度
type Quality string

const (
    Quality360p  Quality = "360p"
    Quality480p  Quality = "480p"
    Quality720p  Quality = "720p"
    Quality1080p Quality = "1080p"
)

// VideoOutput 视频输出
type VideoOutput struct {
    Quality Quality `json:"quality"`
    Path    string  `json:"path"`
    URL     string  `json:"url"`
    Size    int64   `json:"size"`
    Bitrate int     `json:"bitrate"`
}

func (p *VideoProcessor) transcode(ctx context.Context, input string, quality Quality) (VideoOutput, error) {
    output := VideoOutput{
        Quality: quality,
    }

    // 根据清晰度设置参数
    var width, height, bitrate int
    switch quality {
    case Quality360p:
        width, height, bitrate = 640, 360, 800
    case Quality480p:
        width, height, bitrate = 854, 480, 1200
    case Quality720p:
        width, height, bitrate = 1280, 720, 2500
    case Quality1080p:
        width, height, bitrate = 1920, 1080, 5000
    }

    outputPath := fmt.Sprintf("%s_%s.mp4", input, quality)

    // 使用 ffmpeg 转码
    cmd := exec.CommandContext(ctx, "ffmpeg",
        "-i", input,
        "-vf", fmt.Sprintf("scale=%d:%d", width, height),
        "-b:v", fmt.Sprintf("%dk", bitrate),
        "-c:v", "libx264",
        "-preset", "fast",
        "-c:a", "aac",
        "-b:a", "128k",
        outputPath,
    )

    if err := cmd.Run(); err != nil {
        return output, err
    }

    output.Path = outputPath
    return output, nil
}

func (p *VideoProcessor) generateThumbnail(videoPath string) (string, error) {
    thumbnailPath := videoPath + "_thumbnail.jpg"

    // 截取第1秒的帧作为缩略图
    cmd := exec.Command("ffmpeg",
        "-i", videoPath,
        "-ss", "00:00:01",
        "-vframes", "1",
        "-vf", "scale=320:180",
        thumbnailPath,
    )

    if err := cmd.Run(); err != nil {
        return "", err
    }

    return thumbnailPath, nil
}

// VideoResult 视频处理结果
type VideoResult struct {
    VideoID   string        `json:"video_id"`
    Duration  int64         `json:"duration"`
    Size      int64         `json:"size"`
    Outputs   []VideoOutput `json:"outputs"`
    Thumbnail string        `json:"thumbnail"`
}
```

---

## 6. 学习数据分析

### 6.1 学习行为分析

```go
package analytics

import (
    "context"
    "time"
)

// LearningAnalytics 学习分析
type LearningAnalytics struct {
    repo AnalyticsRepository
}

// UserLearningStats 用户学习统计
type UserLearningStats struct {
    UserID          string        `json:"user_id"`
    TotalCourses    int           `json:"total_courses"`
    CompletedCourses int          `json:"completed_courses"`
    TotalDuration   int           `json:"total_duration"` // 总学习时长（分钟）
    AvgDailyTime    int           `json:"avg_daily_time"` // 平均每日学习时长
    LongestStreak   int           `json:"longest_streak"` // 最长连续学习天数
    CurrentStreak   int           `json:"current_streak"` // 当前连续学习天数
    WeakDays        []Weekday     `json:"weak_days"` // 学习薄弱日
    PreferredTime   []TimeSlot    `json:"preferred_time"` // 偏好学习时段
    Categories      []CategoryStat `json:"categories"` // 分类统计
}

// CategoryStat 分类统计
type CategoryStat struct {
    Category string `json:"category"`
    Courses  int    `json:"courses"`
    Duration int    `json:"duration"`
    Progress int    `json:"progress"`
}

// GetUserStats 获取用户统计
func (a *LearningAnalytics) GetUserStats(ctx context.Context, userID string, period Period) (*UserLearningStats, error) {
    // 获取用户学习记录
    records, err := a.repo.GetLearningRecords(ctx, userID, period)
    if err != nil {
        return nil, err
    }

    stats := &UserLearningStats{
        UserID: userID,
    }

    // 统计总课程数和完成数
    courseMap := make(map[string]bool)
    completedMap := make(map[string]bool)
    totalDuration := 0

    for _, record := range records {
        courseMap[record.CourseID] = true
        totalDuration += record.Duration

        if record.Progress >= 100 {
            completedMap[record.CourseID] = true
        }
    }

    stats.TotalCourses = len(courseMap)
    stats.CompletedCourses = len(completedMap)
    stats.TotalDuration = totalDuration

    // 计算平均每日学习时长
    days := period.Days()
    if days > 0 {
        stats.AvgDailyTime = totalDuration / days
    }

    // 计算学习连续性
    stats.LongestStreak, stats.CurrentStreak = a.calculateStreak(records)

    // 分析学习时段偏好
    stats.PreferredTime = a.analyzeTimePreference(records)

    // 分析薄弱学习日
    stats.WeakDays = a.analyzeWeakDays(records)

    // 分类统计
    stats.Categories = a.analyzeCategoryStats(records)

    return stats, nil
}

func (a *LearningAnalytics) calculateStreak(records []LearningRecord) (longest, current int) {
    if len(records) == 0 {
        return 0, 0
    }

    // 按日期分组
    dateMap := make(map[string]bool)
    for _, record := range records {
        date := record.CreatedAt.Format("2006-01-02")
        dateMap[date] = true
    }

    // 计算连续天数
    today := time.Now()
    currentStreak := 0
    longestStreak := 0
    tempStreak := 0

    for i := 0; i < 365; i++ {
        date := today.AddDate(0, 0, -i).Format("2006-01-02")
        if dateMap[date] {
            tempStreak++
            if i == 0 || currentStreak > 0 {
                currentStreak++
            }
        } else {
            if tempStreak > longestStreak {
                longestStreak = tempStreak
            }
            tempStreak = 0
        }
    }

    return longestStreak, currentStreak
}

// TimeSlot 时间段
type TimeSlot struct {
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`
    Count int       `json:"count"`
}

func (a *LearningAnalytics) analyzeTimePreference(records []LearningRecord) []TimeSlot {
    // 统计各时段学习次数
    slots := make(map[int]int) // hour -> count

    for _, record := range records {
        hour := record.CreatedAt.Hour()
        slots[hour]++
    }

    // 找出学习频率最高的时段
    var preferred []TimeSlot
    for hour, count := range slots {
        if count > 5 { // 阈值：超过5次
            slot := TimeSlot{
                Start: time.Date(0, 1, 1, hour, 0, 0, 0, time.UTC),
                End:   time.Date(0, 1, 1, hour+1, 0, 0, 0, time.UTC),
                Count: count,
            }
            preferred = append(preferred, slot)
        }
    }

    return preferred
}

// Weekday 星期
type Weekday string

const (
    Monday    Weekday = "Monday"
    Tuesday   Weekday = "Tuesday"
    Wednesday Weekday = "Wednesday"
    Thursday  Weekday = "Thursday"
    Friday    Weekday = "Friday"
    Saturday  Weekday = "Saturday"
    Sunday    Weekday = "Sunday"
)

func (a *LearningAnalytics) analyzeWeakDays(records []LearningRecord) []Weekday {
    // 统计各星期的学习次数
    weekdayCount := make(map[time.Weekday]int)

    for _, record := range records {
        weekday := record.CreatedAt.Weekday()
        weekdayCount[weekday]++
    }

    // 找出学习次数最少的星期
    avgCount := len(records) / 7
    var weakDays []Weekday

    weekdayNames := map[time.Weekday]Weekday{
        time.Monday:    Monday,
        time.Tuesday:   Tuesday,
        time.Wednesday: Wednesday,
        time.Thursday:  Thursday,
        time.Friday:    Friday,
        time.Saturday:  Saturday,
        time.Sunday:    Sunday,
    }

    for wd := time.Monday; wd <= time.Sunday; wd++ {
        if weekdayCount[wd] < avgCount {
            weakDays = append(weakDays, weekdayNames[wd])
        }
    }

    return weakDays
}

func (a *LearningAnalytics) analyzeCategoryStats(records []LearningRecord) []CategoryStat {
    // 按分类统计
    categoryMap := make(map[string]*CategoryStat)

    for _, record := range records {
        stat, exists := categoryMap[record.Category]
        if !exists {
            stat = &CategoryStat{
                Category: record.Category,
            }
            categoryMap[record.Category] = stat
        }

        stat.Courses++
        stat.Duration += record.Duration
        stat.Progress += record.Progress
    }

    // 计算平均进度
    var stats []CategoryStat
    for _, stat := range categoryMap {
        stat.Progress /= stat.Courses
        stats = append(stats, *stat)
    }

    return stats
}

// LearningRecord 学习记录
type LearningRecord struct {
    UserID    string
    CourseID  string
    Category  string
    Duration  int
    Progress  int
    CreatedAt time.Time
}

// Period 时间段
type Period struct {
    Start time.Time
    End   time.Time
}

func (p Period) Days() int {
    return int(p.End.Sub(p.Start).Hours() / 24)
}
```

---

## 7. 考试评测系统

### 7.1 在线考试

```go
package exam

import (
    "context"
    "time"
)

// Exam 考试
type Exam struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    CourseID    string     `json:"course_id"`
    Duration    int        `json:"duration"` // 考试时长（分钟）
    TotalScore  int        `json:"total_score"`
    PassScore   int        `json:"pass_score"` // 及格分数
    Questions   []Question `json:"questions"`
    StartTime   time.Time  `json:"start_time"`
    EndTime     time.Time  `json:"end_time"`
    Status      string     `json:"status"` // draft/published/closed
}

// Question 试题
type Question struct {
    ID      string        `json:"id"`
    Type    string        `json:"type"` // single/multiple/judge/fill/essay
    Content string        `json:"content"`
    Options []Option      `json:"options,omitempty"`
    Answer  interface{}   `json:"answer"` // 正确答案
    Score   int           `json:"score"`
    Analysis string       `json:"analysis"` // 答案解析
}

// Option 选项
type Option struct {
    ID      string `json:"id"`
    Content string `json:"content"`
}

// ExamSession 考试会话
type ExamSession struct {
    ID          string              `json:"id"`
    ExamID      string              `json:"exam_id"`
    UserID      string              `json:"user_id"`
    Answers     map[string]interface{} `json:"answers"` // questionID -> answer
    Score       int                 `json:"score"`
    Status      string              `json:"status"` // in_progress/submitted/graded
    StartTime   time.Time           `json:"start_time"`
    SubmitTime  time.Time           `json:"submit_time,omitempty"`
    RemainingTime int               `json:"remaining_time"` // 剩余时间（秒）
}

// ExamService 考试服务
type ExamService struct {
    repo ExamRepository
}

// StartExam 开始考试
func (s *ExamService) StartExam(ctx context.Context, examID, userID string) (*ExamSession, error) {
    // 获取考试信息
    exam, err := s.repo.GetExam(ctx, examID)
    if err != nil {
        return nil, err
    }

    // 检查考试时间
    now := time.Now()
    if now.Before(exam.StartTime) {
        return nil, ErrExamNotStarted
    }
    if now.After(exam.EndTime) {
        return nil, ErrExamEnded
    }

    // 检查是否已参加过
    if exists, _ := s.repo.HasSession(ctx, examID, userID); exists {
        return nil, ErrAlreadyTaken
    }

    // 创建考试会话
    session := &ExamSession{
        ID:            generateID(),
        ExamID:        examID,
        UserID:        userID,
        Answers:       make(map[string]interface{}),
        Status:        "in_progress",
        StartTime:     now,
        RemainingTime: exam.Duration * 60,
    }

    if err := s.repo.CreateSession(ctx, session); err != nil {
        return nil, err
    }

    // 启动倒计时
    go s.countdown(session.ID, exam.Duration*60)

    return session, nil
}

// SubmitAnswer 提交答案
func (s *ExamService) SubmitAnswer(
    ctx context.Context,
    sessionID, questionID string,
    answer interface{},
) error {
    session, err := s.repo.GetSession(ctx, sessionID)
    if err != nil {
        return err
    }

    if session.Status != "in_progress" {
        return ErrSessionClosed
    }

    // 保存答案
    session.Answers[questionID] = answer

    return s.repo.UpdateSession(ctx, session)
}

// SubmitExam 提交考试
func (s *ExamService) SubmitExam(ctx context.Context, sessionID string) (*ExamResult, error) {
    session, err := s.repo.GetSession(ctx, sessionID)
    if err != nil {
        return nil, err
    }

    exam, err := s.repo.GetExam(ctx, session.ExamID)
    if err != nil {
        return nil, err
    }

    // 自动批改
    result := s.gradeExam(exam, session)

    // 更新会话状态
    session.Status = "graded"
    session.Score = result.TotalScore
    session.SubmitTime = time.Now()
    s.repo.UpdateSession(ctx, session)

    return result, nil
}

// ExamResult 考试结果
type ExamResult struct {
    SessionID     string                `json:"session_id"`
    UserID        string                `json:"user_id"`
    ExamID        string                `json:"exam_id"`
    TotalScore    int                   `json:"total_score"`
    PassScore     int                   `json:"pass_score"`
    IsPassed      bool                  `json:"is_passed"`
    QuestionResults []QuestionResult    `json:"question_results"`
    SubmitTime    time.Time             `json:"submit_time"`
}

// QuestionResult 题目结果
type QuestionResult struct {
    QuestionID string      `json:"question_id"`
    UserAnswer interface{} `json:"user_answer"`
    IsCorrect  bool        `json:"is_correct"`
    Score      int         `json:"score"`
    Analysis   string      `json:"analysis"`
}

func (s *ExamService) gradeExam(exam *Exam, session *ExamSession) *ExamResult {
    result := &ExamResult{
        SessionID:  session.ID,
        UserID:     session.UserID,
        ExamID:     exam.ID,
        PassScore:  exam.PassScore,
        SubmitTime: time.Now(),
    }

    for _, question := range exam.Questions {
        userAnswer, exists := session.Answers[question.ID]
        if !exists {
            result.QuestionResults = append(result.QuestionResults, QuestionResult{
                QuestionID: question.ID,
                IsCorrect:  false,
                Score:      0,
                Analysis:   question.Analysis,
            })
            continue
        }

        // 判断答案是否正确
        isCorrect := s.checkAnswer(question, userAnswer)
        score := 0
        if isCorrect {
            score = question.Score
            result.TotalScore += score
        }

        result.QuestionResults = append(result.QuestionResults, QuestionResult{
            QuestionID: question.ID,
            UserAnswer: userAnswer,
            IsCorrect:  isCorrect,
            Score:      score,
            Analysis:   question.Analysis,
        })
    }

    result.IsPassed = result.TotalScore >= result.PassScore

    return result
}

func (s *ExamService) checkAnswer(question Question, userAnswer interface{}) bool {
    switch question.Type {
    case "single", "judge":
        return question.Answer == userAnswer
    case "multiple":
        // 比较数组
        correctAnswers, _ := question.Answer.([]string)
        userAnswers, _ := userAnswer.([]string)
        if len(correctAnswers) != len(userAnswers) {
            return false
        }
        for i := range correctAnswers {
            if correctAnswers[i] != userAnswers[i] {
                return false
            }
        }
        return true
    case "fill", "essay":
        // 填空题和问答题需要人工批改
        return false
    }
    return false
}

func (s *ExamService) countdown(sessionID string, duration int) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    remaining := duration

    for range ticker.C {
        remaining--
        if remaining <= 0 {
            // 时间到，自动提交
            s.SubmitExam(context.Background(), sessionID)
            break
        }
    }
}
```

---

## 8. 视频直播与点播

### 8.1 直播系统

```go
package streaming

import (
    "context"
    "sync"
    "time"
)

// LiveStream 直播流
type LiveStream struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    TeacherID   string    `json:"teacher_id"`
    CourseID    string    `json:"course_id"`
    StreamKey   string    `json:"stream_key"` // 推流密钥
    PlayURL     string    `json:"play_url"` // 播放地址
    Status      string    `json:"status"` // preparing/live/ended
    Viewers     int       `json:"viewers"` // 当前观看人数
    MaxViewers  int       `json:"max_viewers"` // 最高观看人数
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time,omitempty"`
}

// LiveStreamManager 直播管理器
type LiveStreamManager struct {
    streams map[string]*LiveStream
    viewers map[string]map[string]*Viewer // streamID -> userID -> Viewer
    mu      sync.RWMutex
}

// Viewer 观众
type Viewer struct {
    UserID    string
    JoinTime  time.Time
    HeartbeatTime time.Time
}

func NewLiveStreamManager() *LiveStreamManager {
    return &LiveStreamManager{
        streams: make(map[string]*LiveStream),
        viewers: make(map[string]map[string]*Viewer),
    }
}

// CreateStream 创建直播
func (m *LiveStreamManager) CreateStream(ctx context.Context, req *CreateStreamRequest) (*LiveStream, error) {
    stream := &LiveStream{
        ID:        generateID(),
        Title:     req.Title,
        TeacherID: req.TeacherID,
        CourseID:  req.CourseID,
        StreamKey: generateStreamKey(),
        Status:    "preparing",
    }

    // 生成播放地址
    stream.PlayURL = fmt.Sprintf("rtmp://live.example.com/live/%s", stream.StreamKey)

    m.mu.Lock()
    m.streams[stream.ID] = stream
    m.viewers[stream.ID] = make(map[string]*Viewer)
    m.mu.Unlock()

    return stream, nil
}

// StartStream 开始直播
func (m *LiveStreamManager) StartStream(ctx context.Context, streamID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    stream, exists := m.streams[streamID]
    if !exists {
        return ErrStreamNotFound
    }

    stream.Status = "live"
    stream.StartTime = time.Now()

    // 通知所有观众
    go m.notifyViewers(streamID, "stream_started")

    return nil
}

// JoinStream 加入直播
func (m *LiveStreamManager) JoinStream(ctx context.Context, streamID, userID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    stream, exists := m.streams[streamID]
    if !exists {
        return ErrStreamNotFound
    }

    if stream.Status != "live" {
        return ErrStreamNotLive
    }

    // 添加观众
    viewer := &Viewer{
        UserID:        userID,
        JoinTime:      time.Now(),
        HeartbeatTime: time.Now(),
    }

    m.viewers[streamID][userID] = viewer
    stream.Viewers++

    if stream.Viewers > stream.MaxViewers {
        stream.MaxViewers = stream.Viewers
    }

    return nil
}

// LeaveStream 离开直播
func (m *LiveStreamManager) LeaveStream(ctx context.Context, streamID, userID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    stream, exists := m.streams[streamID]
    if !exists {
        return ErrStreamNotFound
    }

    delete(m.viewers[streamID], userID)
    stream.Viewers--

    return nil
}

// Heartbeat 心跳
func (m *LiveStreamManager) Heartbeat(ctx context.Context, streamID, userID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    viewers, exists := m.viewers[streamID]
    if !exists {
        return ErrStreamNotFound
    }

    viewer, exists := viewers[userID]
    if !exists {
        return ErrViewerNotFound
    }

    viewer.HeartbeatTime = time.Now()

    return nil
}

// CheckViewers 检查观众活跃状态
func (m *LiveStreamManager) CheckViewers(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            m.removeInactiveViewers()
        }
    }
}

func (m *LiveStreamManager) removeInactiveViewers() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    timeout := 60 * time.Second

    for streamID, viewers := range m.viewers {
        for userID, viewer := range viewers {
            if now.Sub(viewer.HeartbeatTime) > timeout {
                delete(viewers, userID)
                if stream, exists := m.streams[streamID]; exists {
                    stream.Viewers--
                }
            }
        }
    }
}

func (m *LiveStreamManager) notifyViewers(streamID, event string) {
    // 通过WebSocket或消息队列通知观众
}

// EndStream 结束直播
func (m *LiveStreamManager) EndStream(ctx context.Context, streamID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    stream, exists := m.streams[streamID]
    if !exists {
        return ErrStreamNotFound
    }

    stream.Status = "ended"
    stream.EndTime = time.Now()

    // 通知所有观众
    go m.notifyViewers(streamID, "stream_ended")

    return nil
}

// CreateStreamRequest 创建直播请求
type CreateStreamRequest struct {
    Title     string
    TeacherID string
    CourseID  string
}

func generateStreamKey() string {
    // 生成唯一的推流密钥
    return generateID()
}
```

---

## 9. 完整项目：在线学习平台

### 9.1 项目结构

```text
online-learning-platform/
├── cmd/
│   ├── api/              # API服务
│   ├── worker/           # 后台任务
│   └── admin/            # 管理后台
├── internal/
│   ├── user/             # 用户模块
│   ├── course/           # 课程模块
│   ├── learning/         # 学习模块
│   ├── exam/             # 考试模块
│   ├── interactive/      # 互动模块
│   ├── analytics/        # 分析模块
│   └── common/           # 公共模块
├── pkg/
│   ├── auth/             # 认证
│   ├── storage/          # 存储
│   ├── cache/            # 缓存
│   └── queue/            # 消息队列
├── web/                  # 前端代码
├── migrations/           # 数据库迁移
├── docs/                 # 文档
└── deploy/               # 部署配置
```

### 9.2 核心API实现

```go
package api

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// Server API服务器
type Server struct {
    router  *gin.Engine
    user    *UserService
    course  *CourseService
    learning *LearningService
    exam    *ExamService
}

func NewServer(
    user *UserService,
    course *CourseService,
    learning *LearningService,
    exam *ExamService,
) *Server {
    s := &Server{
        router:   gin.Default(),
        user:     user,
        course:   course,
        learning: learning,
        exam:     exam,
    }

    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    // 公开路由
    public := s.router.Group("/api/v1")
    {
        public.POST("/auth/register", s.handleRegister)
        public.POST("/auth/login", s.handleLogin)
        public.GET("/courses", s.handleListCourses)
        public.GET("/courses/:id", s.handleGetCourse)
    }

    // 需要认证的路由
    auth := s.router.Group("/api/v1")
    auth.Use(AuthMiddleware())
    {
        // 用户相关
        auth.GET("/user/profile", s.handleGetProfile)
        auth.PUT("/user/profile", s.handleUpdateProfile)

        // 学习相关
        auth.POST("/courses/:id/enroll", s.handleEnrollCourse)
        auth.POST("/learning/start", s.handleStartLearning)
        auth.POST("/learning/progress", s.handleUpdateProgress)
        auth.GET("/learning/stats", s.handleGetLearningStats)

        // 考试相关
        auth.POST("/exams/:id/start", s.handleStartExam)
        auth.POST("/exams/:id/answer", s.handleSubmitAnswer)
        auth.POST("/exams/:id/submit", s.handleSubmitExam)
        auth.GET("/exams/:id/result", s.handleGetExamResult)
    }

    // 教师路由
    teacher := s.router.Group("/api/v1/teacher")
    teacher.Use(AuthMiddleware(), TeacherMiddleware())
    {
        teacher.POST("/courses", s.handleCreateCourse)
        teacher.PUT("/courses/:id", s.handleUpdateCourse)
        teacher.POST("/courses/:id/publish", s.handlePublishCourse)
        teacher.GET("/courses/:id/students", s.handleGetStudents)
    }
}

// handleRegister 处理注册
func (s *Server) handleRegister(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := s.user.Register(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

// handleEnrollCourse 处理课程注册
func (s *Server) handleEnrollCourse(c *gin.Context) {
    courseID := c.Param("id")
    userID := c.GetString("user_id")

    if err := s.course.Enroll(c.Request.Context(), userID, courseID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "enrolled successfully"})
}

// handleStartLearning 处理开始学习
func (s *Server) handleStartLearning(c *gin.Context) {
    var req StartLearningRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetString("user_id")

    if err := s.learning.StartLearning(c.Request.Context(), userID, req.CourseID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "started"})
}

// Run 启动服务器
func (s *Server) Run(addr string) error {
    return s.router.Run(addr)
}
```

---

## 💡 总结

### 核心要点

1. **技术架构**
   - 微服务架构，模块化设计
   - 高并发处理，支持大量用户
   - 实时互动，WebSocket通信
   - 数据分析，个性化推荐

2. **关键技术**
   - 视频处理：ffmpeg转码、多清晰度
   - 直播系统：RTMP推流、HLS/RTMP播放
   - 实时互动：WebSocket、问答、白板
   - 学习分析：行为追踪、数据挖掘

3. **工程实践**
   - 使用CDN加速内容分发
   - 实现多级缓存策略
   - 异步处理耗时任务
   - 完善的监控告警体系

### 进阶方向

1. **AI赋能教育**
   - 智能推荐课程
   - 自适应学习路径
   - 自动批改作业
   - 学习效果预测

2. **增强互动体验**
   - VR/AR虚拟教室
   - 实时语音识别
   - 智能答疑助手
   - 多人协作学习

3. **数据驱动决策**
   - 学习行为分析
   - 知识图谱构建
   - 学习效果评估
   - 精准营销推荐

---

## 🔗 相关资源

**开源项目**:

- [Moodle](https://github.com/moodle/moodle) - 开源LMS
- [Open edX](https://github.com/openedx) - edX平台
- [Canvas LMS](https://github.com/instructure/canvas-lms)

**参考资料**:

- [在线教育平台架构设计](https://example.com)
- [视频直播技术实践](https://example.com)
- [学习数据分析方法](https://example.com)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月24日  
**文档状态**: ✅ 完成  
**适用版本**: Go 1.21+
