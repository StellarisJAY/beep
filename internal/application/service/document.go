package service

import (
	"beep/internal/errors"
	"beep/internal/models"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"io"
	"mime/multipart"

	"github.com/panjf2000/ants/v2"
)

type DocumentService struct {
	documentRepo      interfaces.DocumentRepo
	knowledgeBaseRepo interfaces.KnowledgeBaseRepo
	fileStore         interfaces.FileStore
	parseService      interfaces.ParseService
	modelService      interfaces.ModelService
	pool              *ants.Pool
}

func NewDocumentService(documentRepo interfaces.DocumentRepo,
	knowledgeBaseRepo interfaces.KnowledgeBaseRepo,
	fileStore interfaces.FileStore,
	parseService interfaces.ParseService,
	modelService interfaces.ModelService,
	pool *ants.Pool) interfaces.DocumentService {
	return &DocumentService{
		documentRepo:      documentRepo,
		knowledgeBaseRepo: knowledgeBaseRepo,
		fileStore:         fileStore,
		parseService:      parseService,
		modelService:      modelService,
		pool:              pool,
	}
}

// CreateFromFile 从文件创建文档
func (d *DocumentService) CreateFromFile(ctx context.Context, knowledgeBaseId int64, file *multipart.FileHeader) error {
	reader, err := file.Open()
	if err != nil {
		return errors.NewInternalServerError("新增文档失败", err)
	}
	defer reader.Close()
	kb, _ := d.knowledgeBaseRepo.FindById(ctx, knowledgeBaseId)
	if kb == nil {
		return errors.NewInternalServerError("新增文档失败，知识库不存在", nil)
	}
	document := &types.Document{
		KnowledgeBaseId:  knowledgeBaseId,
		Name:             file.Filename,
		OriginalFileName: file.Filename,
		FileType:         file.Header.Get("Content-Type"),
		Size:             file.Size,
	}
	key, err := d.fileStore.Upload(ctx, kb.StorageBucketName(), reader)
	if err != nil {
		return errors.NewInternalServerError("新增文档失败，文件上传失败", err)
	}
	document.StoragePath = key
	if err := d.documentRepo.Create(ctx, document); err != nil {
		return errors.NewInternalServerError("新增文档失败，数据库写入失败", err)
	}
	return nil
}

// Delete 删除文档
func (d *DocumentService) Delete(ctx context.Context, id int64) error {
	document, err := d.documentRepo.Get(ctx, id)
	if err != nil {
		return errors.NewInternalServerError("删除文档失败，数据库查询失败", err)
	}
	if document == nil {
		return errors.NewInternalServerError("删除文档失败，文档不存在", nil)
	}
	kb, _ := d.knowledgeBaseRepo.FindById(ctx, document.KnowledgeBaseId)
	if kb == nil {
		return errors.NewInternalServerError("删除文档失败，知识库不存在", nil)
	}
	if err := d.fileStore.Delete(ctx, kb.StorageBucketName(), document.StoragePath); err != nil {
		return errors.NewInternalServerError("删除文档失败，文件删除失败", err)
	}
	if err := d.documentRepo.Delete(ctx, id); err != nil {
		return errors.NewInternalServerError("删除文档失败，数据库删除失败", err)
	}
	return nil
}

// List 列表文档
func (d *DocumentService) List(ctx context.Context, query types.DocumentQuery) ([]*types.Document, int, error) {
	return d.documentRepo.List(ctx, query)
}

// Rename 重命名文档
func (d *DocumentService) Rename(ctx context.Context, req types.RenameDocumentReq) error {
	document, _ := d.documentRepo.Get(ctx, req.Id)
	if document == nil {
		return errors.NewInternalServerError("重命名文档失败，文档不存在", nil)
	}
	document.Name = req.Name
	if err := d.documentRepo.Update(ctx, document); err != nil {
		return errors.NewInternalServerError("重命名文档失败，数据库更新失败", err)
	}
	return nil
}

// Download 下载文档,返回文件下载URL
func (d *DocumentService) Download(ctx context.Context, id int64) (string, error) {
	document, _ := d.documentRepo.Get(ctx, id)
	if document == nil {
		return "", errors.NewInternalServerError("下载文档失败，文档不存在", nil)
	}
	kb, _ := d.knowledgeBaseRepo.FindById(ctx, document.KnowledgeBaseId)
	if kb == nil {
		return "", errors.NewInternalServerError("下载文档失败，知识库不存在", nil)
	}
	url, err := d.fileStore.TempDownloadURL(ctx, kb.StorageBucketName(), document.StoragePath)
	if err != nil {
		return "", errors.NewInternalServerError("下载文档失败，文件下载URL生成失败", err)
	}
	return url, nil
}

// Parse 解析文档，切片、嵌入、存储、生成摘要、构建知识图谱等
func (d *DocumentService) Parse(ctx context.Context, id int64) error {
	document, _ := d.documentRepo.Get(ctx, id)
	if document == nil {
		return errors.NewInternalServerError("解析文档失败，文档不存在", nil)
	}
	if document.ParseStatus == types.ParseStatusParsing {
		return errors.NewInternalServerError("文档正在解析中", nil)
	}
	kb, _ := d.knowledgeBaseRepo.FindById(ctx, document.KnowledgeBaseId)
	if kb == nil {
		return errors.NewInternalServerError("解析文档失败，知识库不存在", nil)
	}
	// 获取嵌入模型
	embeddingModel, err := d.modelService.GetModelDetail(ctx, kb.EmbeddingModel)
	if err != nil {
		return errors.NewInternalServerError("解析文档失败，获取嵌入模型失败", err)
	}
	embedder, err := models.CreateEmbedder(*embeddingModel)
	if err != nil {
		return errors.NewInternalServerError("解析文档失败，创建嵌入模型失败", err)
	}
	// 获取聊天模型
	chatModelDetail, err := d.modelService.GetModelDetail(ctx, kb.ChatModel)
	if err != nil {
		return errors.NewInternalServerError("解析文档失败，获取聊天模型失败", err)
	}
	chatModel, err := models.NewChatModel(*chatModelDetail, types.ChatModelOption{
		ResponseFormat: "json",
	})
	if err != nil {
		return errors.NewInternalServerError("解析文档失败，创建聊天模型失败", err)
	}

	// 读取文档内容
	reader, err := d.fileStore.Download(ctx, kb.StorageBucketName(), document.StoragePath)
	if err != nil {
		return errors.NewInternalServerError("解析文档失败，文件下载失败", err)
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return errors.NewInternalServerError("解析文档失败，文件读取失败", err)
	}

	parseInfo := types.ParseInfo{
		Content:              string(content),
		DocId:                document.ID,
		KbId:                 kb.ID,
		ChunkOptions:         kb.ChunkOptions,
		Embedder:             embedder,
		ChatModel:            chatModel,
		EnableKnowledgeGraph: false,
	}

	if err := d.parseService.Parse(ctx, parseInfo); err != nil {
		return errors.NewInternalServerError("解析文档失败，解析服务失败", err)
	}
	return nil

}
