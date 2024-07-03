package service

import (
	"context"
	"fmt"
	"muzz/internal/db"
	dbMock "muzz/internal/db/mocks"
	"muzz/internal/model"
	"muzz/internal/security"
	securityMock "muzz/internal/security/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

type ServiceTestSuite struct {
	suite.Suite
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

/*
	These tests are not intended to be comprehensive, they're just
	here to demonstrate that I understand go unit testing principles
*/

func (s *ServiceTestSuite) TestLogin() {

	token := "token"

	tests := []struct {
		name     string
		db       db.DB
		security security.Security
		token    string
		err      error
	}{
		{
			name: "Happy Path",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())

				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(&model.UserPassword{}, nil).
					Once()

				m.EXPECT().BeginTx(mock.Anything).
					Return(nil, nil).
					Once()

				m.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().Commit(mock.Anything).
					Return(nil)

				return m
			}(),

			security: func() security.Security {
				m := securityMock.NewSecurity(s.T())
				m.EXPECT().ComparePassword(mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().CreateToken(mock.Anything).
					Return(token, nil)

				return m
			}(),

			token: token,
		},
		{
			name: "s.db.GetUserPassword() errors",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())

				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("getUserPassword() fails")).
					Once()

				return m
			}(),

			token: "",

			err: fmt.Errorf("getUserPassword() fails"),
		},
		{
			name: "User not found",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())
				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()

				return m
			}(),

			// token: "",
		},
		{
			name: "s.security.ComparePassword() errors",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())

				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(&model.UserPassword{}, nil).
					Once()

				m.EXPECT().BeginTx(mock.Anything).
					Return(nil, nil).
					Maybe()

				m.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Maybe()

				return m
			}(),

			security: func() security.Security {
				m := securityMock.NewSecurity(s.T())
				m.EXPECT().ComparePassword(mock.Anything, mock.Anything).
					Return(fmt.Errorf("s.security.ComparePassword() error"))

				return m
			}(),

			err: fmt.Errorf("s.security.ComparePassword() error"),
		},
		{
			name: "s.security.CreateToken() errors",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())

				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(&model.UserPassword{}, nil).
					Once()

				m.EXPECT().BeginTx(mock.Anything).
					Return(nil, nil).
					Maybe()

				m.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Maybe()

				return m
			}(),

			security: func() security.Security {
				m := securityMock.NewSecurity(s.T())
				m.EXPECT().ComparePassword(mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().CreateToken(mock.Anything).
					Return("", fmt.Errorf("s.security.CreateToken() errors"))

				return m
			}(),

			err: fmt.Errorf("s.security.CreateToken() errors"),
		},
		{
			name: "s.db.Login() errors",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())

				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(&model.UserPassword{}, nil).
					Once()

				m.EXPECT().BeginTx(mock.Anything).
					Return(nil, nil).
					Once()

				m.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(fmt.Errorf("s.db.Login() errors")).
					Once()

				return m
			}(),

			security: func() security.Security {
				m := securityMock.NewSecurity(s.T())
				m.EXPECT().ComparePassword(mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().CreateToken(mock.Anything).
					Return(token, nil)

				return m
			}(),

			err: fmt.Errorf("s.db.Login() errors"),
		},
		{
			name: "s.db.Commit() errors",

			db: func() db.DB {
				m := dbMock.NewDB(s.T())

				m.EXPECT().GetUserPassword(mock.Anything, mock.Anything).
					Return(&model.UserPassword{}, nil).
					Once()

				m.EXPECT().BeginTx(mock.Anything).
					Return(nil, nil).
					Once()

				m.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().Commit(mock.Anything).
					Return(fmt.Errorf("s.db.Commit() errors"))

				return m
			}(),

			security: func() security.Security {
				m := securityMock.NewSecurity(s.T())
				m.EXPECT().ComparePassword(mock.Anything, mock.Anything).
					Return(nil)

				m.EXPECT().CreateToken(mock.Anything).
					Return(token, nil)

				return m
			}(),

			err: fmt.Errorf("s.db.Commit() errors"),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			sec := &service{
				db:       test.db,
				security: test.security,
			}

			t, err := sec.Login(context.TODO(), &model.Login{})

			assert.Equal(s.T(), test.token, t)
			assert.Equal(s.T(), test.err, err)
		})
	}
}
