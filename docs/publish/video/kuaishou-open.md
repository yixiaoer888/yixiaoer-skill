# 快手-Open 视频发布参数 (Kuaishou-Open)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 快手标题 | - |
| `description` | `string` | 否 | 快手描述 | - |
| `declaration` | `number` | 否 | 声明: 1-内容为AI生成, 2-演绎情节, 3-个人观点 | - |
| `location` | `Object` | 否 | 快手视频位置 (`PlatformDataItem`) | - |
| `visibleType` | `number` | 否 | 可见类型: 0-公开, 1-私密, 3-好友可见 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `shopping_cart` | `object` | 否 | 关联商品 (与 `mini_app` 互斥) | - |
| `collection` | `Object` | 否 | 合集信息 | - |
| `mini_app` | `object` | 否 | 挂载小程序 | - |
| `nearby_show` | `boolean` | 否 | 是否同城展示 | `true` |
| `allow_same_frame` | `boolean` | 否 | 是否允许同框 | `false` |
| `allow_download` | `boolean` | 否 | 是否允许下载 | `false` |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["KuaishouOpen"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "KS_OPEN_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "快手Open视频",
          "description": "描述 #话题",
          "visibleType": 0,
          "declaration": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `KuaiShouVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaishou.dto.ts`
