export interface Store {
  id: number
  created_at: string
  updated_at: string
  name: string
  slug: string
  description: string
  is_active: boolean
  products?: Product[]
}

export interface Product {
  id: number
  created_at: string
  updated_at: string
  name: string
  description: string
  price: number
  stock: number
  is_active: boolean
  store_id: number
  store?: Store
}

export interface Order {
  id: number
  created_at: string
  updated_at: string
  user_id: number
  store_id: number
  status: string
  total_price: number
  items?: OrderItem[]
}

export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  quantity: number
  price: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}
