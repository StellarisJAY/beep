package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type UserWorkspaceRepo struct {
	db *gorm.DB
}

func NewUserWorkspaceRepo(db *gorm.DB) interfaces.UserWorkspaceRepo {
	return &UserWorkspaceRepo{
		db: db,
	}
}

func (u *UserWorkspaceRepo) Create(ctx context.Context, uw *types.UserWorkspace) error {
	return u.db.WithContext(ctx).Create(uw).Error
}

func (u *UserWorkspaceRepo) FindByUser(ctx context.Context, userId int64) ([]*types.UserWorkspace, error) {
	var uw []*types.UserWorkspace
	err := u.db.WithContext(ctx).Model(&types.UserWorkspace{}).
		Where("user_id = ?", userId).
		Find(&uw).Error
	if err != nil {
		return nil, err
	}
	return uw, nil
}

func (u *UserWorkspaceRepo) Find(ctx context.Context, workspaceId, userId int64) (*types.UserWorkspace, error) {
	var uw *types.UserWorkspace
	err := u.db.WithContext(ctx).Model(&types.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ?", workspaceId, userId).
		First(&uw).Error
	if err != nil {
		return nil, err
	}
	return uw, nil
}

func (u *UserWorkspaceRepo) FindByWorkspace(ctx context.Context, workspaceId int64) (*types.UserWorkspace, error) {
	var uw *types.UserWorkspace
	err := u.db.WithContext(ctx).Model(&types.UserWorkspace{}).
		Where("workspace_id = ?", workspaceId).
		First(&uw).Error
	if err != nil {
		return nil, err
	}
	return uw, nil
}

func (u *UserWorkspaceRepo) Delete(ctx context.Context, id int64) error {
	return u.db.WithContext(ctx).Delete(&types.UserWorkspace{}, id).Error
}

func (u *UserWorkspaceRepo) ListMember(ctx context.Context, workspaceId int64) ([]*types.WorkspaceMember, error) {
	var members []*types.WorkspaceMember
	err := u.db.WithContext(ctx).Table("user_workspaces uw").
		Joins("JOIN users u ON u.id = uw.user_id").
		Where("uw.workspace_id = ?", workspaceId).
		Select("uw.role, u.*").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (u *UserWorkspaceRepo) FindUserDefaultWorkspace(ctx context.Context, userId int64) (*types.Workspace, error) {
	var workspace *types.Workspace
	err := u.db.WithContext(ctx).Model(&types.Workspace{}).
		Joins("JOIN user_workspaces uw ON uw.workspace_id = workspaces.id").
		Where("uw.user_id = ?", userId).Where("uw.role = ?", types.WorkspaceRoleOwner).
		Select("workspaces.*").
		First(&workspace).Error
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (u *UserWorkspaceRepo) FindUserJoinedWorkspace(ctx context.Context, userId int64) ([]*types.Workspace, error) {
	var workspaces []*types.Workspace
	err := u.db.WithContext(ctx).Model(&types.Workspace{}).
		Joins("JOIN user_workspaces uw ON uw.workspace_id = workspaces.id").
		Where("uw.user_id = ?", userId).
		Select("workspaces.*").
		Find(&workspaces).Error
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (u *UserWorkspaceRepo) Update(ctx context.Context, uw *types.UserWorkspace) error {
	return u.db.WithContext(ctx).Model(&types.UserWorkspace{}).Where("id = ?", uw.ID).Updates(uw).Error
}
