package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type ModelService struct {
	modelRepo        interfaces.ModelRepo
	modelFactoryRepo interfaces.ModelFactoryRepo
	encryptService   interfaces.EncryptService
}

var templates []types.ModelFactoryTemplate

func init() {
	initModelFactoryDefault()
}

func NewModelService(modelRepo interfaces.ModelRepo,
	modelFactoryRepo interfaces.ModelFactoryRepo,
	encryptService interfaces.EncryptService) interfaces.ModelService {
	return &ModelService{modelRepo: modelRepo, modelFactoryRepo: modelFactoryRepo, encryptService: encryptService}
}

func (m *ModelService) CreateFactory(ctx context.Context, req types.CreateModelFactoryReq) error {
	index := slices.IndexFunc(templates, func(template types.ModelFactoryTemplate) bool {
		return template.SdkType == string(req.Type)
	})
	if index == -1 {
		return errors.NewBadRequestError("不支持的模型供应商", nil)
	}
	template := templates[index]
	if req.BaseUrl == "" {
		if template.DefaultConfig.BaseUrl != "" {
			req.BaseUrl = template.DefaultConfig.BaseUrl
		} else {
			return errors.NewBadRequestError("模型供应商未配置默认base_url", nil)
		}
	}
	// 配置加密
	encryptedKey, err := m.encryptService.Encrypt(ctx, req.ApiKey)
	if err != nil {
		return errors.NewInternalServerError("新增模型供应商配置失败", err)
	}
	factory := &types.ModelFactory{
		Name:    req.Name,
		Type:    req.Type,
		BaseUrl: req.BaseUrl,
		APIKey:  encryptedKey,
	}
	if err := m.modelFactoryRepo.Create(ctx, factory); err != nil {
		return errors.NewInternalServerError("新增模型供应商失败", err)
	}
	// 添加供应商默认模型
	models := make([]*types.Model, 0, len(template.Models))
	for _, m := range template.Models {
		model := &types.Model{
			Name:         m.Name,
			Type:         m.Type,
			Tags:         m.Tags,
			MaxTokens:    m.MaxTokens,
			FunctionCall: m.FunctionCall,
			FactoryId:    factory.ID,
			Status:       true,
		}
		models = append(models, model)
	}
	if err := m.modelRepo.CreateMany(ctx, models); err != nil {
		return errors.NewInternalServerError("新增模型供应商模型失败", err)
	}
	return nil
}

func (m *ModelService) UpdateFactory(ctx context.Context, req types.UpdateModelFactoryReq) error {
	// 配置加密
	encryptedKey, err := m.encryptService.Encrypt(ctx, req.ApiKey)
	if err != nil {
		return errors.NewInternalServerError("新增模型供应商配置失败", err)
	}
	if err := m.modelFactoryRepo.Update(ctx, &types.ModelFactory{
		BaseEntity: types.BaseEntity{ID: req.Id},
		Name:       req.Name,
		BaseUrl:    req.BaseUrl,
		APIKey:     encryptedKey,
	}); err != nil {
		return errors.NewInternalServerError("更新模型供应商失败", err)
	}
	return nil
}

func (m *ModelService) ListFactory(ctx context.Context) ([]*types.ModelFactory, error) {
	factories, err := m.modelFactoryRepo.List(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError("查询模型供应商失败", err)
	}
	for _, factory := range factories {
		query := types.ListModelQuery{FactoryId: factory.ID}
		models, err := m.modelRepo.List(ctx, query)
		if err != nil {
			return nil, errors.NewInternalServerError("查询模型供应商模型失败", err)
		}
		factory.Models = models
	}
	return factories, nil
}

func (m *ModelService) ListModels(ctx context.Context, query types.ListModelQuery) ([]*types.Model, error) {
	return m.modelRepo.List(ctx, query)
}

func initModelFactoryDefault() {
	templates = []types.ModelFactoryTemplate{}
	entries, err := os.ReadDir("config/models")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// 匹配文件名 model_factory_*.json
		if !strings.HasPrefix(entry.Name(), "model_factory_") || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		// 读取文件内容
		data, err := os.ReadFile("config/models/" + entry.Name())
		if err != nil {
			panic(err)
		}
		var template types.ModelFactoryTemplate
		if err := json.Unmarshal(data, &template); err != nil {
			slog.Error("invalid model factory template", "file", entry.Name(), "err", err)
			continue
		}
		templates = append(templates, template)
	}
}

func (m *ModelService) GetModelDetail(ctx context.Context, id string) (*types.ModelDetail, error) {
	detail, err := m.modelRepo.GetDetail(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerError("获取模型详情失败", err)
	}
	detail.ApiKeyDecrypted, err = m.encryptService.Decrypt(ctx, detail.ApiKey)
	if err != nil {
		return nil, errors.NewInternalServerError("获取模型详情失败，APIKey解密失败", err)
	}
	return detail, nil
}
