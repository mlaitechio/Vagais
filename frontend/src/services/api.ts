import axios, { AxiosInstance } from 'axios';
import {
  User,
  Organization,
  Agent,
  Review,
  Execution,
  Webhook,
  Notification,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  CreateAgentRequest,
  ExecuteAgentRequest,
  CreateReviewRequest,
  SendNotificationRequest,
  SearchAgentsRequest,
  ApiResponse,
  PaginatedResponse,
  DashboardStats,
} from '../types/api';

class ApiService {
  private api: AxiosInstance;
  private baseURL: string;

  constructor() {
    // Use relative path for production, localhost for development
    const isDevelopment = import.meta.env.DEV;
    this.baseURL = isDevelopment 
      ? 'http://localhost:8080/api/v1' 
      : '/api/v1';
    this.api = axios.create({
      baseURL: this.baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.api.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor to handle auth errors
    this.api.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401) {
          // Try to refresh token first
          const refreshToken = localStorage.getItem('refresh_token');
          if (refreshToken) {
            try {
              const response = await this.refreshToken();
              localStorage.setItem('token', response.access_token);
              localStorage.setItem('refresh_token', response.refresh_token);
              // Retry the original request
              const originalRequest = error.config;
              originalRequest.headers.Authorization = `Bearer ${response.access_token}`;
              return this.api(originalRequest);
            } catch (refreshError) {
              // If refresh fails, logout
              localStorage.removeItem('token');
              localStorage.removeItem('refresh_token');
              localStorage.removeItem('user');
              window.location.href = '/login';
            }
          } else {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.api.post<{success: boolean, data: LoginResponse}>('/auth/login', data);
    return response.data.data;
  }

  async register(data: RegisterRequest): Promise<ApiResponse<User>> {
    const response = await this.api.post<{success: boolean, data: User}>('/auth/register', data);
    return { data: response.data.data, message: 'Registration successful', success: true };
  }

  async refreshToken(): Promise<LoginResponse> {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }
    const response = await this.api.post<{success: boolean, data: LoginResponse}>('/auth/refresh', {
      refresh_token: refreshToken
    });
    return response.data.data;
  }

  async logout(): Promise<void> {
    await this.api.post('/auth/logout');
  }

  async validateToken(): Promise<ApiResponse<User>> {
    const response = await this.api.post<ApiResponse<User>>('/auth/validate');
    return response.data;
  }

  async forgotPassword(email: string): Promise<ApiResponse<void>> {
    const response = await this.api.post<ApiResponse<void>>('/auth/forgot-password', { email });
    return response.data;
  }

  async resetPassword(token: string, password: string): Promise<ApiResponse<void>> {
    const response = await this.api.post<ApiResponse<void>>('/auth/reset-password', { token, password });
    return response.data;
  }

  // User endpoints
  async getProfile(): Promise<User> {
    const response = await this.api.get<{success: boolean, data: User}>('/users/profile');
    return response.data.data;
  }

  async updateProfile(data: Partial<User>): Promise<User> {
    const response = await this.api.put<{success: boolean, data: User}>('/users/profile', data);
    return response.data.data;
  }

  async getUserStats(): Promise<DashboardStats> {
    const response = await this.api.get<{success: boolean, data: DashboardStats}>('/users/stats');
    return response.data.data;
  }

  async getUsers(page = 1, limit = 10): Promise<PaginatedResponse<User>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<User>}>(`/users?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getUser(id: string): Promise<User> {
    const response = await this.api.get<{success: boolean, data: User}>(`/users/${id}`);
    return response.data.data;
  }

  async updateUserRole(id: string, role: string): Promise<User> {
    const response = await this.api.put<{success: boolean, data: User}>(`/users/${id}/role`, { role });
    return response.data.data;
  }

  async activateUser(id: string): Promise<User> {
    const response = await this.api.put<{success: boolean, data: User}>(`/users/${id}/activate`);
    return response.data.data;
  }

  async deactivateUser(id: string): Promise<User> {
    const response = await this.api.put<{success: boolean, data: User}>(`/users/${id}/deactivate`);
    return response.data.data;
  }

  // Organization endpoints
  async getOrganizations(page = 1, limit = 10): Promise<PaginatedResponse<Organization>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Organization>}>(`/organizations?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getOrganization(id: string): Promise<Organization> {
    const response = await this.api.get<{success: boolean, data: Organization}>(`/organizations/${id}`);
    return response.data.data;
  }

  async createOrganization(data: Partial<Organization>): Promise<Organization> {
    const response = await this.api.post<{success: boolean, data: Organization}>('/organizations', data);
    return response.data.data;
  }

  async updateOrganization(id: string, data: Partial<Organization>): Promise<Organization> {
    const response = await this.api.put<{success: boolean, data: Organization}>(`/organizations/${id}`, data);
    return response.data.data;
  }

  async getOrganizationUsers(orgId: string, page = 1, limit = 10): Promise<PaginatedResponse<User>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<User>}>(`/organizations/${orgId}/users?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async inviteUserToOrganization(orgId: string, data: { email: string; role: string }): Promise<void> {
    await this.api.post(`/organizations/${orgId}/users`, data);
  }

  // Agent endpoints
  async getAgents(page = 1, limit = 10, filters?: any): Promise<PaginatedResponse<Agent>> {
    const params = new URLSearchParams({ page: page.toString(), limit: limit.toString() });
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Agent>}>(`/agents?${params}`);
    return response.data.data;
  }

  async getAgent(id: string): Promise<Agent> {
    const response = await this.api.get<{success: boolean, data: Agent}>(`/agents/${id}`);
    return response.data.data;
  }

  async createAgent(data: CreateAgentRequest): Promise<Agent> {
    const response = await this.api.post<{success: boolean, data: Agent}>('/agents', data);
    return response.data.data;
  }

  async updateAgent(id: string, data: Partial<CreateAgentRequest>): Promise<Agent> {
    const response = await this.api.put<{success: boolean, data: Agent}>(`/agents/${id}`, data);
    return response.data.data;
  }

  async deleteAgent(id: string): Promise<void> {
    await this.api.delete(`/agents/${id}`);
  }

  async enableAgent(id: string): Promise<Agent> {
    const response = await this.api.post<{success: boolean, data: Agent}>(`/agents/${id}/enable`);
    return response.data.data;
  }

  async disableAgent(id: string): Promise<Agent> {
    const response = await this.api.post<{success: boolean, data: Agent}>(`/agents/${id}/disable`);
    return response.data.data;
  }

  async executeAgent(id: string, data: ExecuteAgentRequest): Promise<Execution> {
    const response = await this.api.post<{success: boolean, data: Execution}>(`/agents/${id}/execute`, data);
    return response.data.data;
  }

  async getAgentCategories(): Promise<string[]> {
    const response = await this.api.get<{success: boolean, data: string[]}>('/agents/categories');
    return response.data.data;
  }

  async getAgentStats(id: string): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>(`/agents/${id}/stats`);
    return response.data.data;
  }

  // Marketplace endpoints
  async searchAgents(params: SearchAgentsRequest): Promise<PaginatedResponse<Agent>> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          value.forEach(v => queryParams.append(key, v.toString()));
        } else {
          queryParams.append(key, value.toString());
        }
      }
    });
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Agent>}>(`/marketplace/search?${queryParams}`);
    return response.data.data;
  }

  async getFeaturedAgents(): Promise<Agent[]> {
    const response = await this.api.get<{success: boolean, data: Agent[]}>('/marketplace/featured');
    return response.data.data;
  }

  async getTrendingAgents(): Promise<Agent[]> {
    const response = await this.api.get<{success: boolean, data: Agent[]}>('/marketplace/trending');
    return response.data.data;
  }

  async getMarketplaceCategories(): Promise<Record<string, number>> {
    const response = await this.api.get<{success: boolean, data: Record<string, number>}>('/marketplace/categories');
    return response.data.data;
  }

  async getMarketplaceStats(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/marketplace/stats');
    return response.data.data;
  }

  async listMarketplaceAgents(params: {
    page?: number;
    limit?: number;
    search?: string;
    category?: string;
    pricing?: string;
    rating?: number;
  }): Promise<PaginatedResponse<Agent>> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Agent>}>(`/marketplace/agents?${queryParams}`);
    return response.data.data;
  }

  async getMarketplaceAgent(id: string): Promise<Agent> {
    const response = await this.api.get<{success: boolean, data: Agent}>(`/marketplace/agents/${id}`);
    return response.data.data;
  }

  async tryMarketplaceAgent(id: string, input: any): Promise<any> {
    const response = await this.api.post<{success: boolean, data: any}>(`/marketplace/agents/${id}/try`, { input });
    return response.data.data;
  }

  async purchaseMarketplaceAgent(id: string, data: {
    pricing_tier: string;
    organization_id?: string;
  }): Promise<any> {
    const response = await this.api.post<{success: boolean, data: any}>(`/marketplace/agents/${id}/purchase`, data);
    return response.data.data;
  }

  async getAgentReviews(agentId: string, page = 1, limit = 10): Promise<{ reviews: Review[]; summary: any }> {
    const response = await this.api.get<{success: boolean, data: { reviews: Review[]; summary: any }}>(`/marketplace/agents/${agentId}/reviews?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async createAgentReview(agentId: string, data: {
    rating: number;
    title?: string;
    comment?: string;
  }): Promise<Review> {
    const response = await this.api.post<{success: boolean, data: Review}>(`/marketplace/agents/${agentId}/reviews`, data);
    return response.data.data;
  }

  // Public APIs (no auth required)
  async getPublicAgents(params: {
    page?: number;
    limit?: number;
    category?: string;
    sort_by?: string;
    sort_order?: string;
  }): Promise<PaginatedResponse<Agent>> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Agent>}>(`/public/agents?${queryParams}`);
    return response.data.data;
  }

