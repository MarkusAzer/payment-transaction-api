package fly

import "testing"

func TestSetStatus(t *testing.T) {
	tt := []struct {
		name   string
		ft     flypayATransaction
		status string
		msg    string
	}{
		{
			"Status Code 1",
			flypayATransaction{Amount: 1000, Currency: "AUD", StatusCode: 1, OrderReference: "2e58bd43-0001", TransactionID: "flypay-a-0001"},
			"authorised",
			"StatusCode 1 should be authorised",
		},
		{
			"Status Code 2",
			flypayATransaction{Amount: 1000, Currency: "AUD", StatusCode: 2, OrderReference: "2e58bd43-0001", TransactionID: "flypay-a-0001"},
			"decline",
			"StatusCode 2 should be decline",
		},
		{
			"Status Code 3",
			flypayATransaction{Amount: 1000, Currency: "AUD", StatusCode: 3, OrderReference: "2e58bd43-0001", TransactionID: "flypay-a-0001"},
			"refunded",
			"StatusCode 3 should be refunded",
		},
	}

	for _, tc := range tt {
		tc.ft.setStatus()
		if tc.ft.Status != tc.status {
			t.Errorf("%v :- %v", tc.name, tc.msg)
		}
	}
}
