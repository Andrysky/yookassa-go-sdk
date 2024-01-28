package yookassa

import (
	"testing"
)

// Для тестирования введите ваш идентификатор магазина и приватный ключ.
var shopid, key = "830911", "test_RYPy6ussTGgAD_MNaFDeGTV6G4wyqWW34H8_5ufNx7g"

// TestKassa_Ping проверяет работоспособность функции Ping.
func TestKassa_Ping(t *testing.T) {
	k := NewKassa(shopid, key)
	ok, err := k.Ping()
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Can't connect to yookassa")
	}
}

// TestKassa_SendPaymentConfig проверяет создание платежа.
func TestKassa_SendPaymentConfig(t *testing.T) {
	k := NewKassa(shopid, key)
	tAmount := Amount{
		Value:    "10.00",
		Currency: "RUB",
	}
	config := NewPaymentConfig(
		tAmount,
		Redirect{
			Type:      TypeRedirect,
			Locale:    "ru_RU",
			Enforce:   true,
			ReturnURL: "https://t.me/ranh_ranepa_bot",
		})

	payment, err := k.SendPaymentConfig(config)
	if err != nil {
		t.Error(err)
	}

	if payment.Amount != tAmount {
		t.Error("Got wrong response.")
	}
}

func TestKassa_GetPaymentInfo(t *testing.T) {
	k := NewKassa(shopid, key)
	tAmount := Amount{
		Value:    "10.00",
		Currency: "RUB",
	}
	config := NewPaymentConfig(
		tAmount,
		Redirect{
			Type:      TypeRedirect,
			Locale:    "ru_RU",
			Enforce:   true,
			ReturnURL: "https://t.me/ranh_ranepa_bot",
		})

	payment, err := k.SendPaymentConfig(config)
	if err != nil {
		t.Error(err)
	}

	id := payment.Id

	payment2, err := k.GetPayment(id)

	t.Log(payment)
	t.Log(payment2)

	if payment.Id != payment2.Id {
		t.Error("Got wrong response.")
	}
}

func TestKassa_GetPaymentErrorsControl(t *testing.T) {
	k := NewKassa(shopid, key)
	config := NewPaymentConfig(
		Amount{
			Value: "-5",
			Currency: "XMR",
		},
		Embedded{})

	_, err := k.SendPaymentConfig(config)
	if err == nil {
		t.Error("Error is not handled.")
	}
}

func TestKassa_GetPaymentErrorsControl(t *testing.T) {
	k := NewKassa(shopid, key)
	body := `{
		"type": "notification",
		"event": "payment.waiting_for_capture",
		"object": {
		  "id": "22d6d597-000f-5000-9000-145f6df21d6f",
		  "status": "waiting_for_capture",
		  "paid": true,
		  "amount": {
			"value": "2.00",
			"currency": "RUB"
		  },
		  "authorization_details": {
			"rrn": "10000000000",
			"auth_code": "000000",
			"three_d_secure": {
			  "applied": true
			}
		  },
		  "created_at": "2018-07-10T14:27:54.691Z",
		  "description": "Заказ №72",
		  "expires_at": "2018-07-17T14:28:32.484Z",
		  "metadata": {},
		  "payment_method": {
			"type": "bank_card",
			"id": "22d6d597-000f-5000-9000-145f6df21d6f",
			"saved": false,
			"card": {
			  "first6": "555555",
			  "last4": "4444",
			  "expiry_month": "07",
			  "expiry_year": "2021",
			  "card_type": "MasterCard",
			"issuer_country": "RU",
			"issuer_name": "Sberbank"
			},
			"title": "Bank card *4444"
		  },
		  "refundable": false,
		  "test": false
		}
	  }`

	event, err := k.GetResponseEvent(body)
	if err == nil {
		t.Error("Error is not handled.")
	}

	if event.Event == "payment.waiting_for_capture"{
		payment, err := k.GetResponseEventPayment(body)
		if err == nil {
			t.Error("Error is not handled.")
		}

		fmt.Println(payment)
	}
}
