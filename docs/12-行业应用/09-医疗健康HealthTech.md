# Go在医疗健康（HealthTech）中的应用

> **简介**: 系统介绍Go语言在电子病历、远程医疗、健康管理、医疗影像等医疗健康领域的架构设计、技术实践与工程落地

---

## 📚 目录

- [Go在医疗健康（HealthTech）中的应用](#go在医疗健康healthtech中的应用)
  - [📚 目录](#-目录)
  - [1. 医疗健康技术概览](#1-医疗健康技术概览)
    - [1.1 行业特点](#11-行业特点)
    - [1.2 Go的优势](#12-go的优势)
  - [2. 电子病历系统（EHR/EMR）](#2-电子病历系统ehremr)
    - [2.1 患者信息管理](#21-患者信息管理)
    - [2.2 病历查询与统计](#22-病历查询与统计)
  - [3. 远程医疗平台](#3-远程医疗平台)
    - [3.1 在线问诊](#31-在线问诊)
    - [3.2 视频会诊](#32-视频会诊)
  - [4. 健康数据管理](#4-健康数据管理)
    - [4.1 健康档案](#41-健康档案)
  - [5. 医疗影像处理](#5-医疗影像处理)
    - [5.1 DICOM影像处理](#51-dicom影像处理)
  - [6. 处方管理系统](#6-处方管理系统)
    - [6.1 电子处方](#61-电子处方)

---

## 1. 医疗健康技术概览

### 1.1 行业特点

**核心需求**:

- 数据安全性（患者隐私、HIPAA合规）
- 高可靠性（7×24小时服务）
- 实时性（急诊、远程会诊）
- 数据互联互通（医院系统集成）
- 审计追溯（完整操作日志）

**技术挑战**:

- 海量医疗数据存储
- 医疗影像大文件处理
- 多系统数据集成
- 隐私数据加密
- 高并发号源抢占

### 1.2 Go的优势

```go
// Go在HealthTech中的优势
优势特性:
✅ 高性能 - 处理海量医疗数据
✅ 并发处理 - 支持大量患者同时访问
✅ 内存安全 - 避免内存泄漏和崩溃
✅ 静态类型 - 减少医疗系统bug
✅ 部署简单 - 便于医院信息化升级
```

---

## 2. 电子病历系统（EHR/EMR）

### 2.1 患者信息管理

```go
package ehr

import (
    "context"
    "time"
)

// Patient 患者信息
type Patient struct {
    ID            string       `json:"id"`
    Name          string       `json:"name"`
    Gender        string       `json:"gender"`
    BirthDate     time.Time    `json:"birth_date"`
    IDNumber      string       `json:"id_number"` // 身份证号（加密）
    Phone         string       `json:"phone"` // 手机号（加密）
    Address       string       `json:"address"`
    BloodType     string       `json:"blood_type"`
    Allergies     []string     `json:"allergies"` // 过敏史
    ChronicDiseases []string   `json:"chronic_diseases"` // 慢性病
    EmergencyContact EmergencyContact `json:"emergency_contact"`
    InsuranceInfo InsuranceInfo `json:"insurance_info"` // 医保信息
    CreatedAt     time.Time    `json:"created_at"`
    UpdatedAt     time.Time    `json:"updated_at"`
}

// EmergencyContact 紧急联系人
type EmergencyContact struct {
    Name     string `json:"name"`
    Relation string `json:"relation"`
    Phone    string `json:"phone"`
}

// InsuranceInfo 医保信息
type InsuranceInfo struct {
    Type       string `json:"type"` // 医保类型
    CardNumber string `json:"card_number"` // 医保卡号
    ExpiryDate time.Time `json:"expiry_date"`
}

// MedicalRecord 病历
type MedicalRecord struct {
    ID            string          `json:"id"`
    PatientID     string          `json:"patient_id"`
    VisitID       string          `json:"visit_id"` // 就诊号
    DoctorID      string          `json:"doctor_id"`
    Department    string          `json:"department"` // 科室
    ChiefComplaint string         `json:"chief_complaint"` // 主诉
    PresentIllness string         `json:"present_illness"` // 现病史
    PastHistory   string          `json:"past_history"` // 既往史
    Examination   *PhysicalExam   `json:"examination"` // 体格检查
    Diagnosis     []Diagnosis     `json:"diagnosis"` // 诊断
    Treatment     *TreatmentPlan  `json:"treatment"` // 治疗方案
    Prescriptions []Prescription  `json:"prescriptions"` // 处方
    LabTests      []LabTest       `json:"lab_tests"` // 检验
    Imaging       []ImagingStudy  `json:"imaging"` // 影像
    Status        string          `json:"status"` // draft/completed/signed
    CreatedAt     time.Time       `json:"created_at"`
    SignedAt      time.Time       `json:"signed_at,omitempty"`
}

// PhysicalExam 体格检查
type PhysicalExam struct {
    Temperature    float64 `json:"temperature"` // 体温
    Pulse          int     `json:"pulse"` // 脉搏
    Respiration    int     `json:"respiration"` // 呼吸
    BloodPressure  string  `json:"blood_pressure"` // 血压
    Height         float64 `json:"height"` // 身高(cm)
    Weight         float64 `json:"weight"` // 体重(kg)
    GeneralCondition string `json:"general_condition"`
    Notes          string  `json:"notes"`
}

// Diagnosis 诊断
type Diagnosis struct {
    Code        string `json:"code"` // ICD-10编码
    Name        string `json:"name"`
    Type        string `json:"type"` // primary/secondary
    Description string `json:"description"`
}

// TreatmentPlan 治疗方案
type TreatmentPlan struct {
    Medications []string `json:"medications"`
    Procedures  []string `json:"procedures"`
    Advice      string   `json:"advice"` // 医嘱
    FollowUp    time.Time `json:"follow_up,omitempty"` // 复诊时间
}

// EHRService 电子病历服务
type EHRService struct {
    repo       EHRRepository
    encryption EncryptionService
    audit      AuditService
}

// CreateMedicalRecord 创建病历
func (s *EHRService) CreateMedicalRecord(
    ctx context.Context,
    req *CreateRecordRequest,
) (*MedicalRecord, error) {
    // 验证医生权限
    if err := s.validateDoctor(ctx, req.DoctorID, req.Department); err != nil {
        return nil, err
    }

    // 获取患者信息
    patient, err := s.repo.GetPatient(ctx, req.PatientID)
    if err != nil {
        return nil, err
    }

    record := &MedicalRecord{
        ID:         generateID(),
        PatientID:  req.PatientID,
        VisitID:    req.VisitID,
        DoctorID:   req.DoctorID,
        Department: req.Department,
        Status:     "draft",
        CreatedAt:  time.Now(),
    }

    // 保存病历
    if err := s.repo.CreateRecord(ctx, record); err != nil {
        return nil, err
    }

    // 记录审计日志
    s.audit.Log(ctx, &AuditLog{
        Action:     "create_medical_record",
        ResourceID: record.ID,
        UserID:     req.DoctorID,
        Timestamp:  time.Now(),
    })

    return record, nil
}

// UpdateMedicalRecord 更新病历
func (s *EHRService) UpdateMedicalRecord(
    ctx context.Context,
    recordID string,
    updates *RecordUpdates,
) error {
    // 获取原病历
    record, err := s.repo.GetRecord(ctx, recordID)
    if err != nil {
        return err
    }

    // 检查病历状态（已签署的不能修改）
    if record.Status == "signed" {
        return ErrRecordSigned
    }

    // 应用更新
    s.applyUpdates(record, updates)
    record.UpdatedAt = time.Now()

    // 保存更新
    if err := s.repo.UpdateRecord(ctx, record); err != nil {
        return err
    }

    // 审计日志
    s.audit.Log(ctx, &AuditLog{
        Action:     "update_medical_record",
        ResourceID: recordID,
        Changes:    updates,
        Timestamp:  time.Now(),
    })

    return nil
}

// SignMedicalRecord 签署病历
func (s *EHRService) SignMedicalRecord(ctx context.Context, recordID, doctorID string) error {
    record, err := s.repo.GetRecord(ctx, recordID)
    if err != nil {
        return err
    }

    // 验证医生权限
    if record.DoctorID != doctorID {
        return ErrUnauthorized
    }

    // 验证病历完整性
    if err := s.validateRecord(record); err != nil {
        return err
    }

    // 更新状态
    record.Status = "signed"
    record.SignedAt = time.Now()

    if err := s.repo.UpdateRecord(ctx, record); err != nil {
        return err
    }

    // 审计日志
    s.audit.Log(ctx, &AuditLog{
        Action:     "sign_medical_record",
        ResourceID: recordID,
        UserID:     doctorID,
        Timestamp:  time.Now(),
    })

    return nil
}

func (s *EHRService) validateRecord(record *MedicalRecord) error {
    if record.ChiefComplaint == "" {
        return ErrMissingChiefComplaint
    }
    if len(record.Diagnosis) == 0 {
        return ErrMissingDiagnosis
    }
    return nil
}

func (s *EHRService) applyUpdates(record *MedicalRecord, updates *RecordUpdates) {
    if updates.ChiefComplaint != "" {
        record.ChiefComplaint = updates.ChiefComplaint
    }
    if updates.PresentIllness != "" {
        record.PresentIllness = updates.PresentIllness
    }
    if len(updates.Diagnosis) > 0 {
        record.Diagnosis = updates.Diagnosis
    }
    if updates.Treatment != nil {
        record.Treatment = updates.Treatment
    }
}
```

### 2.2 病历查询与统计

```go
package ehr

import (
    "context"
    "time"
)

// RecordQuery 病历查询
type RecordQuery struct {
    PatientID  string
    DoctorID   string
    Department string
    StartDate  time.Time
    EndDate    time.Time
    Status     string
    Diagnosis  string // ICD-10编码
    Limit      int
    Offset     int
}

// SearchRecords 搜索病历
func (s *EHRService) SearchRecords(ctx context.Context, query *RecordQuery) ([]*MedicalRecord, int, error) {
    // 验证查询权限
    userID := ctx.Value("user_id").(string)
    if err := s.validateQueryPermission(ctx, userID, query); err != nil {
        return nil, 0, err
    }

    // 执行查询
    records, total, err := s.repo.SearchRecords(ctx, query)
    if err != nil {
        return nil, 0, err
    }

    // 审计日志
    s.audit.Log(ctx, &AuditLog{
        Action:    "search_medical_records",
        UserID:    userID,
        Query:     query,
        Timestamp: time.Now(),
    })

    return records, total, nil
}

// GetPatientHistory 获取患者完整病史
func (s *EHRService) GetPatientHistory(ctx context.Context, patientID string) (*PatientHistory, error) {
    // 获取所有病历
    records, _, err := s.repo.SearchRecords(ctx, &RecordQuery{
        PatientID: patientID,
    })
    if err != nil {
        return nil, err
    }

    // 构建病史时间线
    history := &PatientHistory{
        PatientID: patientID,
        Timeline:  make([]HistoryEvent, 0),
    }

    for _, record := range records {
        event := HistoryEvent{
            Date:       record.CreatedAt,
            Type:       "visit",
            Department: record.Department,
            Doctor:     record.DoctorID,
            Diagnosis:  record.Diagnosis,
            RecordID:   record.ID,
        }
        history.Timeline = append(history.Timeline, event)
    }

    // 按时间排序
    sort.Slice(history.Timeline, func(i, j int) bool {
        return history.Timeline[i].Date.After(history.Timeline[j].Date)
    })

    return history, nil
}

// PatientHistory 患者病史
type PatientHistory struct {
    PatientID string         `json:"patient_id"`
    Timeline  []HistoryEvent `json:"timeline"`
}

// HistoryEvent 历史事件
type HistoryEvent struct {
    Date       time.Time   `json:"date"`
    Type       string      `json:"type"` // visit/surgery/hospitalization
    Department string      `json:"department"`
    Doctor     string      `json:"doctor"`
    Diagnosis  []Diagnosis `json:"diagnosis"`
    RecordID   string      `json:"record_id"`
}

// GetDepartmentStats 获取科室统计
func (s *EHRService) GetDepartmentStats(
    ctx context.Context,
    department string,
    startDate, endDate time.Time,
) (*DepartmentStats, error) {
    stats := &DepartmentStats{
        Department: department,
        Period:     Period{Start: startDate, End: endDate},
    }

    // 统计就诊人数
    stats.TotalVisits, _ = s.repo.CountVisits(ctx, department, startDate, endDate)

    // 统计常见疾病
    stats.CommonDiseases, _ = s.repo.GetCommonDiseases(ctx, department, startDate, endDate, 10)

    // 统计医生工作量
    stats.DoctorWorkload, _ = s.repo.GetDoctorWorkload(ctx, department, startDate, endDate)

    return stats, nil
}

// DepartmentStats 科室统计
type DepartmentStats struct {
    Department     string              `json:"department"`
    Period         Period              `json:"period"`
    TotalVisits    int                 `json:"total_visits"`
    CommonDiseases []DiseaseStats      `json:"common_diseases"`
    DoctorWorkload []DoctorWorkloadStat `json:"doctor_workload"`
}

// DiseaseStats 疾病统计
type DiseaseStats struct {
    Code  string `json:"code"` // ICD-10
    Name  string `json:"name"`
    Count int    `json:"count"`
}

// DoctorWorkloadStat 医生工作量统计
type DoctorWorkloadStat struct {
    DoctorID    string `json:"doctor_id"`
    DoctorName  string `json:"doctor_name"`
    TotalVisits int    `json:"total_visits"`
    AvgDuration int    `json:"avg_duration"` // 平均诊疗时长（分钟）
}

// Period 时间段
type Period struct {
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`
}
```

---

## 3. 远程医疗平台

### 3.1 在线问诊

```go
package telemedicine

import (
    "context"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

// Consultation 在线问诊
type Consultation struct {
    ID            string           `json:"id"`
    PatientID     string           `json:"patient_id"`
    DoctorID      string           `json:"doctor_id"`
    Type          string           `json:"type"` // text/video/audio
    Status        string           `json:"status"` // waiting/in_progress/completed/cancelled
    Symptoms      string           `json:"symptoms"` // 症状描述
    Images        []string         `json:"images"` // 症状图片
    Messages      []ConsultMessage `json:"messages"`
    Diagnosis     string           `json:"diagnosis,omitempty"`
    Prescription  *Prescription    `json:"prescription,omitempty"`
    CreatedAt     time.Time        `json:"created_at"`
    StartTime     time.Time        `json:"start_time,omitempty"`
    EndTime       time.Time        `json:"end_time,omitempty"`
}

// ConsultMessage 问诊消息
type ConsultMessage struct {
    ID        string    `json:"id"`
    SenderID  string    `json:"sender_id"`
    SenderType string   `json:"sender_type"` // patient/doctor
    Type      string    `json:"type"` // text/image/voice
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}

// ConsultationService 问诊服务
type ConsultationService struct {
    repo      ConsultationRepository
    matching  DoctorMatchingService
    websocket *WebSocketManager
}

// CreateConsultation 创建问诊
func (s *ConsultationService) CreateConsultation(
    ctx context.Context,
    req *CreateConsultationRequest,
) (*Consultation, error) {
    // 创建问诊记录
    consult := &Consultation{
        ID:        generateID(),
        PatientID: req.PatientID,
        Type:      req.Type,
        Symptoms:  req.Symptoms,
        Images:    req.Images,
        Status:    "waiting",
        Messages:  make([]ConsultMessage, 0),
        CreatedAt: time.Now(),
    }

    // 如果指定了医生，直接分配
    if req.DoctorID != "" {
        consult.DoctorID = req.DoctorID
    } else {
        // 否则智能匹配医生
        doctor, err := s.matching.MatchDoctor(ctx, req)
        if err != nil {
            return nil, err
        }
        consult.DoctorID = doctor.ID
    }

    // 保存问诊
    if err := s.repo.Create(ctx, consult); err != nil {
        return nil, err
    }

    // 通知医生
    s.notifyDoctor(consult.DoctorID, consult)

    return consult, nil
}

// StartConsultation 开始问诊
func (s *ConsultationService) StartConsultation(ctx context.Context, consultID, doctorID string) error {
    consult, err := s.repo.Get(ctx, consultID)
    if err != nil {
        return err
    }

    if consult.DoctorID != doctorID {
        return ErrUnauthorized
    }

    consult.Status = "in_progress"
    consult.StartTime = time.Now()

    if err := s.repo.Update(ctx, consult); err != nil {
        return err
    }

    // 通知患者
    s.notifyPatient(consult.PatientID, "consultation_started", consult)

    return nil
}

// SendMessage 发送消息
func (s *ConsultationService) SendMessage(
    ctx context.Context,
    consultID, senderID, senderType, content string,
) error {
    consult, err := s.repo.Get(ctx, consultID)
    if err != nil {
        return err
    }

    if consult.Status != "in_progress" {
        return ErrConsultationNotActive
    }

    // 创建消息
    msg := ConsultMessage{
        ID:         generateID(),
        SenderID:   senderID,
        SenderType: senderType,
        Type:       "text",
        Content:    content,
        Timestamp:  time.Now(),
    }

    consult.Messages = append(consult.Messages, msg)

    // 保存
    if err := s.repo.Update(ctx, consult); err != nil {
        return err
    }

    // 实时推送消息
    recipientID := consult.PatientID
    if senderType == "patient" {
        recipientID = consult.DoctorID
    }
    s.websocket.SendMessage(recipientID, msg)

    return nil
}

// CompleteConsultation 完成问诊
func (s *ConsultationService) CompleteConsultation(
    ctx context.Context,
    consultID string,
    diagnosis string,
    prescription *Prescription,
) error {
    consult, err := s.repo.Get(ctx, consultID)
    if err != nil {
        return err
    }

    consult.Status = "completed"
    consult.Diagnosis = diagnosis
    consult.Prescription = prescription
    consult.EndTime = time.Now()

    if err := s.repo.Update(ctx, consult); err != nil {
        return err
    }

    // 通知患者
    s.notifyPatient(consult.PatientID, "consultation_completed", consult)

    return nil
}

func (s *ConsultationService) notifyDoctor(doctorID string, consult *Consultation) {
    // 通过WebSocket、推送通知等方式通知医生
}

func (s *ConsultationService) notifyPatient(patientID, event string, consult *Consultation) {
    // 通知患者
}
```

### 3.2 视频会诊

```go
package telemedicine

import (
    "context"
    "time"
)

// VideoConsultation 视频会诊
type VideoConsultation struct {
    ID          string    `json:"id"`
    ConsultID   string    `json:"consult_id"`
    RoomID      string    `json:"room_id"` // 会议室ID
    Token       string    `json:"token"` // 加入凭证
    Status      string    `json:"status"` // waiting/active/ended
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time,omitempty"`
    Duration    int       `json:"duration"` // 时长（秒）
    RecordingURL string   `json:"recording_url,omitempty"` // 录像地址
}

// VideoService 视频服务
type VideoService struct {
    webrtc WebRTCService
    repo   VideoRepository
}

// CreateVideoRoom 创建视频会议室
func (s *VideoService) CreateVideoRoom(ctx context.Context, consultID string) (*VideoConsultation, error) {
    // 创建WebRTC房间
    room, err := s.webrtc.CreateRoom(ctx)
    if err != nil {
        return nil, err
    }

    video := &VideoConsultation{
        ID:        generateID(),
        ConsultID: consultID,
        RoomID:    room.ID,
        Status:    "waiting",
        StartTime: time.Now(),
    }

    // 生成加入凭证
    video.Token, _ = s.webrtc.GenerateToken(room.ID)

    if err := s.repo.Create(ctx, video); err != nil {
        return nil, err
    }

    return video, nil
}

// JoinVideoRoom 加入视频会议
func (s *VideoService) JoinVideoRoom(
    ctx context.Context,
    roomID, userID, userType string,
) (*JoinInfo, error) {
    // 验证权限
    video, err := s.repo.GetByRoomID(ctx, roomID)
    if err != nil {
        return nil, err
    }

    // 生成加入信息
    token, err := s.webrtc.GenerateToken(roomID)
    if err != nil {
        return nil, err
    }

    joinInfo := &JoinInfo{
        RoomID:    roomID,
        Token:     token,
        ICEServers: s.webrtc.GetICEServers(),
    }

    // 如果是首次加入，更新状态
    if video.Status == "waiting" {
        video.Status = "active"
        s.repo.Update(ctx, video)
    }

    return joinInfo, nil
}

// EndVideoCall 结束视频通话
func (s *VideoService) EndVideoCall(ctx context.Context, roomID string) error {
    video, err := s.repo.GetByRoomID(ctx, roomID)
    if err != nil {
        return err
    }

    video.Status = "ended"
    video.EndTime = time.Now()
    video.Duration = int(video.EndTime.Sub(video.StartTime).Seconds())

    // 停止录制
    recordingURL, _ := s.webrtc.StopRecording(roomID)
    video.RecordingURL = recordingURL

    if err := s.repo.Update(ctx, video); err != nil {
        return err
    }

    // 关闭WebRTC房间
    s.webrtc.CloseRoom(roomID)

    return nil
}

// JoinInfo 加入信息
type JoinInfo struct {
    RoomID     string      `json:"room_id"`
    Token      string      `json:"token"`
    ICEServers []ICEServer `json:"ice_servers"`
}

// ICEServer ICE服务器
type ICEServer struct {
    URLs       []string `json:"urls"`
    Username   string   `json:"username,omitempty"`
    Credential string   `json:"credential,omitempty"`
}
```

---

## 4. 健康数据管理

### 4.1 健康档案

```go
package health

import (
    "context"
    "time"
)

// HealthProfile 健康档案
type HealthProfile struct {
    UserID          string            `json:"user_id"`
    BasicInfo       BasicHealthInfo   `json:"basic_info"`
    VitalSigns      []VitalSignRecord `json:"vital_signs"` // 生命体征记录
    HealthMetrics   HealthMetrics     `json:"health_metrics"` // 健康指标
    MedicalHistory  []MedicalEvent    `json:"medical_history"` // 病史
    Medications     []Medication      `json:"medications"` // 用药记录
    Vaccinations    []Vaccination     `json:"vaccinations"` // 疫苗接种
    Allergies       []string          `json:"allergies"` // 过敏史
    FamilyHistory   []FamilyDisease   `json:"family_history"` // 家族病史
    LifestyleData   LifestyleData     `json:"lifestyle_data"` // 生活方式
    UpdatedAt       time.Time         `json:"updated_at"`
}

// BasicHealthInfo 基本健康信息
type BasicHealthInfo struct {
    Height      float64 `json:"height"` // cm
    Weight      float64 `json:"weight"` // kg
    BMI         float64 `json:"bmi"`
    BloodType   string  `json:"blood_type"`
    RhFactor    string  `json:"rh_factor"` // 阳性/阴性
}

// VitalSignRecord 生命体征记录
type VitalSignRecord struct {
    Timestamp     time.Time `json:"timestamp"`
    Temperature   float64   `json:"temperature"` // 体温
    BloodPressure BP        `json:"blood_pressure"` // 血压
    HeartRate     int       `json:"heart_rate"` // 心率
    Respiration   int       `json:"respiration"` // 呼吸频率
    OxygenSat     int       `json:"oxygen_saturation"` // 血氧饱和度
    Source        string    `json:"source"` // manual/device/hospital
}

// BP 血压
type BP struct {
    Systolic  int `json:"systolic"` // 收缩压
    Diastolic int `json:"diastolic"` // 舒张压
}

// HealthMetrics 健康指标
type HealthMetrics struct {
    BloodGlucose    float64 `json:"blood_glucose"` // 血糖
    Cholesterol     Lipids  `json:"cholesterol"` // 血脂
    LiverFunction   LiverFunctionTest `json:"liver_function"` // 肝功能
    KidneyFunction  RenalFunctionTest `json:"kidney_function"` // 肾功能
}

// Lipids 血脂
type Lipids struct {
    TotalCholesterol float64 `json:"total_cholesterol"` // 总胆固醇
    HDL             float64 `json:"hdl"` // 高密度脂蛋白
    LDL             float64 `json:"ldl"` // 低密度脂蛋白
    Triglycerides   float64 `json:"triglycerides"` // 甘油三酯
}

// LifestyleData 生活方式数据
type LifestyleData struct {
    SmokingStatus   string        `json:"smoking_status"` // never/former/current
    AlcoholUse      string        `json:"alcohol_use"` // none/occasional/regular
    Exercise        ExerciseData  `json:"exercise"`
    Sleep           SleepData     `json:"sleep"`
    Diet            DietData      `json:"diet"`
}

// ExerciseData 运动数据
type ExerciseData struct {
    FrequencyPerWeek int     `json:"frequency_per_week"` // 每周运动次数
    DurationMinutes  int     `json:"duration_minutes"` // 每次运动时长
    IntensityLevel   string  `json:"intensity_level"` // low/moderate/high
    DailySteps       int     `json:"daily_steps"` // 每日步数
}

// SleepData 睡眠数据
type SleepData struct {
    AverageDuration float64 `json:"average_duration"` // 平均睡眠时长（小时）
    Quality         string  `json:"quality"` // poor/fair/good
    Bedtime         string  `json:"bedtime"` // 就寝时间
    WakeTime        string  `json:"wake_time"` // 起床时间
}

// HealthProfileService 健康档案服务
type HealthProfileService struct {
    repo   HealthRepository
    ai     AIHealthAssistant
}

// UpdateVitalSigns 更新生命体征
func (s *HealthProfileService) UpdateVitalSigns(
    ctx context.Context,
    userID string,
    vitals *VitalSignRecord,
) error {
    profile, err := s.repo.GetProfile(ctx, userID)
    if err != nil {
        return err
    }

    vitals.Timestamp = time.Now()
    profile.VitalSigns = append(profile.VitalSigns, *vitals)

    // 分析健康趋势
    analysis := s.ai.AnalyzeVitalSigns(profile.VitalSigns)
    if analysis.HasAbnormality {
        // 发送健康警告
        s.sendHealthAlert(userID, analysis)
    }

    profile.UpdatedAt = time.Now()
    return s.repo.UpdateProfile(ctx, profile)
}

// GetHealthSummary 获取健康摘要
func (s *HealthProfileService) GetHealthSummary(
    ctx context.Context,
    userID string,
) (*HealthSummary, error) {
    profile, err := s.repo.GetProfile(ctx, userID)
    if err != nil {
        return nil, err
    }

    summary := &HealthSummary{
        UserID:    userID,
        BMI:       profile.BasicInfo.BMI,
        BMIStatus: s.getBMIStatus(profile.BasicInfo.BMI),
    }

    // 最近的生命体征
    if len(profile.VitalSigns) > 0 {
        latest := profile.VitalSigns[len(profile.VitalSigns)-1]
        summary.LatestVitalSigns = &latest
    }

    // 健康评分
    summary.HealthScore = s.calculateHealthScore(profile)

    // 健康建议
    summary.Recommendations = s.ai.GenerateRecommendations(profile)

    return summary, nil
}

// HealthSummary 健康摘要
type HealthSummary struct {
    UserID           string           `json:"user_id"`
    BMI              float64          `json:"bmi"`
    BMIStatus        string           `json:"bmi_status"`
    LatestVitalSigns *VitalSignRecord `json:"latest_vital_signs"`
    HealthScore      int              `json:"health_score"` // 0-100
    Recommendations  []string         `json:"recommendations"`
}

func (s *HealthProfileService) getBMIStatus(bmi float64) string {
    if bmi < 18.5 {
        return "偏瘦"
    } else if bmi < 24 {
        return "正常"
    } else if bmi < 28 {
        return "超重"
    } else {
        return "肥胖"
    }
}

func (s *HealthProfileService) calculateHealthScore(profile *HealthProfile) int {
    score := 100

    // 根据BMI扣分
    if profile.BasicInfo.BMI < 18.5 || profile.BasicInfo.BMI > 28 {
        score -= 10
    }

    // 根据血压扣分
    if len(profile.VitalSigns) > 0 {
        latest := profile.VitalSigns[len(profile.VitalSigns)-1]
        if latest.BloodPressure.Systolic > 140 || latest.BloodPressure.Diastolic > 90 {
            score -= 15
        }
    }

    // 根据生活方式扣分
    if profile.LifestyleData.SmokingStatus == "current" {
        score -= 20
    }
    if profile.LifestyleData.Exercise.FrequencyPerWeek < 3 {
        score -= 10
    }

    if score < 0 {
        score = 0
    }

    return score
}

func (s *HealthProfileService) sendHealthAlert(userID string, analysis *HealthAnalysis) {
    // 发送健康警告通知
}
```

---

## 5. 医疗影像处理

### 5.1 DICOM影像处理

```go
package imaging

import (
    "context"
    "fmt"
)

// ImagingStudy 影像检查
type ImagingStudy struct {
    ID            string          `json:"id"`
    PatientID     string          `json:"patient_id"`
    StudyDate     time.Time       `json:"study_date"`
    Modality      string          `json:"modality"` // CT/MRI/X-Ray/Ultrasound
    BodyPart      string          `json:"body_part"`
    Description   string          `json:"description"`
    Status        string          `json:"status"` // scheduled/in_progress/completed
    Series        []ImageSeries   `json:"series"`
    Report        *ImagingReport  `json:"report,omitempty"`
    Annotations   []Annotation    `json:"annotations"`
}

// ImageSeries 影像序列
type ImageSeries struct {
    SeriesID      string    `json:"series_id"`
    SeriesNumber  int       `json:"series_number"`
    Description   string    `json:"description"`
    ImageCount    int       `json:"image_count"`
    Images        []DICOMImage `json:"images"`
}

// DICOMImage DICOM影像
type DICOMImage struct {
    ImageID       string `json:"image_id"`
    InstanceNumber int   `json:"instance_number"`
    URL           string `json:"url"` // 存储URL
    ThumbnailURL  string `json:"thumbnail_url"`
    Width         int    `json:"width"`
    Height        int    `json:"height"`
    WindowCenter  int    `json:"window_center"` // 窗位
    WindowWidth   int    `json:"window_width"` // 窗宽
}

// ImagingReport 影像报告
type ImagingReport struct {
    ID          string    `json:"id"`
    StudyID     string    `json:"study_id"`
    Radiologist string    `json:"radiologist"` // 阅片医生
    Findings    string    `json:"findings"` // 影像所见
    Impression  string    `json:"impression"` // 影像诊断
    CreatedAt   time.Time `json:"created_at"`
    SignedAt    time.Time `json:"signed_at,omitempty"`
}

// Annotation 标注
type Annotation struct {
    ID        string     `json:"id"`
    ImageID   string     `json:"image_id"`
    Type      string     `json:"type"` // rectangle/circle/arrow/text
    Geometry  Geometry   `json:"geometry"`
    Label     string     `json:"label"`
    CreatedBy string     `json:"created_by"`
}

// Geometry 几何信息
type Geometry struct {
    Type       string    `json:"type"`
    Points     []Point   `json:"points"`
}

// Point 坐标点
type Point struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}

// ImagingService 影像服务
type ImagingService struct {
    repo     ImagingRepository
    storage  StorageClient
    ai       AIImagingService
}

// UploadDICOM 上传DICOM文件
func (s *ImagingService) UploadDICOM(
    ctx context.Context,
    studyID string,
    dicomFile []byte,
) (*DICOMImage, error) {
    // 解析DICOM文件
    dicom, err := parseDICOM(dicomFile)
    if err != nil {
        return nil, err
    }

    // 生成缩略图
    thumbnail, err := generateThumbnail(dicom)
    if err != nil {
        return nil, err
    }

    // 上传原始文件
    url, err := s.storage.Upload(ctx, dicomFile, "application/dicom")
    if err != nil {
        return nil, err
    }

    // 上传缩略图
    thumbnailURL, err := s.storage.Upload(ctx, thumbnail, "image/jpeg")
    if err != nil {
        return nil, err
    }

    image := &DICOMImage{
        ImageID:      generateID(),
        URL:          url,
        ThumbnailURL: thumbnailURL,
        Width:        dicom.Width,
        Height:       dicom.Height,
        WindowCenter: dicom.WindowCenter,
        WindowWidth:  dicom.WindowWidth,
    }

    return image, nil
}

// CreateAnnotation 创建标注
func (s *ImagingService) CreateAnnotation(
    ctx context.Context,
    annotation *Annotation,
) error {
    annotation.ID = generateID()
    return s.repo.CreateAnnotation(ctx, annotation)
}

// AIAnalyzeImage AI影像分析
func (s *ImagingService) AIAnalyzeImage(
    ctx context.Context,
    imageID string,
) (*AIAnalysisResult, error) {
    // 获取影像
    image, err := s.repo.GetImage(ctx, imageID)
    if err != nil {
        return nil, err
    }

    // AI分析
    result, err := s.ai.AnalyzeImage(ctx, image.URL)
    if err != nil {
        return nil, err
    }

    // 自动生成标注
    for _, finding := range result.Findings {
        annotation := &Annotation{
            ImageID:  imageID,
            Type:     "rectangle",
            Geometry: finding.Location,
            Label:    finding.Label,
            CreatedBy: "AI",
        }
        s.repo.CreateAnnotation(ctx, annotation)
    }

    return result, nil
}

// AIAnalysisResult AI分析结果
type AIAnalysisResult struct {
    Confidence float64        `json:"confidence"`
    Findings   []AIFinding    `json:"findings"`
    Suggestions []string      `json:"suggestions"`
}

// AIFinding AI发现
type AIFinding struct {
    Label      string   `json:"label"`
    Confidence float64  `json:"confidence"`
    Location   Geometry `json:"location"`
}

func parseDICOM(data []byte) (*DICOMInfo, error) {
    // 使用DICOM库解析
    // 这里简化处理
    return &DICOMInfo{}, nil
}

func generateThumbnail(dicom *DICOMInfo) ([]byte, error) {
    // 生成缩略图
    return nil, nil
}

// DICOMInfo DICOM信息
type DICOMInfo struct {
    Width        int
    Height       int
    WindowCenter int
    WindowWidth  int
}
```

---

## 6. 处方管理系统

### 6.1 电子处方

```go
package prescription

import (
    "context"
    "time"
)

// Prescription 处方
type Prescription struct {
    ID          string            `json:"id"`
    PatientID   string            `json:"patient_id"`
    DoctorID    string            `json:"doctor_id"`
    ConsultID   string            `json:"consult_id,omitempty"`
    RecordID    string            `json:"record_id,omitempty"`
    Type        string            `json:"type"` // western/chinese
    Items       []PrescriptionItem `json:"items"`
    Diagnosis   string            `json:"diagnosis"`
    Notes       string            `json:"notes"` // 医嘱
    Status      string            `json:"status"` // draft/issued/dispensed/cancelled
    IssuedAt    time.Time         `json:"issued_at"`
    ValidUntil  time.Time         `json:"valid_until"`
    Signature   string            `json:"signature"` // 医生签名
}

// PrescriptionItem 处方项目
type PrescriptionItem struct {
    ID            string  `json:"id"`
    DrugName      string  `json:"drug_name"` // 药品名称
    DrugCode      string  `json:"drug_code"` // 药品编码
    Specification string  `json:"specification"` // 规格
    Dosage        string  `json:"dosage"` // 用量
    Frequency     string  `json:"frequency"` // 用法频次
    Duration      string  `json:"duration"` // 疗程
    Quantity      float64 `json:"quantity"` // 数量
    Unit          string  `json:"unit"` // 单位
    Notes         string  `json:"notes"` // 备注
}

// PrescriptionService 处方服务
type PrescriptionService struct {
    repo       PrescriptionRepository
    drugCheck  DrugInteractionChecker
    signature  SignatureService
}

// CreatePrescription 创建处方
func (s *PrescriptionService) CreatePrescription(
    ctx context.Context,
    req *CreatePrescriptionRequest,
) (*Prescription, error) {
    // 验证医生资质
    if err := s.validateDoctor(ctx, req.DoctorID); err != nil {
        return nil, err
    }

    // 检查药物相互作用
    if err := s.drugCheck.CheckInteractions(req.Items); err != nil {
        return nil, err
    }

    // 检查患者过敏史
    patient, _ := s.getPatient(ctx, req.PatientID)
    if err := s.checkAllergies(req.Items, patient.Allergies); err != nil {
        return nil, err
    }

    prescription := &Prescription{
        ID:         generateID(),
        PatientID:  req.PatientID,
        DoctorID:   req.DoctorID,
        ConsultID:  req.ConsultID,
        RecordID:   req.RecordID,
        Type:       req.Type,
        Items:      req.Items,
        Diagnosis:  req.Diagnosis,
        Notes:      req.Notes,
        Status:     "draft",
        IssuedAt:   time.Now(),
        ValidUntil: time.Now().AddDate(0, 0, 3), // 3天有效期
    }

    // 保存处方
    if err := s.repo.Create(ctx, prescription); err != nil {
        return nil, err
    }

    return prescription, nil
}

// IssuePrescription 签发处方
func (s *PrescriptionService) IssuePrescription(
    ctx context.Context,
    prescriptionID, doctorID string,
) error {
    prescription, err := s.repo.Get(ctx, prescriptionID)
    if err != nil {
        return err
    }

    if prescription.DoctorID != doctorID {
        return ErrUnauthorized
    }

    // 生成医生签名
    signature, err := s.signature.Sign(doctorID, prescription)
    if err != nil {
        return err
    }

    prescription.Signature = signature
    prescription.Status = "issued"

    return s.repo.Update(ctx, prescription)
}

// DispensePrescription 配药
func (s *PrescriptionService) DispensePrescription(
    ctx context.Context,
    prescriptionID, pharmacistID string,
) error {
    prescription, err := s.repo.Get(ctx, prescriptionID)
    if err != nil {
        return err
    }

    if prescription.Status != "issued" {
        return ErrInvalidStatus
    }

    // 检查有效期
    if time.Now().After(prescription.ValidUntil) {
        return ErrPrescriptionExpired
    }

    prescription.Status = "dispensed"

    return s.repo.Update(ctx, prescription)
}

func (s *PrescriptionService) checkAllergies(items []PrescriptionItem, allergies []string) error {
    for _, item := range items {
        for _, allergy := range allergies {
            if item.DrugName == allergy {
                return fmt.Errorf("患者对 %s 过敏", allergy)
            }
        }
    }
    return nil
}

// DrugInteractionChecker 药物相互作用检查
type DrugInteractionChecker interface {
    CheckInteractions(items []PrescriptionItem) error
}
```

---

由于响应长度限制，我将在下一个响应中继续完成剩余部分（智能诊断辅助、医疗数据安全、完整项目等）。
