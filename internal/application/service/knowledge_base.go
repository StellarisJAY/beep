package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
)

type KnowledgeBaseService struct {
	knowledgeBaseRepo interfaces.KnowledgeBaseRepo
	userRepo          interfaces.UserRepo
	fileStore         interfaces.FileStore
	vectorStore       interfaces.VectorStore
	documentRepo      interfaces.DocumentRepo
}

func NewKnowledgeBaseService(knowledgeBaseRepo interfaces.KnowledgeBaseRepo,
	userRepo interfaces.UserRepo,
	fileStore interfaces.FileStore,
	vectorStore interfaces.VectorStore,
	documentRepo interfaces.DocumentRepo) interfaces.KnowledgeBaseService {
	return &KnowledgeBaseService{
		knowledgeBaseRepo: knowledgeBaseRepo,
		userRepo:          userRepo,
		fileStore:         fileStore,
		vectorStore:       vectorStore,
		documentRepo:      documentRepo,
	}
}

func (k *KnowledgeBaseService) Create(ctx context.Context, req types.CreateKnowledgeBaseReq) error {
	kb := &types.KnowledgeBase{
		Name:           req.Name,
		Description:    req.Description,
		EmbeddingModel: req.EmbeddingModel,
		ChatModel:      req.ChatModel,
		Public:         *req.Public,
		ChunkOptions:   req.ChunkOptions,
	}

	if kb.ChunkOptions.ChunkSize == 0 {
		kb.ChunkOptions.ChunkSize = 256
	}
	if kb.ChunkOptions.ChunkOverlap == 0 {
		kb.ChunkOptions.ChunkOverlap = kb.ChunkOptions.ChunkSize / 2
	}
	if len(kb.ChunkOptions.Separators) == 0 {
		kb.ChunkOptions.Separators = []string{"\n\n", "\n", ",", "."}
	}

	// 知识库创建
	if err := k.knowledgeBaseRepo.Create(ctx, kb); err != nil {
		return errors.NewInternalServerError("创建知识库失败", err)
	}
	// 向量库创建集合
	if err := k.vectorStore.CreateCollection(ctx, kb.Name, 2048); err != nil {
		return errors.NewInternalServerError("创建知识库失败", err)
	}
	// 存储桶创建
	if err := k.fileStore.CreateBucket(ctx, kb.StorageBucketName()); err != nil {
		return errors.NewInternalServerError("创建知识库失败", err)
	}
	return nil
}

func (k *KnowledgeBaseService) List(ctx context.Context, query types.KnowledgeBaseQuery) ([]*types.KnowledgeBase, int, error) {
	return k.knowledgeBaseRepo.List(ctx, query)
}

func (k *KnowledgeBaseService) Update(ctx context.Context, req types.UpdateKnowledgeBaseReq) error {
	kb, err := k.knowledgeBaseRepo.FindById(ctx, req.Id)
	if err != nil {
		return errors.NewInternalServerError("更新知识库失败", err)
	}
	if kb == nil {
		return errors.NewNotFoundError("知识库不存在", nil)
	}
	if err := k.knowledgeBaseRepo.Update(ctx, kb); err != nil {
		return errors.NewInternalServerError("更新知识库失败", err)
	}
	return nil
}

func (k *KnowledgeBaseService) Delete(ctx context.Context, id string) error {
	kb, err := k.knowledgeBaseRepo.FindById(ctx, id)
	if err != nil {
		return errors.NewInternalServerError("删除知识库失败", err)
	}
	if kb == nil {
		return errors.NewNotFoundError("知识库不存在", nil)
	}
	// 存储桶删除
	if err := k.fileStore.DeleteBucket(ctx, kb.StorageBucketName()); err != nil {
		return errors.NewInternalServerError("删除知识库失败", err)
	}
	// 向量库删除集合
	if err := k.vectorStore.DropCollection(ctx, kb.Name); err != nil {
		return errors.NewInternalServerError("删除知识库失败", err)
	}
	// 文档删除
	if err := k.documentRepo.DeleteByKnowledgeBaseId(ctx, id); err != nil {
		return errors.NewInternalServerError("删除知识库失败", err)
	}
	// 知识库删除
	if err := k.knowledgeBaseRepo.Delete(ctx, id); err != nil {
		return errors.NewInternalServerError("删除知识库失败", err)
	}
	return nil
}
