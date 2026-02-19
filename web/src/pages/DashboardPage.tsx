import { useQuery } from '@tanstack/react-query'
import { Store, Package, TrendingUp, Activity } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { axiosInstance } from '@/lib/axios-client'
import { formatPrice } from '@/lib/utils'
import type { Store as StoreType, Product, PaginatedResponse } from '@/types/storefront'

export function DashboardPage() {
  const { data: storesData } = useQuery({
    queryKey: ['stores'],
    queryFn: () =>
      axiosInstance.get<PaginatedResponse<StoreType>>('/api/v1/stores').then((r) => r.data),
  })

  const { data: productsData } = useQuery({
    queryKey: ['products', 'hot'],
    queryFn: () =>
      axiosInstance.get<PaginatedResponse<Product>>('/api/v1/products').then((r) => r.data),
  })

  const stores = storesData?.data ?? []
  const products = productsData?.data ?? []
  const totalInventoryValue = products.reduce((sum, p) => sum + p.price * p.stock, 0)

  const stats = [
    {
      title: 'Total Stores',
      value: storesData?.total ?? 0,
      description: 'Active storefronts',
      icon: Store,
    },
    {
      title: 'Total Products',
      value: productsData?.total ?? 0,
      description: 'Listed products',
      icon: Package,
    },
    {
      title: 'Inventory Value',
      value: formatPrice(totalInventoryValue),
      description: 'Stock x Price',
      icon: TrendingUp,
    },
    {
      title: 'API Status',
      value: 'Online',
      description: 'Backend health',
      icon: Activity,
    },
  ]

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground mt-1">Overview of your storefront platform</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map(({ title, value, description, icon: Icon }) => (
          <Card key={title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{value}</div>
              <p className="text-xs text-muted-foreground mt-1">{description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Stores</CardTitle>
            <CardDescription>All active storefronts</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {stores.length === 0 ? (
                <p className="text-sm text-muted-foreground">No stores yet. Run: <code className="bg-muted px-1 rounded">make seed</code></p>
              ) : (
                stores.map((store) => (
                  <div key={store.id} className="flex items-center justify-between border-b pb-3 last:border-0">
                    <div>
                      <p className="font-medium">{store.name}</p>
                      <p className="text-sm text-muted-foreground">{store.slug}</p>
                    </div>
                    <Badge variant={store.is_active ? 'default' : 'secondary'}>
                      {store.is_active ? 'Active' : 'Inactive'}
                    </Badge>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Products</CardTitle>
            <CardDescription>Top products by price</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {products.length === 0 ? (
                <p className="text-sm text-muted-foreground">No products yet. Run: <code className="bg-muted px-1 rounded">make seed</code></p>
              ) : (
                [...products]
                  .sort((a, b) => b.price - a.price)
                  .slice(0, 5)
                  .map((product) => (
                    <div key={product.id} className="flex items-center justify-between border-b pb-3 last:border-0">
                      <div>
                        <p className="font-medium">{product.name}</p>
                        <p className="text-sm text-muted-foreground">Stock: {product.stock}</p>
                      </div>
                      <span className="font-bold">{formatPrice(product.price)}</span>
                    </div>
                  ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
