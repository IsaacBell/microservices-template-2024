package analyticsengine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"core/internal/auth"
	"core/internal/biz"
	"core/internal/constants"
	"core/internal/util"
	notifications_biz "core/pkg/notifications/biz"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/moesif/moesifapi-go"
	moesifModels "github.com/moesif/moesifapi-go/models"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type userCtxKey struct{}

type EventRequest struct {
	moesifModels.EventRequestModel
	Time time.Time `json:"time"`
}

type EventMetadata struct {
	notifications_biz.NotificationMetadata
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"` // e.g. "USD"
	Time     float64 `json:"time_seconds"`
}

type EventData struct {
	ActionName    string        `json:"action_name"`
	Request       EventRequest  `json:"request"`
	CompanyId     string        `json:"company_id"`
	UserId        string        `json:"user_id"`
	TransactionId string        `json:"transaction_id"`
	Metadata      EventMetadata `json:"metadata"`
}

var (
	moesif      moesifapi.API
	apiVersion  = strconv.Itoa(constants.ApiVersion)
	moesifAppId = os.Getenv("MOESIF_APP_ID")
	host        = "https://api.moesif.net"
	// headers     = map[string][]string{
	// 	"Content-Type":            []string{"application/json"},
	// 	"Accept":                  []string{"application/json"},
	// 	"Authorization":           []string{"Bearer " + moesifAppId},
	// 	"X-Moesif-Application-Id": []string{moesifAppId},
	// }
	MoesifOptions = map[string]interface{}{
		"Application_Id":           moesifAppId, // Header: X-Moesif-Application-Id
		"Identify_User":            "",
		"Capture_Outoing_Requests": true,
		"Debug":                    true,
	}
)

func UseMoesifClient() *moesifapi.API {
	if moesif == nil {
		moesif = moesifapi.NewAPI(os.Getenv("MOESIF_APP_ID"), &host, 100, 100, 10)
	}
	return &moesif
}

func UpsertMoesifUser(u *biz.User) error {
	UseMoesifClient()

	data := biz.UserToMoesifData(u)
	err := moesif.QueueUser(data)
	return err
}

type Event[Req, Reply any] struct {
	Ip        string
	UserId    string
	CompanyId string
	Url       string
	Method    string
	Operation string
	Error     error
	Request   Req
	Reply     Reply
	Headers   map[string]interface{}
	Timestamp time.Time
}

func NewEvent[Req any, Reply any](uid, cid, op, url, method string, req Req, reply Reply, err error) *Event[Req, Reply] {
	return &Event[Req, Reply]{
		UserId:    uid,
		CompanyId: cid,
		Operation: op,
		Url:       url,
		Method:    method,
		Request:   req,
		Reply:     reply,
		Error:     err,
	}
}

func (e *Event[Req, Reply]) User() *moesifModels.UserModel {
	str := ""
	uagent := ""
	tmp := e.Headers["user-agent"]
	if tmp != nil {
		for key, val := range e.Headers {
			if len(str) > 0 {
				str += "; "
			}
			str += key + ": " + val.([]string)[0]
		}
		uagent = str
	}
	return &moesifModels.UserModel{
		ModifiedTime:    &e.Timestamp,
		IpAddress:       &e.Ip,
		UserId:          e.UserId,
		CompanyId:       &e.CompanyId,
		UserAgentString: &uagent,
	}
}

func (e *Event[Req, Reply]) UpsertUser() error {
	UseMoesifClient()

	err := moesif.QueueUser(e.User())
	return err
}

func (event *Event[Req, Reply]) Submit(metadata map[string]interface{}) error {
	UseMoesifClient()
	MoesifOptions["Identify_User"] = event.UserId

	var responseTime time.Time
	if metadata["time"] == nil {
		responseTime = time.Now().Local()
		metadata["time"] = responseTime
	}
	status := 200
	if event.Error != nil {
		status = 500
	}

	e := moesifModels.EventModel{
		UserId:    &event.UserId,
		CompanyId: &event.CompanyId,
		Request: moesifModels.EventRequestModel{
			// Time: ,
			Uri:     event.Url,
			Verb:    event.Method,
			Headers: event.Headers,
			// Body:       event.Request,
			ApiVersion: &apiVersion,
			IpAddress:  &event.Ip,
		},
		Response: moesifModels.EventResponseModel{
			Time:   &responseTime,
			Status: status,
			// Headers: ,
			Body:      event.Reply,
			IpAddress: &event.Ip,
		},
		Metadata: metadata,
	}

	err := event.UpsertUser()
	if err != nil {
		return err
	}
	util.PrintLnInColor(util.AnsiColorMagenta, "::::::::::Analytics::::::::::\n->   UID: ", event.UserId,
		"\n->   Action: ", event.Url)

	err = moesif.QueueEvent(&e)
	if err == nil {
		util.PrintLnInColor(util.AnsiColorGreen, "->   OK - enqueued event:\n      ", util.AnsiColorMagenta, e)
	} else {
		util.PrintLnInColor(util.AnsiColorRed , "->   err queueing analytics: %v\n", err)
	}
	return err
}

type Reflect interface {
	ProtoReflect() protoreflect.Message
}

