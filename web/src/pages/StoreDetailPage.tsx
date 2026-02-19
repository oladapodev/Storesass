import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, ShoppingCart } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { axiosInstance } from '@/lib/axios-client'
import { formatPrice } from '@/lib/utils'
import type { Store, Product, PaginatedResponse } from '@/types/storefront'

type CartItem = Product & { quantity: number }

export function StoreDetailPage() {
  const { slug } = useParams<{ slug: string }>()
  const [cart, setCart] = useState<CartItem[]>([])
  const [searchQuery, setSearchQuery] = useState('')

  const { data: store, isLoading: storeLoading } = useQuery({
    queryKey: ['stores', slug],
    queryFn: () =>
      axiosInstance.get<Store>(`/api/v1/stores/${slug}`).then((r) => r.data),
    enabled: !!slug,
  })

  const { data: productsData } = useQuery({
    queryKey: ['stores', slug, 'products'],
    queryFn: () =>
      axiosInstance
        .get<PaginatedResponse<Product>>(`/api/v1/stores/${slug}/products`)
        .then((r) => r.data),
    enabled: !!slug,
  })

  const products = productsData?.data ?? []

  const filteredProducts = searchQuery
    ? products.filter(
        (p) =>
          p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          p.description?.toLowerCase().includes(searchQuery.toLowerCase()),
      )
    : products

  const addToCart = (product: Product) => {
    setCart((prev) => {
      const existing = prev.find((item) => item.id === product.id)
      if (existing) {
        return prev.map((item) =>
          item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item,
        )
      }
      return [...prev, { ...product, quantity: 1 }]
    })
  }

  const cartTotal = cart.reduce((sum, item) => sum + item.price * item.quantity, 0)
  const cartCount = cart.reduce((sum, item) => sum + item.quantity, 0)

  if (storeLoading) {
    return (
      <div className="space-y-6">
        <div className="h-8 bg-muted rounded w-1/3 animate-pulse" />
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader>
                <div className="h-5 bg-muted rounded w-3/4 mb-2" />
                <div className="h-4 bg-muted rounded w-full" />
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  if (!store) {
    return (
      <div className="text-center py-16">
        <p className="text-muted-foreground">Store not found.</p>
        <Button asChild variant="outline" className="mt-4">
          <Link to="/stores">
            <ArrowLeft className="mr-2 h-4 w-4" /> Back to Stores
          </Link>
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-4">
        <Button asChild variant="ghost" size="sm">
          <Link to="/stores">
            <ArrowLeft className="mr-2 h-4 w-4" /> All Stores
          </Link>
        </Button>
      </div>

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">{store.name}</h1>
          <p className="text-muted-foreground mt-1">{store.description}</p>
        </div>
        {cartCount > 0 && (
          <div className="flex items-center gap-3 bg-primary/10 border border-primary/20 rounded-lg px-4 py-2">
            <ShoppingCart className="h-5 w-5 text-primary" />
            <span className="font-medium">{cartCount} items</span>
            <span className="text-muted-foreground">|</span>
            <span className="font-bold text-primary">{formatPrice(cartTotal)}</span>
          </div>
        )}
      </div>

      <div className="max-w-sm">
        <Input
          placeholder="Search products..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {filteredProducts.length === 0 ? (
        <div className="text-center py-16 text-muted-foreground">
          <p>No products found.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          {filteredProducts.map((product) => {
            const cartItem = cart.find((item) => item.id === product.id)
            return (
              <Card key={product.id} className="flex flex-col">
                <CardHeader>
                  <div className="flex items-start justify-between gap-2">
                    <CardTitle className="text-base">{product.name}</CardTitle>
                    <Badge variant={product.stock > 0 ? 'secondary' : 'destructive'} className="shrink-0">
                      {product.stock > 0 ? `${product.stock} left` : 'Sold out'}
                    </Badge>
                  </div>
                  <CardDescription className="line-clamp-2">{product.description}</CardDescription>
                </CardHeader>
                <CardContent className="mt-auto">
                  <div className="flex items-center justify-between">
                    <span className="text-xl font-bold">{formatPrice(product.price)}</span>
                    <Button
                      size="sm"
                      onClick={() => addToCart(product)}
                      disabled={product.stock === 0}
                      variant={cartItem ? 'secondary' : 'default'}
                    >
                      {cartItem ? `In cart (${cartItem.quantity})` : 'Add to cart'}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
