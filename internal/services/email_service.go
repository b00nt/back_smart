// internal/services/email_service.go
package services

import (
	"back/internal/models"
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

func SendFeedback(city string, feedback models.Feedback) error {
	var from string
	var to string
	var psswd string
	if city == "moscow" {
		from = "zhmihov.valery@yandex.ru"
		to = "zhmihov.valery@yandex.ru"
		psswd = "dzmeuihwmgwwqyzl"
	} else if city == "saratov" {
		from = "zhmihov.valery@yandex.ru"
		to = "zhmihov.valery@yandex.ru"
		psswd = "dzmeuihwmgwwqyzl"
	}
	subject := "Данные из обратной связи"
	body := fmt.Sprintf("Обратный звонок\n\nИмя: %s\nТелефон: %s\nГород: %s", feedback.Name, feedback.Telephone, feedback.City)

	// Set up the email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// Configure the SMTP dialer
	dialer := gomail.NewDialer("smtp.yandex.ru", 465, from, psswd)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendOrder(city string, order models.Order, customer models.CustomerInfo, items []models.OrderItem) error {
	var from string
	var to string
	var psswd string
	if city == "moscow" {
		from = "zhmihov.valery@yandex.ru"
		to = "zhmihov.valery@yandex.ru"
		psswd = "dzmeuihwmgwwqyzl"
	} else if city == "saratov" {
		from = "zhmihov.valery@yandex.ru"
		to = "zhmihov.valery@yandex.ru"
		psswd = "dzmeuihwmgwwqyzl"
	}

	// Email subject
	subject := fmt.Sprintf("Новый заказ #%d", order.ID)

	// Generate email body
	body := buildOrderEmailBody(order, customer, items)

	// Set up the email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body) // Use HTML for better formatting

	// Configure the SMTP dialer
	dialer := gomail.NewDialer("smtp.yandex.ru", 465, from, psswd)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Helper function to build the email body
func buildOrderEmailBody(order models.Order, customer models.CustomerInfo, items []models.OrderItem) string {
	// Start with customer info
	emailBody := `
	<h2>Детали нового заказа</h2>
	<p><strong>ID заказа:</strong> ` + fmt.Sprintf("%d", order.ID) + `</p>
	<p><strong>Дата заказа:</strong> ` + order.OrderDate.Format("2006-01-02 15:04:05") + `</p>
	<h3>Информация заказчика</h3>
	<ul>
		<li><strong>ФИО:</strong> ` + customer.FullName + `</li>
		<li><strong>Email:</strong> ` + customer.Email + `</li>
		<li><strong>Телефон:</strong> ` + customer.TelephoneNumber + `</li>
		<li><strong>Адрес:</strong> ` + formatAddress(customer) + `</li>
		<li><strong>Комментарии:</strong> ` + emptyOrValue(customer.Comment) + `</li>
	</ul>`

	// Include alternative contact details, if provided
	if customer.AnotherFullName != "" || customer.AnotherTelephoneNumber != "" {
		emailBody += `
		<h4>Альтернативная информация заказчика</h4>
		<ul>
			<li><strong>ФИО:</strong> ` + emptyOrValue(customer.AnotherFullName) + `</li>
			<li><strong>Телефон:</strong> ` + emptyOrValue(customer.AnotherTelephoneNumber) + `</li>
		</ul>`
	}

	// Add order summary
	emailBody += `
	<h3>Заказ</h3>
	<p><strong>Общая сумма:</strong> ` + fmt.Sprintf("%.2f", order.TotalAmount) + `</p>`

	// Add ordered items
	emailBody += `
	<h3>Заказанные товары</h3>
	<table border="1" cellpadding="8" cellspacing="0">
		<tr>
			<th>Название товара</th>
			<th>Количество</th>
			<th>Цена</th>
			<th>Всего</th>
			<th>Характеристики</th>
		</tr>`

	for _, item := range items {
		total := float64(item.Quantity) * item.Price

		// Start building the row for the current item
		emailBody += `
		<tr>
			<td>` + fmt.Sprintf(item.Name) + `</td>
			<td>` + fmt.Sprintf("%d", item.Quantity) + `</td>
			<td>` + fmt.Sprintf("%.2f", item.Price) + `</td>
			<td>` + fmt.Sprintf("%.2f", total) + `</td>
			<td>`

		// Add modification characteristics as a nested list
		if len(item.ModificationCharacteristics) > 0 {
			emailBody += `<ul>`
			for _, char := range item.ModificationCharacteristics {
				emailBody += `<li><strong>` + char.Name + `:</strong> ` + char.Value + `</li>`
			}
			emailBody += `</ul>`
		} else {
			emailBody += `Нет`
		}

		emailBody += `</td>
		</tr>`
	}

	emailBody += `</table>`

	// Closing
	emailBody += `<p></p>`

	return emailBody
}

// Helper function to format address
func formatAddress(customer models.CustomerInfo) string {
	var addressParts []string
	addressParts = append(addressParts, customer.City, customer.Street, "House "+customer.House)
	if customer.Entrance != "" {
		addressParts = append(addressParts, "Entrance "+customer.Entrance)
	}
	if customer.Floor != "" {
		addressParts = append(addressParts, "Floor "+customer.Floor)
	}
	if customer.Apartment != "" {
		addressParts = append(addressParts, "Apartment "+customer.Apartment)
	}
	return strings.Join(addressParts, ", ")
}

// Helper function to handle empty fields gracefully
func emptyOrValue(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

// func buildOrderEmailBody(order models.Order, customer models.CustomerInfo, items []models.OrderItem) string {
// 	// Start with customer info
// 	emailBody := `
// 	<h2>Детали нового заказа</h2>
// 	<p><strong>ID заказа:</strong> ` + fmt.Sprintf("%d", order.ID) + `</p>
// 	<p><strong>Дата заказа:</strong> ` + order.OrderDate.Format("2006-01-02 15:04:05") + `</p>
// 	<h3>Информация заказчика</h3>
// 	<ul>
// 		<li><strong>ФИО:</strong> ` + customer.FullName + `</li>
// 		<li><strong>Email:</strong> ` + customer.Email + `</li>
// 		<li><strong>Телефон:</strong> ` + customer.TelephoneNumber + `</li>
// 		<li><strong>Адрес:</strong> ` + formatAddress(customer) + `</li>
// 		<li><strong>Комментарии:</strong> ` + emptyOrValue(customer.Comment) + `</li>
// 	</ul>`
//
// 	// Include alternative contact details, if provided
// 	if customer.AnotherFullName != "" || customer.AnotherTelephoneNumber != "" {
// 		emailBody += `
// 		<h4>Альтернативная информация заказчика</h4>
// 		<ul>
// 			<li><strong>ФИО:</strong> ` + emptyOrValue(customer.AnotherFullName) + `</li>
// 			<li><strong>Телефон:</strong> ` + emptyOrValue(customer.AnotherTelephoneNumber) + `</li>
// 		</ul>`
// 	}
//
// 	// Add order summary
// 	emailBody += `
// 	<h3>Заказ</h3>
// 	<p><strong>Общая сумма:</strong> ` + fmt.Sprintf("%.2f", order.TotalAmount) + `</p>`
//
// 	// Add ordered items
// 	emailBody += `
// 	<h3>Заказанные товары</h3>
// 	<table border="1" cellpadding="8" cellspacing="0">
// 		<tr>
// 			<th>Название товара</th>
// 			<th>Количество</th>
// 			<th>Цена</th>
// 			<th>Всего</th>
// 		</tr>`
//
// 	for _, item := range items {
// 		total := float64(item.Quantity) * item.Price
// 		emailBody += `
// 		<tr>
// 			<td>` + fmt.Sprintf(item.Name) + `</td>
// 			<td>` + fmt.Sprintf("%d", item.Quantity) + `</td>
// 			<td>` + fmt.Sprintf("%.2f", item.Price) + `</td>
// 			<td>` + fmt.Sprintf("%.2f", total) + `</td>
// 		</tr>`
// 	}
//
// 	emailBody += `</table>`
//
// 	// Closing
// 	emailBody += `<p></p>`
//
// 	return emailBody
// }
