// src/components/Products/ProductList.tsx
import React, { useEffect, useState } from 'react';
import { Product, productService } from '../../services/product';
import { Link } from 'react-router-dom';

const ProductList: React.FC = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [search, setSearch] = useState('');

    const fetchProducts = async () => {
        try {
            setLoading(true);
            const response = await productService.listProducts(page, 10, search);
            // Mengubah ini untuk menyesuaikan dengan response API
            setProducts(response.products);
            setTotalPages(response.meta.total_pages);
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchProducts();
    }, [page, search]);

    const handleDelete = async (id: number) => {
        if (window.confirm('Are you sure you want to delete this product?')) {
            try {
                await productService.deleteProduct(id);
                fetchProducts();
            } catch (err: any) {
                setError(err.message);
            }
        }
    };

    if (loading) return <div className="text-center">Loading...</div>;
    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="container mx-auto p-4">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">Products</h1>
                <Link 
                    to="/products/create"
                    className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                >
                    Add Product
                </Link>
            </div>

            <div className="mb-4">
                <input
                    type="text"
                    placeholder="Search products..."
                    className="w-full p-2 border rounded"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                />
            </div>

            {products.length === 0 ? (
                <div className="text-center text-gray-500">No products found</div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {products.map((product) => (
                        <div key={product.id} className="border rounded-lg p-4 shadow">
                            <h2 className="text-xl font-semibold">{product.name}</h2>
                            <p className="text-gray-600">{product.description}</p>
                            <div className="mt-2">
                                <p className="text-lg font-bold">${product.price.toLocaleString()}</p>
                                <p className="text-sm text-gray-500">Stock: {product.stock}</p>
                            </div>
                            <div className="mt-4 flex justify-end space-x-2">
                                <Link 
                                    to={`/products/${product.id}/edit`}
                                    className="bg-yellow-500 text-white px-3 py-1 rounded hover:bg-yellow-600"
                                >
                                    Edit
                                </Link>
                                <button
                                    onClick={() => handleDelete(product.id)}
                                    className="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
                                >
                                    Delete
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {totalPages > 1 && (
                <div className="mt-6 flex justify-center space-x-2">
                    <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                        className="px-4 py-2 border rounded disabled:opacity-50"
                    >
                        Previous
                    </button>
                    <span className="px-4 py-2">Page {page} of {totalPages}</span>
                    <button
                        onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                        disabled={page === totalPages}
                        className="px-4 py-2 border rounded disabled:opacity-50"
                    >
                        Next
                    </button>
                </div>
            )}
        </div>
    );
};

export default ProductList;