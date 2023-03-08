package user

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/zoox"

	apiUser "github.com/go-zoox/connect/app/api/core/user"
)

const QRCODE_SERVER = "https://login.zcorky.com"

func GenerateDeviceUUID(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		response, err := fetch.Get(fmt.Sprintf("%s/api/qrcode/device/uuid", QRCODE_SERVER), &fetch.Config{
			Headers: fetch.Headers{
				"x-real-ip":       ctx.Get("x-forwarded-for"),
				"x-forwarded-for": ctx.Get("x-forwarded-for"),
				//
				"Accept": "application/json",
			},
			Query: fetch.Query{
				"client_id":     ctx.Query().Get("client_id").String(),
				"redirect_uri":  ctx.Query().Get("redirect_uri").String(),
				"response_type": "code",
				"state":         "_",
				"scope":         "qrcode",
			},
		})
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// example:
		// {
		//   "uuid": "326de618-a855-408c-a425-9b70a07ef82e",
		//   "url": "https://login.zcorky.com/qr/confirm?uuid=326de618-a855-408c-a425-9b70a07ef82e"
		// }
		ctx.String(200, response.String())
	}
}

func GetDeviceStatus(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		response, err := fetch.Get(fmt.Sprintf("%s/api/qrcode/device/status", QRCODE_SERVER), &fetch.Config{
			Headers: fetch.Headers{
				"x-real-ip":       ctx.Get("x-forwarded-for"),
				"x-forwarded-for": ctx.Get("x-forwarded-for"),
				//
				"Accept": "application/json",
			},
			Query: fetch.Query{
				"uuid": ctx.Query().Get("uuid").String(),
			},
		})
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// example:
		//	1. INIT
		// {
		// 	"owner": {
		// 		"name": "哆啦A梦认证",
		// 		"logo": "https://images.unsplash.com/photo-1506047610595-ab105976d1ef?ixlib=rb-1.2.1&q=80&fm=jpg&crop=entropy&cs=tinysrgb&w=400&fit=max&ixid=eyJhcHBfaWQiOjYwMDkyfQ",
		// 		"client_id": "xxxxxxxxxxxxxxx"
		// 	},
		// 	"client": {
		// 		"name": "Euno",
		// 		"logo": "https://resource.zcorky.com/api/open/v1/upload/4e0f9e63-526b-4138-8344-5504b2816aaa.png",
		// 		"client_id": "yyyyyyyyyyyyyyy"
		// 	},
		// 	"user": null,
		// 	"status": "INIT",
		// 	"uuid": "1edb8bcc-4e42-4768-8a01-416338886206",
		// 	"expiresAt": "2022-09-19T17:22:45.113Z"
		// }
		//
		//	2. scan
		// {
		//   "owner": {
		//     "name": "哆啦A梦认证",
		//     "logo": "https://images.unsplash.com/photo-1506047610595-ab105976d1ef?ixlib=rb-1.2.1&q=80&fm=jpg&crop=entropy&cs=tinysrgb&w=400&fit=max&ixid=eyJhcHBfaWQiOjYwMDkyfQ",
		//     "client_id": "xxxxxxxxxxxxxxx"
		//   },
		//   "client": {
		//     "name": "Euno",
		//     "logo": "https://resource.zcorky.com/api/open/v1/upload/4e0f9e63-526b-4138-8344-5504b2816aaa.png",
		//     "client_id": "yyyyyyyyyyyyyyy"
		//   },
		//   "user": {
		//     "nickname": "Zero",
		//     "avatar": "https://s3-imfile.feishucdn.com/static-resource/v1/v2_efc8fe50-072a-430b-9be8-23f32bfa664g~?image_size=72x72&cut_type=&quality=&format=image&sticker_format=.webp",
		//     "username": "zero@zero.com"
		//   },
		//   "status": "SCAN",
		//   "uuid": "1edb8bcc-4e42-4768-8a01-416338886206",
		//   "expiresAt": "2022-09-19T17:22:45.113Z"
		// }
		//
		// 	3. confirmed
		// {
		//   "authorization_code": "8f6b51b38c00c08092d757e5cd2f51de57c4e4ed",
		//   "state": "_"
		// }
		//
		//  4. error
		// {
		//   "code": 4004202,
		//   "message": "QRCode expired",
		//   "result": null
		// }
		ctx.String(200, response.String())
	}
}

type GetDeviceTokenReq struct {
	UUID string `json:"uuid"`
	Code string `json:"code"`
}

func GetDeviceToken(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		var body GetDeviceTokenReq
		if err := ctx.BindJSON(&body); err != nil {
			ctx.Fail(err, 500, "failed to parse body")
			return
		}

		response, err := fetch.Post(fmt.Sprintf("%s/api/qrcode/device/token", QRCODE_SERVER), &fetch.Config{
			Headers: fetch.Headers{
				"x-real-ip":       ctx.Get("x-forwarded-for"),
				"x-forwarded-for": ctx.Get("x-forwarded-for"),
				//
				"Accept": "application/json",
			},
			Body: map[string]string{
				"uuid": body.UUID,
				"code": body.Code,
			},
		})
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// example:
		// 1. success
		// {
		// 	"access_token": "qrcode_access_ce9e1aeae0328820849b514fb1",
		// 	"token_type": "bearer",
		// 	"expires_in": 3600,
		// 	"refresh_token": "qrcode_refresh_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// }
		//
		// 2. error
		// {
		//   "code": 5003003,
		//   "message": "QRCode expired"
		// }
		ctx.String(200, response.String())
	}
}

func GetUser(cfg *config.Config) zoox.HandlerFunc {
	// return func(ctx *zoox.Context) {
	// 	response, err := fetch.Get(fmt.Sprintf("%s/api/qrcode/device/user", QRCODE_SERVER), &fetch.Config{
	// 		Headers: fetch.Headers{
	// 			"Accept":        "application/json",
	// 			"Authorization": ctx.Get("Authorization"),
	// 		},
	// 	})
	// 	if err != nil {
	// 		ctx.String(500, err.Error())
	// 		return
	// 	}

	// 	// example:
	// 	// 1. success
	// 	// {
	// 	//   "nickname": "Zero",
	// 	//   "avatar": "https://s3-imfile.feishucdn.com/static-resource/v1/v2_efc8fe50-072a-430b-9be8-23f32bfa664g~?image_size=72x72&cut_type=&quality=&format=image&sticker_format=.webp",
	// 	//   "description": "",
	// 	//   "_id": "xxxxxxxxxxxxxxxxxxxxxxxxx",
	// 	//   "username": "zero@zero.com",
	// 	//   "email": "zero@zero.com"
	// 	// }
	// 	//
	// 	// 2. error
	// 	// {
	// 	//   "code": 5003003,
	// 	//   "message": "QRCode expired"
	// 	// }
	// 	ctx.String(200, response.String())
	// }

	return apiUser.New(cfg)
}
