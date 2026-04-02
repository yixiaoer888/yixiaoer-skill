# 上传资源 (Upload Resource)

将图片、视频或其它素材上传至蚁小二云端存储 (OSS)，并获取用于发布的资源 Key。这是自动化发布的前置步骤。

## 场景描述 (Usage)

- "帮我把这个视频上传并给我资源 Key。"
- "上传这组图片作为多个账号的封面图。"

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"upload","url":"https://example.com/image.jpg","bucket":"cloud-publish"}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `url` | `string` | **是** | 资源的远程 URL 或本地绝对路径 |
| `bucket` | `string` | **是** | OSS 存储桶。固定值为 `cloud-publish` |
| `contentType` | `string` | 否 | 资源的 MIME 类型。如果通过 URL 上传，程序会尝试自动识别，建议显式指定 (如 `video/mp4`, `image/png`) |
| `size` | `number` | 否 | 资源大小 (字节)。用于素材库容量检查 (如果适用) |

## 输出结果 (Output)

输出生成的资源标识，供发布脚本引用：
```json
{
  "key": "cloud-publish/2026/03/26/66b2xxx/xxx.jpg",
  "name": "xxx.jpg"
}
```
**注意**: 在发布文章或视频时，请直接传入返回的 `key` 字符串作为封面或图片地址。

> [!IMPORTANT]
> **发布合规性提醒**:
> 所有的封面图、图文图片、视频文件均**严禁直接使用外部网络 URL**，必须通过本项目提供的上传接口进行处理并获取 `key` 后进行发布。不遵守此规范将直接导致任务失败。

