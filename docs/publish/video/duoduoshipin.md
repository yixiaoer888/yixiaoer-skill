# 多多视频 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| description | string | 否 | 多多视频描述 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| shopping_cart | object | 否 | 关联商品信息（购物车） | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Duoduoshipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PDD_VIDEO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 30
        },
        "contentPublishForm": {
          "formType": "task",
          "description": "这是多多视频的商品分享推荐内容。",
          "shopping_cart": {
            "goods_id": "998877",
            "source": "pdd"
          }
        }
      }
    ]
  }
}
```
