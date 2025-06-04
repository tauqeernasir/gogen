# GoGen - OpenAPI Client Generator

A fast, type-safe API client generator that transforms OpenAPI 3.0 specifications into production-ready client libraries.

## 🚀 Core Features

- **Multi-Source Input**: Load specs from files, URLs, or stdin
- **Type-Safe Generation**: Fully typed TypeScript clients
- **NestJS Compatible**: Handles NestJS Swagger schemas seamlessly
- **Zero Configuration**: Works out of the box with sensible defaults

## 📦 Installation

```bash
go install github.com/yourusername/gogen@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/gogen.git
cd gogen && go build -o gogen main.go
```

## 🛠️ Usage

### Generate Client

```bash
# From local file
gogen -spec ./openapi.json -name myapi

# From URL
gogen -spec https://api.example.com/openapi.json -name myapi

# From stdin
curl -s https://api.example.com/openapi.json | gogen -spec - -name myapi
```

### Use Generated Client

```typescript
import { MyApiClient } from "./generated-client";

const client = new MyApiClient({
  baseURL: "https://api.example.com",
});

// Set auth token
client.setAuthToken("your-token");

// Make requests with grouped parameters
await client.createUser(
  { page: 1, limit: 10 }, // query params
  { "x-api-key": "key" }, // headers
  { name: "John", email: "john@example.com" } // request body
);
```

## 📋 Command Options

```bash
-spec    OpenAPI spec (file path, URL, or '-' for stdin)
-name    Project name (required)
-output  Output directory (default: ./generated-client)
-lang    Language: typescript (default: typescript)
-templates Custom templates directory
-prettier Run prettier after generation (default: true)
```

## 🎯 Generated Output

```
generated-client/
├── package.json # NPM package
├── client.ts # Main client class
├── types.ts # TypeScript interfaces
└── index.ts # Exports
```

**Type-safe methods:**

```typescript
// Generated from OpenAPI spec
interface User { id: number; name: string; email?: string }

async createUser(
  query?: { page?: number; limit?: number },
  headers?: { 'x-api-key'?: string },
  data: CreateUserRequest
): Promise<User>
```

## 🚀 Quick Start

```bash
# 1. Generate client
gogen -spec https://petstore.swagger.io/v2/swagger.json -name petstore

# 2. Install & build
cd generated-client && npm install && npm run build

# 3. Use in your project
import { PetstoreClient } from './generated-client';
const client = new PetstoreClient({ baseURL: 'https://petstore.swagger.io/v2' });
```

## 📄 License

MIT License - see [LICENSE](https://github.com/tauqeernasir/gogen/blob/main/LICENSE) for details.

---
