package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	serviceName = "servico-b"
	zipkinURL   = "http://zipkin:9411/api/v2/spans" // Endereço do Zipkin no Docker Compose
)

// ViaCEPResponse representa a estrutura da resposta da API ViaCEP
type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

// WeatherAPIResponse representa a estrutura simplificada da resposta da WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
}

// CEPRequest representa a estrutura da requisição de CEP
type CEPRequest struct {
	CEP string `json:"cep"`
}

// WeatherResponse representa a estrutura da resposta final
type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func initTracer() *trace.TracerProvider {
	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		log.Fatalf("falha ao criar exportador Zipkin: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Erro ao encerrar tracer provider: %v", err)
		}
	}()

	http.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(handleWeather), "handleWeather"))

	log.Printf("Serviço B rodando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(serviceName).Start(r.Context(), "handleWeatherRequest")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	cep := req.CEP

	// Valida se o CEP tem 8 dígitos e é uma string
	if len(cep) != 8 || !regexp.MustCompile(`^[0-9]{8}$`).MatchString(cep) {
		span.SetAttributes(attribute.String("cep.validation.error", "invalid zipcode format"))
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	span.SetAttributes(attribute.String("cep", cep))

	// 1. Consultar ViaCEP para obter a cidade
	city, err := getCityFromCEP(ctx, cep)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "can not find zipcode" {
			w.WriteHeader(http.StatusNotFound) // 404
			json.NewEncoder(w).Encode(map[string]string{"message": "can not find zipcode"})
		} else {
			http.Error(w, fmt.Sprintf("Erro ao buscar cidade: %v", err), http.StatusInternalServerError)
		}
		return
	}

	span.SetAttributes(attribute.String("city", url.QueryEscape(city)))

	// 2. Consultar WeatherAPI para obter a temperatura
	tempC, err := getTemperatureFromWeatherAPI(ctx, url.QueryEscape(city))
	if err != nil {
		span.RecordError(err)
		http.Error(w, fmt.Sprintf("Erro ao buscar temperatura: %v", err), http.StatusInternalServerError)
		return
	}

	// 3. Calcular Fahrenheit e Kelvin
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15 // Usando 273.15 para maior precisão

	response := WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(response)
}

func getCityFromCEP(ctx context.Context, cep string) (string, error) {
	ctx, span := otel.Tracer(serviceName).Start(ctx, "getCityFromCEP")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição ViaCEP: %w", err)
	}

	client := otelhttp.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao chamar ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ViaCEP retornou status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta ViaCEP: %w", err)
	}

	var viaCEPResp ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResp); err != nil {
		return "", fmt.Errorf("erro ao fazer unmarshal da resposta ViaCEP: %w", err)
	}

	if viaCEPResp.Erro {
		return "", fmt.Errorf("can not find zipcode")
	}

	span.SetAttributes(attribute.String("viacep.city", viaCEPResp.Localidade))

	return viaCEPResp.Localidade, nil
}

func getTemperatureFromWeatherAPI(ctx context.Context, city string) (float64, error) {
	ctx, span := otel.Tracer(serviceName).Start(ctx, "getTemperatureFromWeatherAPI")
	defer span.End()

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=f66eea1616bb484fad3183814252705&q=%s", city)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar requisição WeatherAPI: %w", err)
	}

	client := otelhttp.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("erro ao chamar WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("WeatherAPI error response: %s", string(body))
		return 0, fmt.Errorf("WeatherAPI retornou status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta WeatherAPI: %w", err)
	}

	var weatherResp WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return 0, fmt.Errorf("erro ao fazer unmarshal da resposta WeatherAPI: %w", err)
	}

	span.SetAttributes(attribute.Float64("weatherapi.temp_c", weatherResp.Current.TempC))
	return weatherResp.Current.TempC, nil
}
