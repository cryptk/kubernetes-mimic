package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	dockerparser "github.com/novln/docker-parser"
	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

type WebhookServer struct {
	server        *http.Server
	mirrorsConfig *map[string]string
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func (whsvr *WebhookServer) generatePatch(index int, container corev1.Container, newImage string) (patch []patchOperation) {
	patch = []patchOperation{
		{
			Op:    "add",
			Path:  fmt.Sprintf("/metadata/annotations/kubernetes-mimic.io~1%s-original-image", container.Name),
			Value: container.Image,
		},
		{
			Op:    "replace",
			Path:  fmt.Sprintf("/spec/containers/%d/image", index),
			Value: newImage,
		},
	}

	return patch
}

func (whsvr *WebhookServer) updateMirrorsConfig(newconfig *map[string]string) {
	log.Debug("Callback called, updating mirrors config")

	whsvr.mirrorsConfig = newconfig
}

func (whsvr *WebhookServer) mutate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	log.Debugf("Mutate called: %v", whsvr.mirrorsConfig)

	var pod corev1.Pod

	req := ar.Request

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		log.WithError(err).Error("Could not unmarshal raw object")

		return &admissionv1.AdmissionResponse{
			Allowed:  true,
			Warnings: []string{err.Error()},
		}
	}

	var patches []patchOperation

	for i := 0; i < len(pod.Spec.Containers); i++ {
		container := pod.Spec.Containers[i]

		log.WithFields(log.Fields{
			"totalContainers": len(pod.Spec.Containers),
			"containerNumber": i,
			"image":           container.Image,
		}).Info("Container Image Found")

		parsedImage, err := dockerparser.Parse(container.Image)
		if err != nil {
			log.WithError(err).Error("Failed to parse docker image name")
			return nil
		}

		log.WithFields(log.Fields{
			"name":       parsedImage.Name(),
			"registry":   parsedImage.Registry(),
			"remote":     parsedImage.Remote(),
			"repository": parsedImage.Repository(),
			"shortname":  parsedImage.ShortName(),
			"tag":        parsedImage.Tag(),
		}).Debug("Image name parsed")

		mirror, ok := (*whsvr.mirrorsConfig)[parsedImage.Registry()]
		if !ok {
			log.WithField("remote", parsedImage.Remote()).Info("No mirror configured for image")

			// Just because we don't have a mirror configured doesn't mean we should prevent the pod from running
			return &admissionv1.AdmissionResponse{
				Allowed: true,
			}
		}

		newImage := fmt.Sprintf("%s/%s", mirror, parsedImage.Name())
		log.WithField("image", newImage).Info("Mirrored located")

		patches = append(patches, whsvr.generatePatch(i, container, newImage)...)
	}

	log.WithFields(log.Fields{
		"kind":           req.Kind,
		"namespace":      req.Namespace,
		"name":           req.Name,
		"podName":        pod.Name,
		"uid":            req.UID,
		"patchOperation": req.Operation,
		"userInfo":       req.UserInfo,
	}).Debug("Processing AdmissionReview")

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		log.Info("Failed to marshal JSON for patch response")

		return &admissionv1.AdmissionResponse{
			Allowed: true,
		}
	}

	log.Debugf("AdmissionResponse: patch=%v\n", string(patchBytes))

	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// Serve method for webhook server
func (whsvr *WebhookServer) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte

	if r.Body == nil {
		log.Error("Request body is nil")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Request body read failure")
		return
	}

	if len(body) == 0 {
		log.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)

		return
	}

	// verify the content type is accurate
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		log.WithField("content-type", contentType).Error("invalid content-type, expected application/json")
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)

		return
	}

	var admissionResponse *admissionv1.AdmissionResponse

	ar := admissionv1.AdmissionReview{}

	_, _, err = deserializer.Decode(body, nil, &ar)
	if err != nil {
		log.WithError(err).Error("can't decode body")

		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		admissionResponse = whsvr.mutate(&ar)
	}

	admissionReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       ar.Kind,
			APIVersion: ar.APIVersion,
		},
	}

	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		log.WithError(err).Error("Can't encode response")
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}

	log.Info("Ready to write response ...")

	if _, err := w.Write(resp); err != nil {
		log.WithError(err).Error("Can't write response")
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
