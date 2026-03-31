# 小红书 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 视频标题 | - |
| description | string | 否 | 视频描述 | - |
| declaration | number | 否 | 内容类型申明：1-虚构演绎仅供娱乐, 2-笔记含 AI 合成内容 | - |
| type | number | 否 | 创作类型：1-原创, 0-不申明 | 0 |
| location | object | 否 | 视频位置，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| visibleType | number | 是 | 可见类型：0-公开, 1-私密, 3-好友可见 | 0 |
| collection | object | 否 | 合集信息 | - |
| group | object | 否 | 群聊信息 | - |
| bind_live_info | object | 否 | 直播预告信息 | - |
| shopping_cart | object[] | 否 | 关联商品信息 | - |

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
  "platforms": ["Xiaohongshu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1440,
          "duration": 30
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "小红书笔记标题",
          "description": "这是在小红书分享的一段精彩视频 #好物分享 #生活",
          "type": 1,
          "visibleType": 0,
          "declaration": 2
        }
      }
    ]
  }
}
```
