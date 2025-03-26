package service

import (
	"github.com/Corray333/keep_it/internal/storage"
)

const (
	CodeRequestTypeSignUp = iota + 1
	CodeRequestTypeLogIn
	CodeRequestTypeChangePassword
)

type repository interface {
	storage.Transactioner

	verificationCodeGetter
	signUper
	logIner
	tokensRenewer

	userGetter
}

type UserService struct {
	transactioner          storage.Transactioner
	signUper               signUper
	logIner                logIner
	verificationCodeGetter verificationCodeGetter
	tokensRenewer          tokensRenewer
	userGetter             userGetter
}

type Option func(*UserService)

func WithTransactioner(tx storage.Transactioner) Option {
	return func(s *UserService) {
		s.transactioner = tx
	}
}

func WithSignUper(s signUper) Option {
	return func(u *UserService) {
		u.signUper = s
	}
}

func WithLogIner(l logIner) Option {
	return func(u *UserService) {
		u.logIner = l
	}
}

func WithVerificationCodeGetter(v verificationCodeGetter) Option {
	return func(s *UserService) {
		s.verificationCodeGetter = v
	}
}

func WithTokensRenewer(t tokensRenewer) Option {
	return func(s *UserService) {
		s.tokensRenewer = t
	}
}

func WithUserGetter(u userGetter) Option {
	return func(s *UserService) {
		s.userGetter = u
	}
}

func WithRepository(repo repository) Option {
	return func(s *UserService) {
		s.transactioner = repo
		s.signUper = repo
		s.logIner = repo
		s.verificationCodeGetter = repo
		s.tokensRenewer = repo
		s.userGetter = repo
	}
}

func New(options ...Option) *UserService {
	s := &UserService{}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *UserService) Run() {

}
