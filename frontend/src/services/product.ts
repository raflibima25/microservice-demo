// src/services/product.ts
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';

export interface Product {
    id: number;
    name: string;
    description: string;
    price: number;
    stock: number;
    created_at: string;
    updated_at: string;
}

export interface CreateProductRequest {
    name: string;
    description: string;
    price: number;
    stock: number;
}

export interface UpdateProductRequest {
    name?: string;
    description?: string;
    price?: number;
    stock?: number;
}

export interface ListProductsResponse {
    products: Product[];
    meta: {
        total: number;
        page: number;
        limit: number;
        total_pages: number;
    };
}

class ProductService {
    async createProduct(data: CreateProductRequest): Promise<Product> {
        const response = await axios.post<Product>(`${API_URL}/products`, data);
        return response.data;
    }

    async getProduct(id: number): Promise<Product> {
        const response = await axios.get<Product>(`${API_URL}/products/${id}`);
        return response.data;
    }

    async updateProduct(id: number, data: UpdateProductRequest): Promise<Product> {
        const response = await axios.put<Product>(`${API_URL}/products/${id}`, data);
        return response.data;
    }

    async deleteProduct(id: number): Promise<void> {
        await axios.delete(`${API_URL}/products/${id}`);
    }

    async listProducts(page: number = 1, limit: number = 10, search: string = ''): Promise<ListProductsResponse> {
        const response = await axios.get<ListProductsResponse>(`${API_URL}/products`, {
            params: { page, limit, search }
        });
        return response.data;
    }
}

export const productService = new ProductService();