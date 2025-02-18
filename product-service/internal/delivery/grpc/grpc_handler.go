package grpc

import (
	"context"
	pb "grpc/pb/product"
	"product-service/internal/domain"
	"time"
)

type GRPCProductHandler struct {
	pb.UnimplementedProductServiceServer
	productUseCase domain.ProductUseCase
}

func NewGRPCProductHandler(productUseCase domain.ProductUseCase) *GRPCProductHandler {
	return &GRPCProductHandler{productUseCase: productUseCase}
}

func (h *GRPCProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := h.productUseCase.Create(
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
	)
	if err != nil {
		return nil, err
	}

	return convertToProtoProduct(product), nil
}

func (h *GRPCProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := h.productUseCase.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return convertToProtoProduct(product), nil
}

func (h *GRPCProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := h.productUseCase.List(req.Page, req.PerPage, req.Search)
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertToProtoProduct(&product)
	}

	meta := &pb.Meta{
		Total:      int32(total),
		Page:       req.Page,
		PerPage:    req.PerPage,
		TotalPages: int32((total + int64(req.PerPage) - 1) / int64(req.PerPage)),
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Meta:     meta,
	}, nil
}

func (h *GRPCProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	product, err := h.productUseCase.Update(
		req.Id,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
	)
	if err != nil {
		return nil, err
	}

	return convertToProtoProduct(product), nil
}

func (h *GRPCProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.productUseCase.Delete(req.Id)
	if err != nil {
		return &pb.DeleteProductResponse{Success: false}, err
	}

	return &pb.DeleteProductResponse{Success: true}, nil
}

// serve starts the gRPC server
func (h *GRPCProductHandler) Serve(address string) error {
	server := NewGRPCProductServer(address)
	server.RegisterServices(h)
	return server.Start()
}

// helper func to conver domain Product to proto Product
func convertToProtoProduct(product *domain.Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}
