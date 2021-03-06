package bradescoNetEmpresa

import (
	"testing"
	"time"

	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
	. "github.com/smartystreets/goconvey/convey"
)

const baseMockJSON = `
{
    "bankNumber": 237,
   "authentication": {
            "Username": "",
            "Password": ""
        },
        "agreement": {
            "agreementNumber": 5822351,
            "wallet": 9,
            "agency": "1111",
            "account": "0062145"
        },
        "title": {
           
            "expireDate": "2050-12-30",
            "amountInCents": 200,
            "ourNumber": 12345678901,
            "instructions": "Não receber após a data de vencimento.",
            "documentNumber": "1234567890"
        },
        "recipient": {
            "name": "Empresa - Boletos",
            "document": {
                "type": "CNPJ",
                "number": "29799428000128"
            },
            "address": {
                "street": "Avenida Miguel Estefno, 2394",
                "complement": "Água Funda",
                "zipCode": "04301-002",
                "city": "São Paulo",
                "stateCode": "SP"
            }
        },
        "buyer": {
            "name": "Usuario Teste",
            "email": "p@p.com",
            "document": {
                "type": "CNPJ",
                "number": "29.799.428/0001-28"
            },
            "address": {
                "street": "Rua Teste",
                "number": "2",
                "complement": "SALA 1",
                "zipCode": "20931-001",
                "district": "Centro",
                "city": "Rio de Janeiro",
                "stateCode": "RJ"
            }
        }
}
`

func TestRegisterBoleto(t *testing.T) {
	env.Config(true, true, true)
	input := new(models.BoletoRequest)
	if err := util.FromJSON(baseMockJSON, input); err != nil {
		t.Fail()
	}
	bank := New()
	go mock.Run("9098")
	time.Sleep(2 * time.Second)

	Convey("Deve-se processar um boleto BradescoNetEmpresa com sucesso", t, func() {
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldNotBeEmpty)
		So(output.DigitableLine, ShouldNotBeEmpty)
		So(output.Errors, ShouldBeEmpty)
	})

	Convey("Deve-se exibir uma mensagem de erro, caso o registro não aconteça com sucesso", t, func() {
		input.Title.AmountInCents = 201
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})

	Convey("Deve-se exibir uma mensagem de erro, caso o status code for 500", t, func() {
		input.Title.AmountInCents = 202
		output, err := bank.ProcessBoleto(input)
		So(err, ShouldBeNil)
		So(output.BarCodeNumber, ShouldBeEmpty)
		So(output.DigitableLine, ShouldBeEmpty)
		So(output.Errors, ShouldNotBeEmpty)
	})
}
