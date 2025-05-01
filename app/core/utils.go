package core

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
)

// Função para comparar dois floats com uma tolerância
func AlmostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
func GerarTokenApartirDeAppKey() string {
    app_key := fmt.Sprintf("%s",facades.Config().Env("APP_KEY"))
    h := hmac.New(sha256.New, []byte(app_key))
    var data string
    h.Write([]byte(data))
    token :=  hex.EncodeToString(h.Sum(nil))
    return token
}

func MesesEntreDatas(data1, data2 time.Time) int {
	// Calcula a diferença em anos e meses
	anos := data2.Year() - data1.Year()
	meses := int(data2.Month()) - int(data1.Month())

	// Calcula o total de meses
	totalMeses := anos*12 + meses

	// Ajusta se o dia da segunda data for menor que o dia da primeira data
	if data2.Day() < data1.Day() {
		totalMeses--
	}
    if totalMeses < 0 {
        totalMeses = totalMeses * -1
    }
	return totalMeses
}
func PorcentagemValorInicialParaValorFinal(valor_inicial, valor_final float64) float64 {
    return (valor_final / valor_inicial) * 100
}
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
func GetTaxaSelic() float64 {
    valorSelic := 13.25 // valor padrão

    result, err := HttpRequest("https://www.bcb.gov.br/api/servico/sitebcb//taxaselic/ultima?withCredentials=true", "GET",
        map[string]string{"content-type":"text/plain"}, "")
    if err != nil {
        return valorSelic
    }

    // Usando map[string]interface{} para evitar structs
    var response map[string]interface{}
    if err := ConverterJson(result, &response); err != nil {
        return valorSelic
    }

    conteudo, ok := response["conteudo"].([]interface{})
    if !ok || len(conteudo) == 0 {
        return valorSelic
    }

    primeiroItem, ok := conteudo[0].(map[string]interface{})
    if !ok {
        return valorSelic
    }

    if metaSelic, ok := primeiroItem["MetaSelic"].(float64); ok {
        valorSelic = metaSelic
    }

    return valorSelic
}
func QuantidadeDiasDeUmMes(data time.Time) int {
	firstOfNextMonth := time.Date(data.Year(), data.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	numDays := lastDayOfMonth.Day()
    return numDays
}
func EhMobile(userAgent string) bool {
	// Lista de palavras-chave que indicam um dispositivo móvel
	mobileKeywords := []string{
		"Android",
		"iPhone",
		"iPad",
		"Windows Phone",
		"BlackBerry",
		"Mobile",
	}

	// Verifica se o User-Agent contém alguma das palavras-chave
	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}

	return false
}
func RenderView(name string, data map[string]interface{}) ([]byte, error) {
    t, err := template.ParseFiles("resources/views/" + name + ".html")
    var buf bytes.Buffer
    if err != nil {
        return buf.Bytes(), err
    }
    if err := t.ExecuteTemplate(&buf, name, data); err != nil {
        return buf.Bytes(), err
    }
    return buf.Bytes(), nil
}
