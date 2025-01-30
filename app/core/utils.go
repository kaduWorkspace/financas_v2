package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)
// Função para calcular os dias restantes até o final do ano
func DiasRestantesNoAno(t time.Time) int {
	// Obter o último dia do ano
	ultimoDiaDoAno := time.Date(t.Year(), 12, 31, 23, 59, 59, 0, t.Location())

	// Calcular a diferença em dias
	diasRestantes := int(ultimoDiaDoAno.Sub(t).Hours() / 24)
	return diasRestantes
}
func DiasNoAnoV2(t time.Time) int {
	ano := t.Year()
	// Verifica se o ano é bissexto
	if (ano%4 == 0 && ano%100 != 0) || (ano%400 == 0) {
		return 366 // Ano bissexto
	}
	return 365 // Ano normal
}
func Dias_no_ano(year int) int {
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := startOfYear.AddDate(1, 0, 0) // Primeiro dia do próximo ano
	return int(endOfYear.Sub(startOfYear).Hours() / 24) // Diferença em dias
}
func FormataData(data time.Time) string {
    return data.Format("02/01/2006")
}
func FormatarValorMonetario(value float64) string {
	str := strconv.FormatFloat(value, 'f', 2, 64)

	str = strings.Replace(str, ".", ",", 1)

	parts := strings.Split(str, ",")
	intPart := parts[0]
	decimalPart := parts[1]

	var result string
	count := 0
	for i := len(intPart) - 1; i >= 0; i-- {
		result = string(intPart[i]) + result
		count++
		if count%3 == 0 && i > 0 {
			result = "." + result
		}
	}

	return result + "," + decimalPart
}
func HttpRequest(url string, method string, headers map[string]string, body string) (string, error) {
	// Define o método padrão como GET, caso nenhum seja passado
	if method == "" {
		method = "GET"
	}

	// Cria o payload para o corpo da requisição
	var requestBody *bytes.Reader
	if body != "" {
		requestBody = bytes.NewReader([]byte(body))
	} else {
		requestBody = bytes.NewReader(nil)
	}

	// Cria a requisição HTTP
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao criar a requisição: %w", err)
	}

	// Adiciona os cabeçalhos à requisição
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Realiza a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao realizar a requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler a resposta: %w", err)
	}

	return string(respBody), nil
}
func ConverterJson[v any](input string, destination *v) error {
	// Verifica se o destino é um ponteiro válido
	destValue := reflect.ValueOf(destination)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return fmt.Errorf("o destino deve ser um ponteiro válido")
	}

	// Desserializa o JSON para o destino
    if err := json.Unmarshal([]byte(input), destination); err != nil {
		return fmt.Errorf("falha ao desserializar JSON para destino: %w", err)
	}

	return nil
}
func DiasNoAno(data string) (int, error) {
	partes := strings.Split(data, "/")
	if len(partes) != 3 {
		return 0, fmt.Errorf("data inválida, o formato deve ser dd/mm/yyyy")
	}
	ano, err := strconv.Atoi(partes[2])
	if err != nil {
		return 0, fmt.Errorf("ano inválido: %v", err)
	}
	if (ano%4 == 0 && ano%100 != 0) || (ano%400 == 0) {
		return 366, nil
	}
	return 365, nil
}
