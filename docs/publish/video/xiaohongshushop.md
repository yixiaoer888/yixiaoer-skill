# 小红书店铺 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| location | object | 否 | 视频位置，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| shoppingCart | object[] | 否 | 关联商品列表 | - |

## 2. 复杂对象结构

### PlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| id | string | 是 | ID |
| text | string | 是 | 文本内容 |
| raw | object | 是 | 平台原始数据 |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Xiaohongshushop"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_SHOP_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1440,
          "duration": 30
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "小红书店铺新品推荐",
          "description": "这是关于店铺新品发布的推广视频。",
          "shoppingCart": [
            {
              "goods_id": "112233",
              "price": 99.9
            }
          ]
        }
      }
    ]
  }
}
```
