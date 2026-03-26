# 发布视频 (Publish Video Engine)

全平台视频发布统一接口。支持将视频同步分发到抖音、快手、视频号、B站等多个平台的多个账号。

## 场景描述 (Usage)
- "帮我把这个视频发布到我的抖音和快手。"
- "使用视频 (Key: xxx) 发布到视频号，标题为：我们的故事。"

## 参数定义 (Parameters)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `video_key` | `string` | 是 | 视频素材的 OSS Key (通过上传资源获得) |
| `title` | `string` | 是 | 视频标题 |
| `description` | `string` | 否 | 视频描述/简介 |
| `platforms` | `string[]` | 是 | 目标平台列表 |
| `account_ids` | `string[]` | 否 | 目标账号列表 |

## 脚本逻辑 (Backend)
- **脚本路径**: `../scripts/publish-video.ts`
- **核心逻辑**: 存证视频素材 -> 构建平台 DTO -> 提交任务集。
