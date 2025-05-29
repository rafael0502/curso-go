package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	serviceName = "servico-a"
	zipkinURL   = "http://zipkin:9411/api/v2/spans" // Endereço do Zipkin no Docker Compose
)

// CEPRequest representa a estrutura da requisição de CEP
type CEPRequest struct {
	CEP string `json:"cep"`
}

// CEPResponse representa a estrutura da resposta do Serviço B
type CEPResponse struct {
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

	http.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(handleCEP), "handleCEP"))

	log.Printf("Serviço A rodando na porta 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(serviceName).Start(r.Context(), "handleCEPRequest")
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

	jsonBody, _ := json.Marshal(map[string]string{"cep": cep})
	reqServiceB, err := http.NewRequestWithContext(ctx, "POST", "http://servicoB:8081/weather", bytes.NewBuffer(jsonBody))
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Erro ao criar requisição para o Serviço B", http.StatusInternalServerError)
		return
	}
	reqServiceB.Header.Set("Content-Type", "application/json")

	client := otelhttp.DefaultClient
	respServiceB, err := client.Do(reqServiceB)
	if err != nil {
		span.RecordError(err)
		http.Error(w, fmt.Sprintf("Erro ao se comunicar com o Serviço B: %v", err), http.StatusInternalServerError)
		return
	}
	defer respServiceB.Body.Close()

	bodyBytes, err := ioutil.ReadAll(respServiceB.Body)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Erro ao ler resposta do Serviço B", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(respServiceB.StatusCode)
	w.Write(bodyBytes)
}