type RequestProto struct {
	Reflect
	Ok bool
}

func getUserJwtClaims(session any) (map[string]interface{}, error) {
	var err error
	if session == nil {
		return map[string]interface{}{}, nil
	}

	var u map[string]interface{}
	if session != nil {
		switch s := session.(type) {
		case []byte:
			err = json.Unmarshal(s, &u)
			if err != nil {
				util.PrintLnInColor(util.AnsiColorRed, "err: %v\n", err)
				return nil, err
			}
		case map[string]interface{}:
			u = s
		default:
			util.PrintLnInColor(util.AnsiColorRed, "unexpected session type: %T\n", s)
			return map[string]interface{}{}, nil
		}
	}

	// out := util.ConvertMapToStringVals(u)

	return u, nil
}

func setMetadataFromJwtClaims(
	claims map[string]interface{},
	metadata map[string]interface{},
) map[string]interface{} {
	for key, val := range claims {
		if val != nil && val != "" {
			switch v := val.(type) {
			case string:
				metadata[key] = v
			case bool:
				metadata[key] = strconv.FormatBool(v)
			case float64:
				metadata[key] = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				metadata[key] = strconv.Itoa(v)
			case int64:
				metadata[key] = strconv.FormatInt(v, 10)
			default:
				metadata[key] = fmt.Sprintf("%v", v)
			}
		}
	}
	return metadata
}

func getHeaders(headers interface{}) map[string]interface{} {
	if headers == nil {
		return map[string]interface{}{}
	}

	mapp := map[string]interface{}{}
	if reqHeaders, ok := headers.(transport.Header); ok {
		var keys []string = reqHeaders.Keys()
		for _, key := range keys {
			mapp[key] = reqHeaders.Values(key)
		}
	}

	return mapp
}

func MoesifMiddleware(authCtx *auth.AuthCtx) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {

		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			session := authCtx.Ctx.Value("session")

			// storing these as vars for convenience
			var (
				url        string
				operation  string
				method     string
				kind       string
				userId     string
				email      string
				first_name string
				last_name  string
				username   string
				companyId  string
				ip         string
				headers    map[string]interface{}
			)

			// Extract request headers
			if tr, ok := transport.FromServerContext(ctx); ok {
				kind = tr.Kind().String()
				fmt.Printf("kind: %v\n", kind)
				operation = tr.Operation()

				tmp := ctx.Value("method")
				if tmp != nil {
					method = tmp.(string)
				}
				if replyMsg, ok := reply.(RequestProto); ok {
					method = string(replyMsg.ProtoReflect().Descriptor().FullName())
					ok = replyMsg.Ok
				}

				userId = tr.ReplyHeader().Get("X-Auth-User-ID")
				companyId = tr.ReplyHeader().Get("X-Auth-Company-ID")

				url = tr.Endpoint()

				tmp = ctx.Value("host")
				if tmp != nil {
					host = ctx.Value("host").(string)
				}
				fmt.Printf("tr.RequestHeader(): %v\n", tr.RequestHeader())
				headers = getHeaders(tr.RequestHeader())

				ip = tr.RequestHeader().Get("X-Forwarded-For")
				if ip == "" {
					ip = tr.RequestHeader().Get("X-Real-IP")
				}
				if ip == "" {
					tmp := ctx.Value("request_ip")
					if tmp != nil {
						ip = tmp.(string)
					}
				}
			}

			info, err := getUserJwtClaims(session)
			if err == nil {
				if info["first_name"] != "" {
					first_name = info["first_name"].(string)
				}
				if info["last_name"] != "" {
					last_name = info["last_name"].(string)
				}
				if info["sub"] != "" {
					userId = info["sub"].(string)
				}
				if info["user_email"] != "" {
					email = info["user_email"].(string)
				}
				if info["username"] != "" {
					username = info["username"].(string)
				}
			} else {
				util.PrintLnInColor(util.AnsiColorRed, "Error retrieving user's JWT claims: ", err)
			}

			// Create metadata for the event
			tmp := map[string]interface{}{
				// amount: 100.0,
				// currency: "USD",
				"time_unix":  float64(time.Now().Local().Unix()),
				"kind":       kind,
				"username":   username,
				"first_name": first_name,
				"last_name":  last_name,
				"email":      email,
			}
			metadata := setMetadataFromJwtClaims(tmp, info)

			// Run the server operation
			reply, err = handler(ctx, req)

			userInfo, ok := ctx.Value("session").(map[string]string)
			if ok {
				fmt.Printf("userInfo[\"user_id\"]: %v\n", userInfo["user_id"])
				fmt.Printf("userInfo[\"user_id\"]: %v\n", userInfo["user_id"])
				if userInfo["company_id"] != "" {
					companyId = userInfo["company_id"]
				}
			}

			// Create an event with the extracted information
			event := NewEvent(userId, companyId, url, operation, method, req, reply, err)
			event.Ip = ip
			event.Operation = operation
			event.Headers = headers

			if err != nil {
				return nil, err
			}

			err = event.Submit(metadata)
			if err != nil {
				return nil, err
			}
			// ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

			return reply, err
		}
	}
}
