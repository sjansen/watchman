package protocol

// Request is the interface used to encode a request PDU.
type Request interface {
	Args() []interface{}
}

// Response is the interface common to all response PDUs.
type Response interface {
	PDU() ResponsePDU
	Version() string
	Warning() string
}

type response struct {
	pdu     ResponsePDU
	version string
	warning string
}

func (r *response) init(pdu ResponsePDU) {
	r.pdu = pdu
	if x, ok := pdu["version"]; ok {
		if version, ok := x.(string); ok {
			r.version = version
		}
	}
	if x, ok := pdu["warning"]; ok {
		if warning, ok := x.(string); ok {
			r.warning = warning
		}
	}
}

func (r *response) PDU() ResponsePDU {
	return r.pdu
}

// Version returns the Watchman server version.
func (r *response) Version() string {
	return r.version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (r *response) Warning() string {
	return r.warning
}

// ResponsePDU provides access to response data decoded to primitive Go values.
type ResponsePDU map[string]interface{}

// IsUnilateral indicates if the ResponsePDU was sent in response to
// the immediately previous request, or an older request such as the
// subscribe command.
func (pdu ResponsePDU) IsUnilateral() bool {
	if x, ok := pdu["unilateral"]; ok {
		if unilateral, ok := x.(bool); ok {
			return unilateral
		}
	}
	return false
}

// The ResponseTranslator type is an adapter that converts a
// ResponsePDU to a more exact Response.
type ResponseTranslator func(ResponsePDU) Response
