import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { ArrowRight, Store } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { axiosInstance } from '@/lib/axios-client'
import type { Store as StoreType, PaginatedResponse } from '@/types/storefront'

export function StoresPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['stores'],
    queryFn: () =>
      axiosInstance.get<PaginatedResponse<StoreType>>('/api/v1/stores').then((r) => r.data),
  })

  const stores = data?.data ?? []

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Stores</h1>
          <p className="text-muted-foreground mt-1">Browse all active storefronts</p>
        </div>
        <Badge variant="secondary">{data?.total ?? 0} stores</Badge>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader>
                <div className="h-5 bg-muted rounded w-3/4 mb-2" />
                <div className="h-4 bg-muted rounded w-full" />
              </CardHeader>
              <CardContent>
                <div className="h-8 bg-muted rounded w-1/3" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : stores.length === 0 ? (
        <div className="text-center py-16">
          <Store className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <p className="text-muted-foreground">No stores found. Run the seed command to add demo stores.</p>
          <code className="text-sm bg-muted px-2 py-1 rounded mt-2 inline-block">make seed</code>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {stores.map((store) => (
            <Card key={store.id} className="hover:shadow-md transition-shadow group">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle className="group-hover:text-primary transition-colors">
                    {store.name}
                  </CardTitle>
                  {store.is_active && <Badge variant="default">Active</Badge>}
                </div>
                <CardDescription>{store.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <Button asChild variant="outline" size="sm">
                  <Link to={`/stores/${store.slug}`}>
                    Browse Products <ArrowRight className="ml-2 h-3 w-3" />
                  </Link>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
