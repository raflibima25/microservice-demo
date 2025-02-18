import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface LoginCredentials {
    username: string;
    password: string;
}

export interface RegisterCredentials {
    username: string;
    email: string;
    password: string;
}

export interface User {
    id: number;
    username: string;
    email: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}

class AuthService {
    private setToken(token: string) {
        localStorage.setItem('token', token);
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    }

    private clearToken() {
        localStorage.removeItem('token');
        delete axios.defaults.headers.common['Authorization'];
    }

    async login(credentials: LoginCredentials): Promise<AuthResponse> {
        const response = await axios.post<AuthResponse>(`${API_URL}/auth/login`, credentials);
        this.setToken(response.data.token);
        return response.data;
    }

    async register(credentials: RegisterCredentials): Promise<AuthResponse> {
        const response = await axios.post<AuthResponse>(`${API_URL}/auth/register`, credentials);
        this.setToken(response.data.token);
        return response.data;
    }

    async logout() {
        try {
            await axios.post(`${API_URL}/auth/logout`);
        } finally {
            this.clearToken();
        }
    }

    async getCurrentUser(): Promise<User | null> {
        try {
            const token = localStorage.getItem('token');
            if (!token) return null;

            axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
            const response = await axios.get<User>(`${API_URL}/auth/me`)
            return response.data;
        } catch (error) {
            this.clearToken();
            return null;
        }
    }

    isAuthenticated(): boolean {
        return !!localStorage.getItem('token');
    }
}

export const authService = new AuthService();