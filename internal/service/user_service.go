package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ome-app-back/api/middleware"
	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// UserService 处理用户相关业务逻辑
type UserService struct {
	userDAO       *dao.AppUserDAO
	userWeightDAO *dao.UserWeightDAO
	userGoalDAO   *dao.UserGoalDAO
}

// NewUserService 创建用户服务实例
func NewUserService(userDAO *dao.AppUserDAO, userWeightDAO *dao.UserWeightDAO, userGoalDAO *dao.UserGoalDAO) *UserService {
	return &UserService{
		userDAO:       userDAO,
		userWeightDAO: userWeightDAO,
		userGoalDAO:   userGoalDAO,
	}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// Register 用户注册
func (s *UserService) Register(req RegisterRequest) (*RegisterResponse, error) {
	// 检查手机号是否已存在
	if req.Phone != "" {
		existUser, _ := s.userDAO.GetByPhone(req.Phone)
		if existUser != nil {
			return nil, errors.New("手机号已注册")
		}
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		existUser, _ := s.userDAO.GetByEmail(req.Email)
		if existUser != nil {
			return nil, errors.New("邮箱已注册")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.AppUser{
		UserName:     req.UserName,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		// BirthDate字段不设置，将使用数据库默认值NULL
	}

	if err := s.userDAO.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 生成JWT Token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &RegisterResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Account  string `json:"account"` // 可以是手机号或邮箱
	Password string `json:"password"`
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	UserID            int64  `json:"user_id"`
	UserName          string `json:"user_name"`
	Token             string `json:"token"`
	IsProfileComplete bool   `json:"is_profile_complete"`
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	var user *model.AppUser
	var err error

	// 通过手机号或邮箱查找用户
	if len(req.Account) > 0 {
		if strings.Contains(req.Account, "@") {
			user, err = s.userDAO.GetByEmail(req.Account)
		} else {
			user, err = s.userDAO.GetByPhone(req.Account)
		}
	}

	if user == nil || err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 检查用户档案是否完善
	isProfileComplete := user.HeightCM > 0 && !user.BirthDate.IsZero() && user.Sex != ""

	// 生成JWT Token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 打印用户信息用于调试
	fmt.Printf("用户登录: ID=%d, UserName=%s, Email=%s\n", user.ID, user.UserName, user.Email)

	return &LoginResponse{
		UserID:            user.ID,
		UserName:          user.UserName,
		Token:             token,
		IsProfileComplete: isProfileComplete,
	}, nil
}

// UpdateProfileRequest 更新用户档案请求
type UpdateProfileRequest struct {
	UserID    int64   `json:"user_id"`
	HeightCM  float64 `json:"height_cm"`
	BirthDate string  `json:"birth_date"` // 格式 YYYY-MM-DD
	Sex       string  `json:"sex"`        // male/female/other
	WeightKG  float64 `json:"weight_kg"`
}

// UpdateProfile 更新用户基本档案
func (s *UserService) UpdateProfile(req UpdateProfileRequest) error {
	// 获取用户
	user, err := s.userDAO.GetByID(req.UserID)
	if err != nil {
		return errors.New("获取用户信息失败")
	}

	// 更新身高
	if req.HeightCM > 0 {
		if req.HeightCM < 50 || req.HeightCM > 300 {
			return errors.New("身高超出合理范围(50-300cm)")
		}
		user.HeightCM = req.HeightCM
	}

	// 更新出生日期
	if req.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return errors.New("无效的日期格式")
		}
		user.BirthDate = birthDate
	}

	// 更新性别
	if req.Sex != "" {
		user.Sex = req.Sex
	}

	// 保存用户信息
	if err := s.userDAO.Update(user); err != nil {
		return errors.New("更新基本信息失败")
	}

	// 记录体重
	if req.WeightKG > 0 {
		if err := s.userWeightDAO.Create(req.UserID, req.WeightKG); err != nil {
			return errors.New("记录体重失败")
		}
	}

	return nil
}

// UpdateGoalRequest 更新健康目标请求
type UpdateGoalRequest struct {
	UserID           int64    `json:"user_id"`
	GoalType         string   `json:"goal_type"` // lose_fat/keep_fit/gain_muscle
	TargetWeightKG   float64  `json:"target_weight_kg"`
	WeeklyChangeKG   float64  `json:"weekly_change_kg"`
	TargetDate       string   `json:"target_date"` // 格式 YYYY-MM-DD
	DietType         string   `json:"diet_type"`   // normal/vegetarian/low_carb等
	TastePreferences []string `json:"taste_preferences"`
	FoodIntolerances []string `json:"food_intolerances"`
}

// UpdateGoal 更新用户健康目标
func (s *UserService) UpdateGoal(req UpdateGoalRequest) error {
	// 解析日期
	targetDate, err := time.Parse("2006-01-02", req.TargetDate)
	if err != nil {
		return errors.New("无效的日期格式")
	}

	// 创建或更新用户目标
	err = s.userGoalDAO.CreateOrUpdate(
		req.UserID,
		req.GoalType,
		req.TargetWeightKG,
		req.WeeklyChangeKG,
		targetDate,
		req.DietType,
		req.TastePreferences,
		req.FoodIntolerances,
	)

	if err != nil {
		return errors.New("更新健康目标失败")
	}

	return nil
}
