package user

import (
	"my_web/backend/internal/response"
	"my_web/backend/internal/zerrors"
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
		h.Fail(ctx, zerrors.ErrRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		h.Fail(ctx, zerrors.ErrRequest)
		return
	}

	err := h.serv.Register(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.Fail(ctx, err)
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
		h.Fail(ctx, zerrors.ErrRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		h.Fail(ctx, zerrors.ErrRequest)
		return
	}
	token, err := h.serv.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.Fail(ctx, err)
		return
	}
	h.Success(ctx, gin.H{
		"token": token,
	})
}

func (h *Handler) getProfile(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	data, err := h.serv.getProfile(ctx, id)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, data)
}
