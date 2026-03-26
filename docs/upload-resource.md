# 上传资源 (Upload Resource)

将图片、视频或其它素材上传至蚁小二云端存储 (OSS)，并获取用于发布的资源 Key。

## 场景描述 (Usage)

- "帮我把这个视频上传一下，我要发布到抖音。"
- "先上传这张图片作为多篇文章的封面。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `url` | `string` | 是 | 资源的远程 URL 或本地路径（本地需脚本运行环境支持读取） |
| `bucket` | `'cloud-publish'\|'material-library'` | 否 | 存储桶类型。默认 `'cloud-publish'` |

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/upload-resource.ts`
- **流程**:
  1. 识别 `url` 参数：若以 `http` 开头则抓取远程资源，否则视为本地路径读取文件。
  2. 调用 `GET /api/storages/[bucket]/upload-url` 获取预签名上传地址及 Key。
  3. 使用 `PUT` 请求将文件流发送至预签名地址。
- **调用示例**: 
  - 远程: `node upload-resource.ts --url="https://example.com/item.jpg"`
  - 本地: `node upload-resource.ts --url="./my-image.png"`

## 输出结果 (Output)

成功时输出包含资源 Key 的 JSON 对象：
```json
{
  "key": "cloud-publish/2026/03/26/xxx.jpg",
  "name": "xxx.jpg"
}
```
该 Key 可在后续的“发布百家号文章”等能力中作为参数使用。
