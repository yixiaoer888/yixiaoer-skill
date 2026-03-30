# 小红书视频发布参数 (Xiaohongshu Video)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 笔记标题 (最多20字) | - |
| `description` | `string` | 否 | 笔记内容及话题标签 | - |
| `declaration` | `number` | 否 | 内容申明: 1-虚构演绎, 2-笔记含AI合成内容 | - |
| `type` | `number` | 否 | 创作类型: 1-原创, 0-不申明 | `0` |
| `visibleType` | `number` | **是** | 可见性: 0-公开, 1-私密, 3-好友 | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `collection` | `Object` | 否 | 笔记所属合集 | - |
| `group` | `Object` | 否 | 关联群聊信息 | - |
| `bind_live_info` | `Object` | 否 | 关联直播预告 | - |
| `shopping_cart` | `Array` | 否 | 笔记挂载商品信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["小红书"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "我的小红书视频",
          "description": "这是内容描述 #话题",
          "visibleType": 0,
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `XiaoHongShuVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/xiaohongshu.dto.ts`
