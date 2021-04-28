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
	mirrorsConfig *map[string]string
	server        *http.Server
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func (whsvr *WebhookServer) generateImagePatch(containerNum int, origImage string, newImage string) (patch patchOperation) {
	return patchOperation{
		Op:    "replace",
		Path:  fmt.Sprintf("/spec/containers/%d/image", containerNum),
		Value: newImage,
	}
}

func (whsvr *WebhookServer) updateMirrorsConfig(newconfig *map[string]string) {
	log.Debug("Callback called, updating mirrors config")
	whsvr.mirrorsConfig = newconfig
}

func (whsvr *WebhookServer) mutate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	log.Debugf("Mutate called: %v", whsvr.mirrorsConfig)
	req := ar.Request
	var pod corev1.Pod

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		log.WithError(err).Error("Could not unmarshal raw object")
		return &admissionv1.AdmissionResponse{
			Allowed:  true,
			Warnings: []string{err.Error()},
		}
	}

	var patches []patchOperation

	for i := 0; i < len(pod.Spec.Containers); i++ {

		thisImage := pod.Spec.Containers[i].Image

		log.WithFields(log.Fields{
			"totalContainers": len(pod.Spec.Containers),
			"containerNumber": i,
			"image":           thisImage,
		}).Info("Container Image Found")

		parsedImage, err := dockerparser.Parse(thisImage)
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
		log.Infof("Mirrored image located at: %s", newImage)

		patches = append(patches, whsvr.generateImagePatch(i, thisImage, newImage))
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
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Errorf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *admissionv1.AdmissionResponse
	ar := admissionv1.AdmissionReview{}
	_, _, err := deserializer.Decode(body, nil, &ar)
	if err != nil {
		log.Errorf("Can't decode body: %v", err)
		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	admissionResponse = whsvr.mutate(&ar)

	admissionReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       ar.Kind,
			APIVersion: ar.APIVersion,
			// Kind:       "AdmissionReview",
			// APIVersion: "admission.k8s.io/v1",
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

	log.Info("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		log.WithError(err).Error("Can't write response")
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
