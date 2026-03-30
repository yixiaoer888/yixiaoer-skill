# 上传资源 (Upload Resource)

将图片、视频或其它素材上传至蚁小二云端存储 (OSS)，并获取用于发布的资源 Key。这是自动化发布的前置步骤。

## 场景描述 (Usage)

- "帮我把这个视频上传并给我资源 Key。"
- "上传这组图片作为多个账号的封面图。"

## 调用指令 (Command)

```bash
node scripts/upload-resource.ts --payload='{"url":"https://example.com/image.jpg","bucket":"cloud-publish"}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `url` | `string` | **是** | 资源的远程 URL 或本地绝对路径 |
| `bucket` | `string` | 否 | OSS 存储桶。可选：`cloud-publish` (默认), `material-library` |

## 输出结果 (Output)

输出生成的资源标识，供发布脚本引用：
```json
{
  "key": "cloud-publish/2026/03/26/66b2xxx/xxx.jpg",
  "name": "xxx.jpg"
}
```
**注意**: 在发布文章或视频时，请直接传入返回的 `key` 字符串作为封面或图片地址。

