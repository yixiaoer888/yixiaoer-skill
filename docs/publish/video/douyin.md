# 抖音 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 抖音视频标题 | - |
| description | string | 否 | 抖音视频描述 | - |
| horizontalCover | object | 否 | 抖音视频横板封面，使用 `OldCover` 结构 | - |
| statement | number | 否 | 抖音视频声明：3-内容由AI生成, 4-可能引人不适, 5-虚构演绎仅供娱乐, 6-危险行为请勿模仿 | - |
| location | object | 否 | 抖音视频位置，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| allow_save | number | 否 | 保存权限：0-不允许, 1-允许 | 0 |
| shoppingCart | object[] | 否 | 购物车列表 | - |
| groupShopping | object | 否 | 团购信息 | - |
| collection | object | 否 | 合集信息 | - |
| sub_collection | object | 否 | 合集选集 | - |
| sync_apps | object[] | 否 | 同时发布到的应用列表 | - |
| hot_event | object | 否 | 热点事件 | - |
| challenge | object | 否 | 挑战/话题 | - |
| mini_app | object | 否 | 挂载小程序（与购物车互斥） | - |
| music | object | 否 | 背景音乐信息 | - |
| cooperation_info | object | 否 | 共创信息 | - |
| game | object | 否 | 游戏挂载信息 | - |

## 2. 复杂对象结构

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| width | number | 是 | 宽度 |
| height | number | 是 | 高度 |
| size | number | 是 | 文件大小 (Byte) |
| path | string | 否 | 文件绝对路径 (与 `key` 二选一) |
| key | string | 否 | OSS 对象存储 Key (与 `path` 二选一) |

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
  "platforms": ["Douyin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DOUYIN_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "记录美好生活",
          "description": "这是我在抖音的第一条视频 #美好生活 #见闻",
          "statement": 3,
          "allow_save": 1,
          "location": {
            "id": "123",
            "text": "上海市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```
