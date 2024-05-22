package types

type SetLockedRequest struct {
	Request
	Params SetLockedRequestParams `json:"params"`
}

type SetLockedRequestParams struct {
	IsLocked bool `json:"is_locked"`
}

type SetLockedResponse struct {
	Response
	Result []struct {
		LockState string `json:"lock_state"`
	} `json:"result"`
}
