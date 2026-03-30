# 小红书商家号视频发布参数 (Xiaohongshushop)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `shoppingCart` | `Array` | 否 | 关联商品列表 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Xiaohongshushop"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_SHOP_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "商家视频",
          "description": "描述内容 #商家"
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `XiaoHongShuShopVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/xiaohongshushop.dto.ts`
