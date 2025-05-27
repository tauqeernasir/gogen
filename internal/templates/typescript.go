package templates

import (
	"strings"
	"text/template"
)

// loadTypeScriptTemplates loads embedded TypeScript templates
func (tm *Manager) loadTypeScriptTemplates() error {
	templates := map[string]string{
		"typescript/package.json": `{
  "name": "{{.ProjectName | ToLower}}-client",
  "version": "{{.Version}}",
  "description": "{{.Description}}",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "prepublishOnly": "npm run build"
  },
  "dependencies": {
    "axios": "^1.6.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0"
  },
  "files": ["dist/"],
  "keywords": ["api", "client", "typescript", "axios"],
  "license": "MIT"
}`,

		"typescript/tsconfig.json": `{
  "compilerOptions": {
    "target": "ES2018",
    "module": "commonjs",
    "lib": ["ES2018"],
    "outDir": "./dist",
    "rootDir": "./",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["*.ts"],
  "exclude": ["node_modules", "dist"]
}`,

		"typescript/client": `import axios, { AxiosInstance, AxiosResponse, AxiosRequestConfig } from 'axios';
import { {{range $i, $t := .Types}}{{if $i}}, {{end}}{{$t.Name}}{{end}} } from './types';

export interface {{.ClientClassName}}Config {
  baseURL: string;
  timeout?: number;
  headers?: Record<string, string>;
}

export class {{.ClientClassName}} {
  private client: AxiosInstance;

  constructor(config: {{.ClientClassName}}Config) {
    this.client = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout || 30000,
      headers: {
        'Content-Type': 'application/json',
        ...config.headers,
      },
    });
  }

  public setAuthToken(token: string): void {
    this.client.defaults.headers.common['Authorization'] = ` + "`Bearer ${token}`" + `;
  }

  public removeAuthToken(): void {
    delete this.client.defaults.headers.common['Authorization'];
  }

{{range .Methods}}
  /**
   * {{.Summary}}
   * {{.Description}}
   */
  public async {{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.Name}}{{if not $p.Required}}?{{end}}: {{$p.Type}}{{end}}{{if .RequestBody}}{{if .Parameters}}, {{end}}data{{if not .RequestBody.Required}}?{{end}}: {{.RequestBody.Type}}{{end}}): Promise<{{.ResponseType}}> {
    const config: AxiosRequestConfig = {
      method: '{{.HTTPMethod}}',
      url: ` + "`{{.Path}}`" + `,{{if .RequestBody}}
      data: data,{{end}}{{if .Parameters}}
      params: { {{range $i, $p := .Parameters}}{{if eq $p.In "query"}}{{if $i}}, {{end}}{{$p.Name}}{{end}}{{end}} },{{end}}
    };

    const response: AxiosResponse<{{.ResponseType}}> = await this.client.request(config);
    return response.data;
  }
{{end}}
}`,

		"typescript/types": `// Generated types from OpenAPI specification
{{range .Types}}
{{if eq .Type "interface"}}export interface {{.Name}} {
{{range .Properties}}
  {{.Name}}{{if not .Required}}?{{end}}: {{.Type}};
{{end}}
}
{{else}}export type {{.Name}} = {{.Type}};{{end}}
{{end}}`,

		"typescript/index": `export { {{.ClientClassName}} } from './client';
export * from './types';`,

		"typescript/README.md": `# {{.ProjectName}} Client

TypeScript/JavaScript client for {{.ProjectName}} API.

## Installation

` + "```bash" + `
npm install {{.ProjectName | ToLower}}-client
` + "```" + `

## Usage

` + "```typescript" + `
import { {{.ClientClassName}} } from '{{.ProjectName | ToLower}}-client';

const client = new {{.ClientClassName}}({
  baseURL: '{{.BaseURL}}',
  timeout: 30000,
});

// Set authentication token if needed
client.setAuthToken('your-jwt-token');
` + "```" + `

## License

MIT`,
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	for name, content := range templates {
		tmpl, err := template.New(name).Funcs(funcMap).Parse(content)
		if err != nil {
			return err
		}
		tm.templates[name] = tmpl
	}

	return nil
}
