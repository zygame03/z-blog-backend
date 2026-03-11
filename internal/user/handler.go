package user

import (
	"errors"
	"my_web/backend/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	response.BaseHandler
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/user")
	r.POST("/register", h.register)
	r.POST("/login", h.login)
	r.GET("/profile/:id", h.getProfile)
}

func (h *Handler) register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Fail(ctx, response.ErrRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		h.Fail(ctx, response.ErrRequest)
		return
	}

	err := h.serv.Register(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExist) {
			h.Fail(ctx, response.ErrUserExist)
			return
		}
		h.Fail(ctx, response.ErrDBOp)
		return
	}
	h.Success(ctx, "")
}

func (h *Handler) login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Fail(ctx, response.ErrRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		h.Fail(ctx, response.ErrRequest)
		return
	}
	token, err := h.serv.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.Fail(ctx, response.ErrUserExist)
			return
		}
		if errors.Is(err, ErrInvalidPassword) {
			h.Fail(ctx, response.ErrPassword)
			return
		}
		h.Fail(ctx, response.ErrDBOp)
		return
	}
	h.Success(ctx, gin.H{
		"token": token,
	})
}

func (h *Handler) getProfile(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.Fail(ctx, response.ErrRequest)
		return
	}

	data, err := h.serv.getProfile(ctx, id)
	if err != nil {
		h.Fail(ctx, response.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}
