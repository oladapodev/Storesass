import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { ArrowRight, Zap, Shield, BarChart3 } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { axiosInstance } from '@/lib/axios-client'
import { formatPrice } from '@/lib/utils'
import type { Store, Product, PaginatedResponse } from '@/types/storefront'

export function LandingPage() {
  const { data: storesData } = useQuery({
    queryKey: ['stores', 'featured'],
    queryFn: () =>
      axiosInstance.get<PaginatedResponse<Store>>('/api/v1/stores').then((r) => r.data),
  })

  const { data: productsData } = useQuery({
    queryKey: ['products', 'hot'],
    queryFn: () =>
      axiosInstance.get<PaginatedResponse<Product>>('/api/v1/products').then((r) => r.data),
  })

  const stores = storesData?.data ?? []
  const products = productsData?.data ?? []

  return (
    <div className="space-y-16">
      {/* Hero */}
      <section className="py-16 text-center space-y-6">
        <Badge variant="secondary" className="text-sm">
          Open Source SaaS Template
        </Badge>
        <h1 className="text-5xl font-bold tracking-tight">
          Launch your storefront
          <br />
          <span className="text-primary">in minutes</span>
        </h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          A production-ready multi-tenant storefront template built with Go, React, and modern tooling.
          Clean architecture, type-safe API, beautiful UI.
        </p>
        <div className="flex items-center justify-center gap-4">
          <Button asChild size="lg">
            <Link to="/stores">
              Browse Stores <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
          <Button variant="outline" size="lg" asChild>
            <a href="http://localhost:8080/swagger/index.html" target="_blank" rel="noreferrer">
              API Docs
            </a>
          </Button>
        </div>
      </section>

      {/* Features */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[
          {
            icon: Zap,
            title: 'Blazing Fast',
            description: 'Go backend with Redis caching. Auto-generated TanStack Query hooks for optimal data fetching.',
          },
          {
            icon: Shield,
            title: 'Type Safe',
            description: 'OpenAPI codegen pipeline: swag -> openapi-typescript -> Orval generates hooks and types automatically.',
          },
          {
            icon: BarChart3,
            title: 'Clean Architecture',
            description: 'Handler -> Service -> Repository separation. SOLID principles, domain-driven naming, testable layers.',
          },
        ].map(({ icon: Icon, title, description }) => (
          <Card key={title}>
            <CardHeader>
              <Icon className="h-8 w-8 text-primary mb-2" />
              <CardTitle>{title}</CardTitle>
              <CardDescription>{description}</CardDescription>
            </CardHeader>
          </Card>
        ))}
      </section>

      {/* Featured Stores */}
      {stores.length > 0 && (
        <section className="space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold">Featured Stores</h2>
            <Button variant="ghost" asChild>
              <Link to="/stores">
                View all <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {stores.slice(0, 2).map((store) => (
              <Card key={store.id} className="hover:shadow-md transition-shadow">
                <CardHeader>
                  <CardTitle>{store.name}</CardTitle>
                  <CardDescription>{store.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <Button asChild variant="outline" size="sm">
                    <Link to={`/stores/${store.slug}`}>
                      Visit Store <ArrowRight className="ml-2 h-3 w-3" />
                    </Link>
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
      )}

      {/* Hot Products */}
      {products.length > 0 && (
        <section className="space-y-6">
          <h2 className="text-2xl font-bold">Popular Products</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            {products.slice(0, 6).map((product) => (
              <Card key={product.id} className="hover:shadow-md transition-shadow">
                <CardHeader>
                  <CardTitle className="text-base">{product.name}</CardTitle>
                  <CardDescription className="line-clamp-2">{product.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <span className="font-bold text-lg">{formatPrice(product.price)}</span>
                    <Badge variant={product.stock > 0 ? 'default' : 'destructive'}>
                      {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
                    </Badge>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
      )}
    </div>
  )
}
