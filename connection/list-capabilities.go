package connection

/*
$ watchman list-capabilities
{
    "version": "3.8.0",
    "capabilities": [
        "field-mode",
        "term-allof",
        "cmd-trigger"
    ]
}
*/

type ListCapabilitiesRequest struct{}

func (req *ListCapabilitiesRequest) Args() []interface{} {
	return []interface{}{"list-capabilities"}
}

type listCapabilitiesResponse struct {
	response
	Capabilities []string
}

type ListCapabilitiesResponse struct {
	listCapabilitiesResponse
}

func (res *ListCapabilitiesResponse) Version() string {
	return res.response.Version
}

func (res *ListCapabilitiesResponse) Warning() string {
	return res.response.Warning
}

func (res *ListCapabilitiesResponse) Capabilities() []string {
	return res.listCapabilitiesResponse.Capabilities
}
