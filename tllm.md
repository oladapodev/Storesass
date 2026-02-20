# Orval & Type-Safe API Guide

## 📍 Where are Orval types and hooks?
The generated code lives in `web/src/api/`. 
- **Hooks & API Methods**: Found in directories like `web/src/api/stores/` and `web/src/api/products/`.
- **Domain Models**: Shared TypeScript interfaces describing your database entities (e.g., `Store`, `Product`) are in `web/src/api/model/`.

## 🛠️ How are they used?
Instead of manually configuring `useQuery` and defining types for every fetch, you import the generated hooks directly:

```tsx
// ❌ OLD WAY (Manual & Error Prone)
const { data } = useQuery({
  queryKey: ['stores'],
  queryFn: () => axios.get<PaginatedResponse<Store>>('/api/v1/stores')
})

// ✅ NEW WAY (Orval Generated)
import { useGetStores } from '@/api/stores/stores'

const { data } = useGetStores({ limit: 10 }) // Fully typed parameters and response!
```

## 🚀 How this improves the project
1. **End-to-End Type Safety**: If you change a field name in a Go struct (e.g. `internal/domain/models.go`) and run `make codegen`, TypeScript will immediately show red squiggles everywhere that field is used in the frontend. No more runtime "undefined" errors.
2. **Automated Documentation**: The types are derived directly from your Swagger annotations in Go. Your code *is* your documentation.
3. **Zero Boilerplate**: You don't have to write Axios wrappers, interface definitions, or TanStack Query configurations. Orval creates standard, battle-tested hooks for Every. Single. Endpoint.
4. **Consistency**: Whether you are fetching products or creating a store, the API interaction pattern is identical across the entire team.
5. **LLM Friendly**: Because the entire API surface is defined in TypeScript files, AI assistants (like Copilot) can perfectly understand your backend capabilities without you having to explain them.

## 🔄 The Codegen Loop
1. Edit Go `models.go` or `handler.go`.
2. Run `make codegen`.
3. Use the updated hooks/types in your React components.
