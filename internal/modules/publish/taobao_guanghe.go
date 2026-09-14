package publish

import (
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func ValidateTaobaoGuanghePublish(platform, publishType, channel string, payload map[string]interface{}) error {
	if platformutil.CanonicalKey(platform) != "taobaoguanghe" {
		return nil
	}
	if publishType != "video" && publishType != "imageText" {
		return taobaoGuangheError("taobao_guanghe_invalid_content_type", "淘宝光合仅支持 video 或 imageText 发布", map[string]interface{}{"type": publishType})
	}
	publishArgs := ExtractPublishArgs(payload)
	forms, _ := publishArgs["accountForms"].([]interface{})
	for i, rawForm := range forms {
		form, _ := rawForm.(map[string]interface{})
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		if publishType == "imageText" {
			if strings.TrimSpace(stringField(cpf, "title")) == "" && strings.TrimSpace(stringField(cpf, "desc")) == "" {
				return taobaoGuangheError("taobao_guanghe_content_empty", "淘宝光合图文标题和正文不能同时为空", map[string]interface{}{"accountForm": i})
			}
			images, _ := form["images"].([]interface{})
			if len(images) < 1 || len(images) > 9 {
				return taobaoGuangheError("taobao_guanghe_images_invalid", "淘宝光合图文必须包含 1 至 9 张图片", map[string]interface{}{"accountForm": i, "count": len(images)})
			}
			for imageIndex, rawImage := range images {
				image, _ := rawImage.(map[string]interface{})
				width, widthOK := numericValue(image["width"])
				height, heightOK := numericValue(image["height"])
				if !widthOK || !heightOK || width < 720 || height < 720 {
					return taobaoGuangheError("taobao_guanghe_image_dimensions", "淘宝光合图文图片宽高均不得小于 720 像素", map[string]interface{}{"accountForm": i, "image": imageIndex, "width": image["width"], "height": image["height"]})
				}
			}
		}
		if publishType == "video" {
			cover := objectField(form, "cover")
			if cover != nil && stringField(form, "coverKey") != stringField(cover, "key") {
				return taobaoGuangheError("taobao_guanghe_cover_mismatch", "淘宝光合视频 coverKey 必须与 cover.key 一致", map[string]interface{}{"accountForm": i})
			}
		}
		if err := validateTaobaoGuangheGoods(cpf["shopping_cart"], i); err != nil {
			return err
		}
	}
	return nil
}

func validateTaobaoGuangheGoods(value interface{}, accountIndex int) error {
	if value == nil {
		return nil
	}
	items, ok := value.([]interface{})
	if !ok {
		return taobaoGuangheError("taobao_guanghe_goods_invalid", "淘宝光合 shopping_cart 必须是数组", map[string]interface{}{"accountForm": accountIndex})
	}
	if len(items) > 6 {
		return taobaoGuangheError("taobao_guanghe_goods_limit", "淘宝光合每个账号最多关联 6 件商品", map[string]interface{}{"accountForm": accountIndex, "count": len(items)})
	}
	seen := map[string]bool{}
	for i, rawItem := range items {
		item, ok := rawItem.(map[string]interface{})
		id := stringField(item, "yixiaoerId")
		name := stringField(item, "yixiaoerName")
		_, rawOK := item["raw"].(map[string]interface{})
		if !ok || strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" || !rawOK {
			return taobaoGuangheError("taobao_guanghe_goods_invalid", "淘宝光合商品必须包含非空 yixiaoerId、yixiaoerName 和对象类型 raw", map[string]interface{}{"accountForm": accountIndex, "item": i})
		}
		if seen[id] {
			return taobaoGuangheError("taobao_guanghe_goods_invalid", "淘宝光合同一账号内商品 ID 不得重复", map[string]interface{}{"accountForm": accountIndex, "item": i, "yixiaoerId": id})
		}
		seen[id] = true
	}
	return nil
}

func taobaoGuangheError(code, message string, details interface{}) error {
	return yxerrors.New(yxerrors.ValidationType, code, message, details).
		WithCategory("taobao_guanghe_goods").
		WithHint("请重新执行淘宝光合商品查询或修正对应发布字段后，再运行 validate。")
}
