# 脚本公共 API 使用说明 (`api.ts`)

为了简化脚本开发并保持代码一致性，我们抽象了公共的 `api.ts` 模块。所有位于 `scripts/` 目录下的 TypeScript 脚本应当优先引入此模块进行 API 调用和参数解析。

## 引入模块

在脚本中引入常用的工具函数：

```typescript
import { getPayload, callApi, uploadResource, handleError } from './api';
```

## 核心功能介绍

### 1. `getPayload<T>()`

**作用**：解析命令行参数 `--payload='{...}'` 中的 JSON 内容并返回。

**示例**：
```typescript
interface MyPayload {
  name: string;
}
const payload = getPayload<MyPayload>();
console.log(payload.name);
```

### 2. `callApi(endpoint: string, options?: RequestInit)`

**作用**：封装了 `fetch` 请求，自动处理 `Authorization` (API Key) 头和 `Content-Type: application/json`。支持相对路径或完整 URL。

**示例**：
```typescript
const result = await callApi('/taskSets/v2', {
  method: 'POST',
  body: JSON.stringify(taskBody)
});
```

### 3. `uploadResource(urlOrPath: string, bucket?: string)`

**作用**：跨平台的资源上传辅助函数。支持输入本地文件路径或远程 HTTP URL，并将其上传到指定的 OSS Bucket（默认 `cloud-publish`），最后返回文件的统一 `key`。

**示例**：
```typescript
const key = await uploadResource('https://example.com/video.mp4');
```

### 4. `handleError(error: any, context: string)`

**作用**：统一的错误输出格式处理。当脚本发生异常时，它会输出 JSON 格式的错误信息并以状态码 `1` 退出程序。

**示例**：
```typescript
try {
  // logic...
} catch (error) {
  handleError(error, "submit the task");
}
```

## 脚本编写规范

1. **环境依赖**：确保系统环境变量中配置了 `YIXIAOER_API_KEY`。可选配置 `YIXIAOER_API_URL`。
2. **错误处理**：所有 `main` 函数内的逻辑应当被 `try...catch` 包裹，并调用 `handleError`。
3. **输出格式**：作为自动化工具，正常执行的输出内容必须为 JSON 格式（通过 `console.log(JSON.stringify(...))` 输出）。

---

*蚁小二 开源开发团队*
