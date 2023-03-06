package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ilyadubrovsky/product-rest-api/internal/config"
	"github.com/ilyadubrovsky/product-rest-api/pkg/logging"
	"net/http"
)

const (
	productsURL = "/products"
	productURL  = "/products/:id"
)

type service interface {
	GetAllProducts(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, id string) (Product, error)
	CreateProduct(ctx context.Context, dto CreateProductDTO) (string, error)
	FullyUpdateProductByID(ctx context.Context, id string, dto UpdateProductDTO) error
	PartiallyUpdateProductByID(ctx context.Context, id string, dto UpdateProductDTO) error
	DeleteProductByID(ctx context.Context, id string) error
}

type Handler struct {
	logger  *logging.Logger
	cfg     *config.Config
	service service
}

func NewHandler(logger *logging.Logger, cfg *config.Config, service service) *Handler {
	return &Handler{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET(productsURL, h.getAllProducts)
	r.GET(productURL, h.getProduct)
	r.POST(productsURL, h.createProduct)
	r.PUT(productURL, h.fullyUpdateProduct)
	r.PATCH(productURL, h.partiallyUpdateProduct)
	r.DELETE(productURL, h.deleteProduct)
}

func (h *Handler) getAllProducts(c *gin.Context) {
	prdcts, err := h.service.GetAllProducts(context.TODO())
	if err != nil {
		switch err {
		case ErrNotFound:
			c.String(http.StatusNotFound, err.Error())
			return
		case ErrInternalServer:
			c.String(http.StatusInternalServerError, err.Error())
			return
		default:
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, prdcts)
}

func (h *Handler) getProduct(c *gin.Context) {
	id := c.Params.ByName("id")

	prdct, err := h.service.GetProductByID(context.TODO(), id)
	if err != nil {
		switch err {
		case ErrBadRequest:
			c.String(http.StatusBadRequest, err.Error())
			return
		case ErrNotFound:
			c.String(http.StatusNotFound, err.Error())
			return
		case ErrInternalServer:
			c.String(http.StatusInternalServerError, err.Error())
			return
		default:
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, prdct)
}

func (h *Handler) createProduct(c *gin.Context) {
	dto := CreateProductDTO{}
	if err := c.BindJSON(&dto); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateProduct(context.TODO(), dto)
	if err != nil {
		switch err {
		case ErrInternalServer:
			c.String(http.StatusInternalServerError, err.Error())
			return
		case ErrBadRequest:
			c.String(http.StatusBadRequest, err.Error())
			return
		default:
			c.String(http.StatusInternalServerError, err.Error())
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) fullyUpdateProduct(c *gin.Context) {
	id := c.Params.ByName("id")

	dto := UpdateProductDTO{}
	if err := c.BindJSON(&dto); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if dto.Type == "" || dto.InStock == 0 || dto.Name == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.FullyUpdateProductByID(context.TODO(), id, dto); err != nil {
		switch err {
		case ErrNotFound:
			c.String(http.StatusNotFound, err.Error())
			return
		case ErrBadRequest:
			c.String(http.StatusBadRequest, err.Error())
			return
		case ErrInternalServer:
			c.String(http.StatusInternalServerError, err.Error())
			return
		default:
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) partiallyUpdateProduct(c *gin.Context) {
	id := c.Params.ByName("id")

	dto := UpdateProductDTO{}
	if err := c.BindJSON(&dto); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.service.PartiallyUpdateProductByID(context.TODO(), id, dto); err != nil {
		switch err {
		case ErrNotFound:
			c.String(http.StatusNotFound, err.Error())
			return
		case ErrBadRequest:
			c.String(http.StatusBadRequest, err.Error())
			return
		case ErrInternalServer:
			c.String(http.StatusInternalServerError, err.Error())
			return
		default:
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) deleteProduct(c *gin.Context) {
	id := c.Params.ByName("id")

	if err := h.service.DeleteProductByID(context.TODO(), id); err != nil {
		switch err {
		case ErrNotFound:
			c.String(http.StatusNotFound, err.Error())
			return
		case ErrBadRequest:
			c.String(http.StatusBadRequest, err.Error())
			return
		case ErrInternalServer:
			c.String(http.StatusInternalServerError, err.Error())
			return
		default:
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Status(http.StatusNoContent)
}
