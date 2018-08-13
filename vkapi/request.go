package vkapi

type RequestParams map[string]interface{}

type VKAPIRequest interface {
	MethodName() string
	Params() (RequestParams, error)
}
