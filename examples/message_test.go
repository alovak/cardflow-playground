package examples

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

var spec *iso8583.MessageSpec = &iso8583.MessageSpec{
	Name: "ISO 8583 CardFlow Playgroud ASCII Specification",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Length:      8,
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		}),
		2: field.NewString(&field.Spec{
			Length:      19,
			Description: "Primary Account Number (PAN)",
			Enc:         encoding.BCD,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewString(&field.Spec{
			Length:      6,
			Description: "Amount",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
			Pad:         padding.Left('0'),
		}),
		4: field.NewString(&field.Spec{
			Length:      12,
			Description: "Transmission Date & Time", // YYMMDDHHMMSS
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		5: field.NewString(&field.Spec{
			Length:      2,
			Description: "Approval Code",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		6: field.NewString(&field.Spec{
			Length:      6,
			Description: "Authorization Code",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		7: field.NewString(&field.Spec{
			Length:      3,
			Description: "Currency",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		8: field.NewString(&field.Spec{
			Length:      4,
			Description: "Card Verification Value (CVV)",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		9: field.NewString(&field.Spec{
			Length:      4,
			Description: "Card Expiration Date",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
		10: field.NewComposite(&field.Spec{
			Length:      999,
			Description: "Acceptor Information",
			Pref:        prefix.ASCII.LLL,
			Tag: &field.TagSpec{
				Length: 2,
				Enc:    encoding.ASCII,
				Sort:   sort.StringsByInt,
			},
			Subfields: map[string]field.Field{
				"01": field.NewString(&field.Spec{
					Length:      99,
					Description: "Merchant Name",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}),
				"02": field.NewString(&field.Spec{
					Length:      4,
					Description: "Merchant Category Code (MCC)",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"03": field.NewString(&field.Spec{
					Length:      10,
					Description: "Merchant Postal Code",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}),
				"04": field.NewString(&field.Spec{
					Length:      299,
					Description: "Merchant Website",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LLL,
				}),
			},
		}),
		11: field.NewString(&field.Spec{
			Length:      6,
			Description: "Systems Trace Audit Number (STAN)",
			Enc:         encoding.BCD,
			Pref:        prefix.BCD.Fixed,
		}),
	},
}

func TestMessagePackingAndUnpacking(t *testing.T) {
	// We use field tags to map the struct fields to the ISO 8583 fields
	type AcceptorInformation struct {
		MerchantName         string `iso8583:"01"`
		MerchantCategoryCode string `iso8583:"02"`
		MerchantPostalCode   string `iso8583:"03"`
		MerchantWebsite      string `iso8583:"04"`
	}

	type AuthorizationRequest struct {
		MTI                 string               `iso8583:"0"`
		PAN                 string               `iso8583:"2"`
		Amount              int64                `iso8583:"3"`
		TransactionDatetime string               `iso8583:"4"`
		Currency            string               `iso8583:"7"`
		CVV                 string               `iso8583:"8"`
		ExpirationDate      string               `iso8583:"9"`
		AcceptorInformation *AcceptorInformation `iso8583:"10"`
		STAN                string               `iso8583:"11"`
	}

	// Create a new message
	requestMessage := iso8583.NewMessage(spec)

	// use time from our example
	timeFromExample := "240812160140"
	processingTime, err := time.Parse("060102150405", timeFromExample)
	require.NoError(t, err)

	// Set the message fields
	err = requestMessage.Marshal(&AuthorizationRequest{
		MTI:                 "0100",
		PAN:                 "4242424242424242",
		Amount:              1000,
		TransactionDatetime: processingTime.Format("060102150405"),
		Currency:            "840",
		CVV:                 "7890",
		ExpirationDate:      "2512",
		AcceptorInformation: &AcceptorInformation{
			MerchantName:         "Merchant Name",
			MerchantCategoryCode: "1234",
			MerchantPostalCode:   "1234567890",
			MerchantWebsite:      "https://www.merchant.com",
		},
		STAN: "000001",
	})
	require.NoError(t, err)

	// Pack the message
	packed, err := requestMessage.Pack()
	require.NoError(t, err)

	// Unpack the message
	requestMessage = iso8583.NewMessage(spec)
	err = requestMessage.Unpack(packed)
	require.NoError(t, err)

	// Unmarshal the message fields
	var authorizationRequest AuthorizationRequest
	err = requestMessage.Unmarshal(&authorizationRequest)
	require.NoError(t, err)

	// Check the message fields
	require.Equal(t, "0100", authorizationRequest.MTI)
	require.Equal(t, "4242424242424242", authorizationRequest.PAN)
	require.Equal(t, int64(1000), authorizationRequest.Amount)
	require.Equal(t, timeFromExample, authorizationRequest.TransactionDatetime)
	require.Equal(t, "840", authorizationRequest.Currency)
	require.Equal(t, "7890", authorizationRequest.CVV)
	require.Equal(t, "2512", authorizationRequest.ExpirationDate)
	require.Equal(t, "Merchant Name", authorizationRequest.AcceptorInformation.MerchantName)
	require.Equal(t, "1234", authorizationRequest.AcceptorInformation.MerchantCategoryCode)
	require.Equal(t, "1234567890", authorizationRequest.AcceptorInformation.MerchantPostalCode)
	require.Equal(t, "https://www.merchant.com", authorizationRequest.AcceptorInformation.MerchantWebsite)
	require.Equal(t, "000001", authorizationRequest.STAN)

	// Here is the example of the packed message
	examplePackedMessage := "010073E000000000000031364242424242424242001000240812160140084078902512303636303131334D65726368616E74204E616D653032313233343033313031323334353637383930303430323468747470733A2F2F7777772E6D65726368616E742E636F6D000001"

	// Check the packed message
	// using %X to convert the byte slice to a hex string in uppercase
	require.Equal(t, examplePackedMessage, fmt.Sprintf("%X", packed))

	// to make it right, let's filter the value of CVV field when we output it
	filterCVV := iso8583.FilterField("8", iso8583.FilterFunc(func(in string, data field.Field) string {
		if len(in) == 0 {
			return in
		}
		return in[0:1] + strings.Repeat("*", len(in)-1)
	}))

	// don't forget to apply default filter
	filters := append(iso8583.DefaultFilters(), filterCVV)

	err = iso8583.Describe(requestMessage, os.Stdout, filters...)
	require.NoError(t, err)
}
