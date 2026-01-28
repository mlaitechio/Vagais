// API Types based on backend models and endpoints

export interface User {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role: string;
  is_active: boolean;
  email_verified: boolean;
  avatar: string;
  organization_id?: string;
  credits: number;
  last_login_at?: string;
  preferences: any;
  created_at: string;
  updated_at: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  description: string;
  website: string;
  logo: string;
  is_active: boolean;
  plan: string;
  created_at: string;
  updated_at: string;
}

export interface Agent {
  id: string;
  name: string;
  description: string;
  slug: string;
  version: string;
  status: 'draft' | 'published' | 'archived';
  type: string;
  category: string;
  tags: string[];
  config: any;
  llm_provider: string;
  llm_model: string;
  embedding_provider: string;
  embedding_model: string;
  creator_id: string;
  creator: User;
  organization_id?: string;
  organization?: Organization;
  is_public: boolean;
  is_enabled: boolean;
  price: number;
  currency: string;
  pricing_model: 'free' | 'one-time' | 'subscription';
  rating: number;
  review_count: number;
  usage_count: number;
  downloads: number;
  icon: string;
  screenshots: string[];
  documentation: string;
  repository: string;
  video_url: string;
  how_it_works: string;
  created_at: string;
  updated_at: string;
}

export interface Review {
  id: string;
  agent_id: string;
  agent: Agent;
  user_id: string;
  user: User;
  rating: number;
  title: string;
  content: string;
  is_verified: boolean;
  is_helpful: number;
  response?: string;
  created_at: string;
  updated_at: string;
}

export interface Execution {
  id: string;
  agent_id: string;
  agent: Agent;
  user_id: string;
  user: User;
  organization_id?: string;
  organization?: Organization;
  status: string;
  input: any;
  output: any;
  error?: string;
  duration: number;
  cost: number;
  credits_used: number;
  ip_address: string;
  user_agent: string;
  session_id: string;
  created_at: string;
  updated_at: string;
}

export interface Webhook {
  id: string;
  name: string;
  url: string;
  events: string[];
  secret?: string;
  is_active: boolean;
  organization_id?: string;
  organization?: Organization;
  user_id: string;
  user: User;
  last_triggered?: string;
  failure_count: number;
  headers: any;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: string;
  user_id: string;
  user: User;
  organization_id?: string;
  organization?: Organization;
  type: string;
  title: string;
  message: string;
  status: 'unread' | 'read';
  priority: 'low' | 'normal' | 'high';
  category: string;
  read_at?: string;
  metadata: any;
  created_at: string;
  updated_at: string;
}

// API Request/Response Types

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
  expires_at: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  password: string;
  organization_name?: string;
}

export interface CreateAgentRequest {
  name: string;
  description: string;
  category: string;
  tags: string[];
  config: any;
  llm_provider: string;
  llm_model: string;
  embedding_provider: string;
  embedding_model: string;
  is_public: boolean;
  price: number;
  currency: string;
  pricing_model: string;
  icon?: string;
  screenshots?: string[];
  documentation?: string;
  repository?: string;
}

export interface ExecuteAgentRequest {
  input: any;
  session_id?: string;
}

export interface CreateReviewRequest {
  agent_id: string;
  rating: number;
  title: string;
  content: string;
}

export interface SendNotificationRequest {
  user_id: string;
  type: string;
  title: string;
  message: string;
  priority: string;
  category: string;
  metadata?: any;
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ErrorResponse {
  error: string;
  message: string;
  status_code: number;
}

// Dashboard Stats
export interface DashboardStats {
  total_agents: number;
  total_executions: number;
  recent_activities: any[];
  usage_trends: any[];
}

// Marketplace Types
export interface MarketplaceFilters {
  category?: string;
  price_range?: [number, number];
  rating?: number;
  tags?: string[];
  provider?: string;
  pricing_model?: string;
}

export interface SearchAgentsRequest {
  query?: string;
  category?: string;
  tags?: string[];
  price_min?: number;
  price_max?: number;
  rating_min?: number;
  sort_by?: 'rating' | 'downloads' | 'created_at' | 'price';
  sort_order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}
