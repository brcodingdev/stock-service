package app_test

import (
	"bytes"
	"github.com/brcodingdev/stock-service/internal/pkg/app"
	mocks "github.com/brcodingdev/stock-service/internal/pkg/app/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandleStockRequest(t *testing.T) {
	cases := map[string]struct {
		key               string
		stockURL          string
		priceEvalExpected string
		csvResponse       string
	}{
		"call_stock_service_success": {
			key:               "aapl.us",
			stockURL:          "stooq.com",
			priceEvalExpected: "AAPL.US quote is $182.91 per share",
			csvResponse: `Symbol,Date,Time,Open,High,Low,Close,Volume
			AAPL.US,2023-09-06,22:00:11,188.4,188.85,181.47,182.91,67242023`,
		},
		"call_stock_service_not_found": {
			key:               "notfound.us",
			stockURL:          "stooq.com",
			priceEvalExpected: "notfound.us quote is not available",
			csvResponse: `Symbol,Date,Time,Open,High,Low,Close,Volume
			N/D,2023-09-06,22:00:11,N/D,N/D,N/D,N/D,N/D`,
		},
	}

	for caseTitle, tc := range cases {
		t.Run(caseTitle, func(t *testing.T) {
			httpMock := mocks.HTTPClient{}

			body := io.NopCloser(bytes.NewReader([]byte(tc.csvResponse)))

			response := http.Response{
				Status:     "200",
				StatusCode: http.StatusOK,
				Body:       body,
			}

			httpMock.On(
				"Do",
				mock.Anything,
			).Return(&response, nil)

			stockApp := app.NewStockApp(
				tc.stockURL, &httpMock,
			)

			res := stockApp.HandleStockRequest(tc.key)

			assert.Equal(
				t,
				strings.TrimSpace(tc.priceEvalExpected),
				strings.TrimSpace(res),
			)
		})
	}
}
