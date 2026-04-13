# 📄 资源上传 动作 参数 (Upload Resource Action)

将图片、视频或其它素材上传至蚁小二云端存储 (OSS)，获取用于发布的资源 Key。这是所有自动化发布任务的必须前置步骤。

> [!CAUTION]
> **发布合规性禁令**：严禁在 `publish` 接口中直接透传外部 HTTP URL。所有外部资源必须通过本文档定义的 `upload` 动作转换为 `key` 后方可使用。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户提供了原始 URL 或本地路径，并准备进行发布或存入素材库。
- **典型提示词**：
  - “帮我传这张图到蚁小二”
  - “获取这个视频的上传 Key”
  - “上传这个封面”

## 2. 交互协议 (Interactive Protocol)

1. **类型匹配协议**：Agent 应识别文件后缀，并自动映射为标准的 `contentType` (如 `.png` -> `image/png`)。
2. **桶策略选择**：
   - 默认发布 -> `bucket: "cloud-publish"`。
   - 明确说“存入素材库” -> `bucket: "material-library"`。
3. **两步连动协议**：执行完上传并获得 Key 后，Agent 应立即根据场景引导用户进行下一步（发布或素材入库登记）。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `upload` | 固定值。 |
| **`url`** | `string` | **是** | - | 资源的远程 URL 或本地绝对路径。 |
| `bucket` | `string` | 否 | `cloud-publish` | 存储桶：`cloud-publish` (发布用) 或 `material-library` (入库用)。 |
| `contentType` | `string` | 否 | - | MIME 类型。必须与随后的物理上传 Header 强一致。 |

### 3.1 返回结果结构 (Response)

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `key` | `string` | **核心资源标识**。在后续发布表单中填入此值而非原始 URL。 |
| `uploadUrl` | `string` | 物理上传的目标地址。 |

## 4. 执行指令示例 (Command)

```bash
# 上传一张封面图
node scripts/api.ts --payload='{"action":"upload","url":"https://example.com/cover.jpg","bucket":"cloud-publish","contentType":"image/jpeg"}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **SignatureDoesNotMatch** | 获取 Key 时的 `contentType` 与实际物理上传时的 Header 不一致。 | 确保 Content-Type 字符串（含大小写）完全重合。 |
| **Timeout / 连接超时** | 蚁小二服务器无法访问该远程 URL，或文件过大。 | 检查 URL 是否公开，或尝试下载到本地后使用本地路径重传。 |
| **403 Forbidden** | API Key 权限不足。 | 重置并更新环境变量中的 `YIXIAOER_API_KEY`。 |

---
> [!IMPORTANT]
> **Key 的引用**：返回的 `key` 适用于所有发布组件（如 `images` 数组、`video.key` 或 `coverKey`）。
