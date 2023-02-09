// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//

package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	"Simple-Douyin/cmd/user/dal/db"

	user "Simple-Douyin/kitex_gen/user"
	"Simple-Douyin/pkg/errno"
)

type CreateUserService struct {
	ctx context.Context
}

// NewCreateUserService new CreateUserService
func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{ctx: ctx}
}

// CreateUser create user info.
func (s *CreateUserService) CreateUser(req *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {
	users, err := db.QueryUser(s.ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if len(users) != 0 {
		return nil, errno.UserAlreadyExistErr
	}

	h := md5.New()
	if _, err = io.WriteString(h, req.Password); err != nil {
		return nil, err
	}
	passWord := fmt.Sprintf("%x", h.Sum(nil))

	users, err = db.CreateUser(s.ctx, []*db.User{{
		Username: req.Username,
		Password: passWord,
	}})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, err
	}
	resp := &user.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		UserId:     users[0].Id,
	}
	return resp, err
}
