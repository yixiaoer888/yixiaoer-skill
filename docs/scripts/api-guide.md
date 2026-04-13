# 📄 脚本公共 API 开发指南 (Script API Guide)

本指南定义了 `scripts/api.ts` 公共模块的使用规范。所有开发者在为蚁小二 (YiXiaoEr) 编写新 Action 脚本时，应当优先引入此模块以确保“零依赖”和“自诊断”目标的实现。

> [!IMPORTANT]
> **安全准则**：严禁在脚本中硬编码 `API Key`。必须通过环境变量 `YIXIAOER_API_KEY` 进行读取。

## 1. 核心工具函数 (Core Functions)

### 1.1 `getPayload<T>()`
解析命令行参数 `--payload='{...}'` 中的 JSON 内容。Agent 在调用脚本时应遵循此传参规范。

### 1.2 `callApi(endpoint, options)`
封装 `fetch` 请求，自动注入鉴权头和内容类型。支持错误自动捕获。

### 1.3 `uploadResource(urlOrPath, ...)`
资产搬运的核心函数。支持将远程 URL 或本地路径物理转储到 OSS，并返回系统互认的 `key`。

## 2. 交互协议 (Interactive Protocol)

1. **输入一致性**：脚本必须支持解析标准化 Payload 结构。
2. **输出标准化**：脚本成功执行后，必须输出符合 DTO 定义的 JSON 字符串。
3. **静默执行**：除了最终的结果 JSON，严禁在 `stdout` 中打印调试信息，以防 Agent 解析失败。

## 3. 错误分类原则 (Error Standards)

| 错误代码 (errorCode) | 含义 | 处理策略 |
| :--- | :--- | :--- |
| `YIXIAOER_USAGE_ERR` | 参数或 JSON 格式错误 | Agent 重新读文档核对字段必填性。 |
| `YIXIAOER_REMOTE_ERR` | 远端 API 或后端逻辑错误 | 按 [🛡️ 排障手册](../troubleshooting-guide.md) 排查。 |
| `YIXIAOER_AUTH_ERR` | 鉴权或账号状态失效 | 引导用户通过 `accounts` 更新状态。 |

## 4. 示例代码 (Code Example)

```typescript
import { getPayload, callApi, handleError } from './api';

async function main() {
  try {
    const payload = getPayload<{ action: string }>();
    const result = await callApi('/some-endpoint');
    console.log(JSON.stringify(result));
  } catch (err) {
    handleError(err, "Execute custom action", "YIXIAOER_USAGE_ERR");
  }
}
main();
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **TS2307: Cannot find module** | 未安装 peer dependencies。 | 虽然项目倾向零依赖，但开发环境需确保 `typescript` 和 `@types/node` 可用。 |
| **API 返回 401** | `YIXIAOER_API_KEY` 配置错误。 | 重新检查环境变量配置。 |

---
> [!NOTE]
> **维护说明**：本指南随 `api.ts` 版本同步更新。若接口发生 Breaking Change，必须在此文档首部进行显著声明。
