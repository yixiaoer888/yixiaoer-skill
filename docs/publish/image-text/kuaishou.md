# 快手图文发布参数 (Kuaishou Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`) | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |
| `visibleType` | `number` | **是** | 可见类型: 0-公开, 1-私密, 3-好友 | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |
| `collection` | `Object` | 否 | 合集信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["快手"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "KS_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>快手图文测试 <topic text='好心情'>#好心情</topic></p>",
          "images": [
            { "key": "img_key_1", "size": 1024, "width": 1080, "height": 1920 }
          ],
          "visibleType": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `KuaiShouDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaishou.dto.ts`