  async searchPublicAgents(params: {
    q?: string;
    category?: string;
    min_rating?: number;
    max_price?: number;
    page?: number;
    limit?: number;
    sort_by?: string;
    sort_order?: string;
  }): Promise<PaginatedResponse<Agent>> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Agent>}>(`/public/agents/search?${queryParams}`);
    return response.data.data;
  }

  async getPublicCategories(): Promise<Record<string, number>> {
    const response = await this.api.get<{success: boolean, data: Record<string, number>}>('/public/agents/categories');
    return response.data.data;
  }

  // Review endpoints
  async getReviews(agentId: string, page = 1, limit = 10): Promise<PaginatedResponse<Review>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Review>}>(`/reviews/agent/${agentId}?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async createReview(data: CreateReviewRequest): Promise<Review> {
    const response = await this.api.post<{success: boolean, data: Review}>('/reviews', data);
    return response.data.data;
  }

  async getReview(id: string): Promise<Review> {
    const response = await this.api.get<{success: boolean, data: Review}>(`/reviews/${id}`);
    return response.data.data;
  }

  async updateReview(id: string, data: Partial<CreateReviewRequest>): Promise<Review> {
    const response = await this.api.put<{success: boolean, data: Review}>(`/reviews/${id}`, data);
    return response.data.data;
  }

  async deleteReview(id: string): Promise<void> {
    await this.api.delete(`/reviews/${id}`);
  }

  async markReviewHelpful(id: string): Promise<Review> {
    const response = await this.api.post<{success: boolean, data: Review}>(`/reviews/${id}/helpful`);
    return response.data.data;
  }

  // Runtime endpoints
  async getExecutions(page = 1, limit = 10): Promise<PaginatedResponse<Execution>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Execution>}>(`/runtime/executions?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getExecution(id: string): Promise<Execution> {
    const response = await this.api.get<{success: boolean, data: Execution}>(`/runtime/executions/${id}`);
    return response.data.data;
  }

  async cancelExecution(id: string): Promise<void> {
    await this.api.post(`/runtime/executions/${id}/cancel`);
  }

  async getExecutionStats(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/runtime/executions/stats');
    return response.data.data;
  }

  async getAgentExecutionStats(agentId: string): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>(`/runtime/agents/${agentId}/executions/stats`);
    return response.data.data;
  }

  async getActiveExecutions(): Promise<Execution[]> {
    const response = await this.api.get<{success: boolean, data: Execution[]}>('/runtime/executions/active');
    return response.data.data;
  }

  // Notification endpoints
  async getNotifications(page = 1, limit = 10): Promise<PaginatedResponse<Notification>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Notification>}>(`/notifications?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getNotification(id: string): Promise<Notification> {
    const response = await this.api.get<{success: boolean, data: Notification}>(`/notifications/${id}`);
    return response.data.data;
  }

  async sendNotification(data: SendNotificationRequest): Promise<Notification> {
    const response = await this.api.post<{success: boolean, data: Notification}>('/notifications', data);
    return response.data.data;
  }

  async markNotificationAsRead(id: string): Promise<Notification> {
    const response = await this.api.put<{success: boolean, data: Notification}>(`/notifications/${id}/read`);
    return response.data.data;
  }

  async markAllNotificationsAsRead(): Promise<void> {
    await this.api.put('/notifications/read-all');
  }

  async deleteNotification(id: string): Promise<void> {
    await this.api.delete(`/notifications/${id}`);
  }

  async getUnreadCount(): Promise<number> {
    const response = await this.api.get<{success: boolean, data: number}>('/notifications/unread-count');
    return response.data.data;
  }

  async sendBulkNotification(data: { user_ids: string[]; notification: SendNotificationRequest }): Promise<void> {
    await this.api.post('/notifications/bulk', data);
  }

  async getNotificationStats(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/notifications/stats');
    return response.data.data;
  }

  // Integration endpoints
  async getWebhooks(page = 1, limit = 10): Promise<PaginatedResponse<Webhook>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Webhook>}>(`/integrations/webhooks?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getWebhook(id: string): Promise<Webhook> {
    const response = await this.api.get<{success: boolean, data: Webhook}>(`/integrations/webhooks/${id}`);
    return response.data.data;
  }

  async createWebhook(data: Partial<Webhook>): Promise<Webhook> {
    const response = await this.api.post<{success: boolean, data: Webhook}>('/integrations/webhooks', data);
    return response.data.data;
  }

  async updateWebhook(id: string, data: Partial<Webhook>): Promise<Webhook> {
    const response = await this.api.put<{success: boolean, data: Webhook}>(`/integrations/webhooks/${id}`, data);
    return response.data.data;
  }

  async deleteWebhook(id: string): Promise<void> {
    await this.api.delete(`/integrations/webhooks/${id}`);
  }

  async getLLMProviders(): Promise<any[]> {
    const response = await this.api.get<{success: boolean, data: any[]}>('/integrations/llm-providers');
    return response.data.data;
  }

  async testLLMConnection(data: any): Promise<any> {
    const response = await this.api.post<{success: boolean, data: any}>('/integrations/llm-providers/test', data);
    return response.data.data;
  }

  async getIntegrationStats(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/integrations/stats');
    return response.data.data;
  }

  // Admin endpoints
  async getAllUsers(page = 1, limit = 10): Promise<PaginatedResponse<User>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<User>}>(`/admin/users?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getAllOrganizations(page = 1, limit = 10): Promise<PaginatedResponse<Organization>> {
    const response = await this.api.get<{success: boolean, data: PaginatedResponse<Organization>}>(`/admin/organizations?page=${page}&limit=${limit}`);
    return response.data.data;
  }

  async getBlockedDomains(): Promise<string[]> {
    const response = await this.api.get<{success: boolean, data: string[]}>('/admin/domains/blocked');
    return response.data.data;
  }

  async addBlockedDomain(domain: string): Promise<void> {
    await this.api.post('/admin/domains/blocked', { domain });
  }

  async removeBlockedDomain(domain: string): Promise<void> {
    await this.api.delete(`/admin/domains/blocked/${domain}`);
  }

  async getSystemStats(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/admin/stats');
    return response.data.data;
  }

  async getSystemHealth(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/admin/health');
    return response.data.data;
  }

  async getSystemMetrics(): Promise<any> {
    const response = await this.api.get<{success: boolean, data: any}>('/admin/metrics');
    return response.data.data;
  }

  // Health check
  async healthCheck(): Promise<any> {
    const response = await this.api.get<any>('/health');
    return response.data;
  }
}

export const apiService = new ApiService();
export default apiService;
