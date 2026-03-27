# 上传资源 (Upload Resource)

将图片、视频或其它素材上传至蚁小二云端存储 (OSS)，并获取用于发布的资源 Key。这是自动化发布的前置步骤。

## 场景描述 (Usage)

- "帮我把这个视频上传并给我资源 Key。"
- "上传这组图片作为多个账号的封面图。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `url` | `string` | **是** | 资源的远程 URL 或本地绝对路径 |
| `bucket` | `TeamBucketNamesEnum` | 否 | 存储桶。可选：`cloud-publish` (默认), `material-library`, `wechat-publish` |

### 存储桶说明 (Buckets)
- `cloud-publish`: 用于云端发布的视频、封面、图文素材。
- `material-library`: 素材库存储。
- `wechat-publish`: 微信公众号专用资产存储。

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/upload-resource.ts`
- **流程 (2-Step Upload)**:
  1. **申请地址**: 调用 `GET /api/storages/[bucket]/upload-url?fileKey=...`。
     - 获取 `serviceUrl`: 预签名的阿里云或私有云 PUT 写入地址。
     - 获取 `key`: 资源在 OSS 上的全局唯一路径。
  2. **物理写入**: 脚本将文件 Buffer 通过 `PUT` 请求直接发送至 `serviceUrl`。
- **调用示例**: 
  - `node scripts/upload-resource.ts --url="https://example.com/item.jpg" --bucket=cloud-publish`

## 输出结果 (Output)

输出生成的资源标识，供发布脚本引用：
```json
{
  "key": "cloud-publish/2026/03/26/66b2xxx/xxx.jpg",
  "name": "xxx.jpg"
}
```
**注意**: 在发布文章或视频时，请直接传入返回的 `key` 字符串作为封面或图片地址。

