package types

import "gorm.io/gorm"

type ParseStatus string

const (
	ParseStatusNotParsed ParseStatus = "notParsed" // 未解析
	ParseStatusParsing   ParseStatus = "parsing"   // 解析中
	ParseStatusParsed    ParseStatus = "parsed"    // 已解析
	ParseStatusFailed    ParseStatus = "failed"    // 解析失败
)

// Document 知识库文档
type Document struct {
	BaseEntity
	KnowledgeBaseId  int64       `json:"knowledge_base_id" gorm:"not null;"`                   // 知识库id
	Name             string      `json:"name" gorm:"not null;type:varchar(255)"`               // 文档名称
	OriginalFileName string      `json:"original_file_name" gorm:"not null;type:varchar(255)"` // 原始文件名
	FileType         string      `json:"file_type" gorm:"not null;type:varchar(255)"`          // 文件类型
	Size             int64       `json:"size" gorm:"not null;"`
	ParseStatus      ParseStatus `json:"parse_status" gorm:"not null;type:varchar(255)"` // 解析状态
	ParseLog         string      `json:"parse_log" gorm:"not null;type:text;"`           // 解析日志
	StoragePath      string      `json:"storage_path" gorm:"not null;type:varchar(255)"` // 存储路径
	NumChunks        int64       `json:"num_chunks" gorm:"not null;type:bigint;"`        // 文档切片数量
	CreateBy         int64       `json:"create_by" gorm:"not null;"`
	WorkspaceId      int64       `json:"workspace_id" gorm:"not null;type:bigint;"`
}

func (*Document) TableName() string {
	return "documents"
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if err := d.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if d.WorkspaceId == 0 {
		workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(int64)
		if ok {
			d.WorkspaceId = workspaceId
		}
	}
	if d.CreateBy == 0 {
		createBy, ok := tx.Statement.Context.Value(UserIdContextKey).(int64)
		if ok {
			d.CreateBy = createBy
		}
	}
	d.ParseStatus = ParseStatusNotParsed
	return nil
}

type DocumentQuery struct {
	BaseQuery
	Name            string      `form:"name"`
	KnowledgeBaseId int64       `form:"knowledge_base_id" binding:"required"`
	ParseStatus     ParseStatus `form:"parse_status"`
}

type RenameDocumentReq struct {
	Id   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
