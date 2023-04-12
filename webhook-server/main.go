package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	port                = GetEnvDefault("WEBHOOK_PORT", "8443")
	env                 = GetEnvDefault("ENV", "dev")
	skippableNamespaces = func() []string {
		def := []string{metav1.NamespacePublic, metav1.NamespaceSystem}
		v := os.Getenv("SKIP_NAMESPACE")
		if v == "" {
			return def
		}

		return append(def, strings.Split(v, ",")...)
	}()
	tlsDir          = os.Getenv("TLS_DIR")
	allowScheduling = GetEnvDefault("ALLOW_SCHEDULING", "false") == "true"
)

const (
	jsonContentType = "application/json"
)

func main() {
	certFile := filepath.Join(tlsDir, "tls.crt")
	keyFile := filepath.Join(tlsDir, "tls.key")

	mux := http.NewServeMux()
	mux.Handle("/validate", validateHandler())
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Println("listening")

	// TODO: add graceful shutdown
	// TODO: add logger
	// TODO: unit testing
	if env == "dev" {
		log.Fatal(server.ListenAndServe())
		return
	}

	log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}

func validateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only POST method is supported
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(responseBody("invalid method %s, only POST requests are allowed", r.Method))

			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBody("could not ready body %v", err))

			return
		}
		if contentType := r.Header.Get("Content-Type"); contentType != jsonContentType {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(
				responseBody(
					"unsupported content type %s, only %s is supported",
					contentType,
					jsonContentType,
				),
			)

			return
		}

		var review v1.AdmissionReview

		err = json.Unmarshal(body, &review)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBody("could not deserialize request: %v", err))

			return
		}

		// Step 3: Construct the AdmissionReview response.
		admissionReviewResponse := v1.AdmissionReview{
			// Since the admission webhook now supports multiple API versions, we need
			// to explicitly include the API version in the response.
			// This API version needs to match the version from the request exactly, otherwise
			// the API server will be unable to process the response.
			// Note: a v1beta1 AdmissionReview is JSON-compatible with the v1 version, that's why
			// we do not need to differentiate during unmarshaling or in the actual logic.
			TypeMeta: review.TypeMeta,
			Response: &v1.AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: true,
			},
		}

		if !skipNamespace(review.Request.Namespace) {
			raw := review.Request.Object.Raw
			pod := corev1.Pod{}

			err := json.Unmarshal(raw, &pod)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(responseBody("could not unmarshal pod spec: %v", err))
				return
			}

			if _, exists := pod.Labels["team"]; !exists {
				admissionReviewResponse.Response.Allowed = allowScheduling
				if allowScheduling {
					admissionReviewResponse.Response.Warnings = []string{
						"Team label not set on pod",
					}
				}
				admissionReviewResponse.Response.Result = &metav1.Status{
					Status:  "Failure",
					Message: "Team label not set on pod",
				}

				log.Printf("Team label not set on pod: %s\n", pod.Name)
			}
		}

		response, err := json.Marshal(admissionReviewResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(responseBody("could not marshal JSON response: %v", err))

			return
		}

		w.Write(response)
	})
}

func responseBody(format string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(format, args...))
}

func GetEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}

func skipNamespace(ns string) bool {
	for _, n := range skippableNamespaces {
		if n == ns {
			return true
		}
	}

	return false
}
