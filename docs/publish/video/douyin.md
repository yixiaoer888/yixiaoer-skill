# 抖音视频发布参数 (Douyin Video)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 抖音视频标题 | - |
| `description` | `string` | **是** | 抖音视频描述 | - |
| `horizontalCover` | `Object` | **是** | 视频横板封面 (`key`, `width`, `height`, `size`) | - |
| `statement` | `number` | 否 | 创作声明: 3-AI生成, 4-引人不适, 5-虚构娱乐, 6-危险动作 | - |
| `location` | `Object` | 否 | 平台物理位置对象 (`yixiaoerId`, `raw`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (10位) | - |
| `allow_save` | `number` | 否 | 允许保存: 0-不允许, 1-允许 | `0` |
| `music` | `Object` | 否 | 音乐挂载对象 | - |
| `mini_app` | `Object` | 否 | 挂载小程序信息 | - |
| `collection` | `Object` | 否 | 挂载合集信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "video": {
          "key": "video_oss_key",
          "size": 1024,
          "width": 1920,
          "height": 1080
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "测试视频标题",
          "description": "这是抖音视频描述 #话题",
          "horizontalCover": {
            "key": "cover_oss_key",
            "size": 100,
            "width": 1280,
            "height": 720
          }
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DouYinVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`
