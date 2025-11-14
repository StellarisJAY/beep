package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"encoding/json"
	errors2 "errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WorkspaceService struct {
	workspaceRepo     interfaces.WorkspaceRepo
	userWorkspaceRepo interfaces.UserWorkspaceRepo
	userRepo          interfaces.UserRepo
	redis             *redis.Client
}

func NewWorkspaceService(workspaceRepo interfaces.WorkspaceRepo,
	userWorkspaceRepo interfaces.UserWorkspaceRepo,
	userRepo interfaces.UserRepo,
	redis *redis.Client) interfaces.WorkspaceService {
	return &WorkspaceService{
		workspaceRepo:     workspaceRepo,
		userWorkspaceRepo: userWorkspaceRepo,
		userRepo:          userRepo,
		redis:             redis,
	}
}

func (w *WorkspaceService) FindById(ctx context.Context, id string) (*types.Workspace, error) {
	return w.workspaceRepo.FindById(ctx, id)
}

func (w *WorkspaceService) ListMembers(ctx context.Context, workspaceId string) ([]*types.WorkspaceMember, error) {
	return w.userWorkspaceRepo.ListMember(ctx, workspaceId)
}

func (w *WorkspaceService) ListByUserId(ctx context.Context, userId string) ([]*types.Workspace, error) {
	return w.userWorkspaceRepo.FindUserJoinedWorkspace(ctx, userId)
}

func (w *WorkspaceService) InviteMember(ctx context.Context, req types.InviteWorkspaceMemberReq) error {
	workspace, _ := w.workspaceRepo.FindById(ctx, req.WorkspaceId)
	if workspace == nil {
		return errors.NewNotFoundError("工作空间不存在", nil)
	}
	if req.Role == "" || req.Role == types.WorkspaceRoleOwner {
		return errors.NewBadRequestError("无效的角色", nil)
	}
	for _, email := range req.Emails {
		user, err := w.userRepo.FindByEmail(ctx, email)
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			// TODO 邀请用户不存在，发送邀请邮件
			continue
		}
		if err != nil {
			return errors.NewInternalServerError("邀请失败", err)
		}

		// 检查用户是否已加入工作空间
		exist, _ := w.userWorkspaceRepo.Find(ctx, req.WorkspaceId, user.ID)
		if exist != nil {
			continue
		}

		// 加入工作空间
		if err := w.userWorkspaceRepo.Create(ctx, &types.UserWorkspace{
			UserId:      user.ID,
			WorkspaceId: req.WorkspaceId,
			Role:        req.Role,
		}); err != nil {
			return errors.NewInternalServerError("邀请失败", err)
		}
	}
	return nil
}

func (w *WorkspaceService) SwitchWorkspace(ctx context.Context, id string) error {
	// 获取登录用户的id
	userId, _ := ctx.Value(types.UserIdContextKey).(string)
	token, _ := ctx.Value(types.AccessTokenContextKey).(string)
	// 查询是否已加入工作空间
	uw, err := w.userWorkspaceRepo.Find(ctx, id, userId)
	if err != nil || uw == nil {
		return errors.NewNotFoundError("无法切换到工作空间", nil)
	}
	// 修改用户登录信息
	loginInfo := types.LoginInfo{
		UserId:      userId,
		WorkspaceId: id,
	}
	data, _ := json.Marshal(loginInfo)
	if err := w.redis.Set(ctx, "access_token:"+token, data, time.Hour).Err(); err != nil {
		return errors.NewInternalServerError("登录信息缓存失败", err)
	}
	return nil
}

func (w *WorkspaceService) SetRole(ctx context.Context, req types.SetWorkspaceRoleReq) error {
	if req.Role == "" || req.Role == types.WorkspaceRoleOwner {
		return errors.NewBadRequestError("无效的角色", nil)
	}
	// 获取登录用户的id，判断是否是admin或owner
	userId, _ := ctx.Value(types.UserIdContextKey).(string)
	uw, _ := w.userWorkspaceRepo.Find(ctx, req.WorkspaceId, userId)
	if uw == nil || (uw.Role != types.WorkspaceRoleAdmin && uw.Role != types.WorkspaceRoleOwner) {
		return errors.NewUnauthorizedError("没有权限", nil)
	}
	// 检查目标用户是否存在
	targetUW, err := w.userWorkspaceRepo.Find(ctx, req.WorkspaceId, req.UserId)
	if err != nil || targetUW == nil {
		return errors.NewNotFoundError("目标用户不存在", nil)
	}
	// 修改角色
	targetUW.Role = req.Role
	if err := w.userWorkspaceRepo.Update(ctx, targetUW); err != nil {
		return errors.NewInternalServerError("设置角色失败", err)
	}
	return nil
}
