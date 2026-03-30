# 多多视频发布参数 (DuoduoVideo)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | 否 | 多多视频描述 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `shopping_cart` | `object` | 否 | 关联商品信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["DuoduoVideo"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PDD_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "description": "多多视频内容 #拼多多"
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DuoDuoShiPinVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/duoduoshipin.dto.ts`
